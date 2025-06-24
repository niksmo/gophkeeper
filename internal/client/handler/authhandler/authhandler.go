package authhandler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/niksmo/gophkeeper/internal/client/command/synccommand"
	"github.com/niksmo/gophkeeper/internal/client/handler"
	"github.com/niksmo/gophkeeper/internal/client/service/authservice"
	"github.com/niksmo/gophkeeper/internal/client/service/syncservice"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type (
	UserRegistrar interface {
		RegisterUser(ctx context.Context, login, password string) error
	}

	UserAuthorizer interface {
		AuthorizeUser(ctx context.Context, login, password string) error
	}

	SyncCloser interface {
		CloseSynchronization(context.Context) error
	}
)

type SignupHandler struct {
	l logger.Logger
	s UserRegistrar
	w io.Writer
}

func NewSignup(l logger.Logger, s UserRegistrar, w io.Writer) *SignupHandler {
	return &SignupHandler{l, s, w}
}

func (h *SignupHandler) Handle(ctx context.Context, fv synccommand.AuthFlags) {
	const op = "SignupHandler.Handle"

	log := h.l.WithOp(op)

	// TODO: verify login and password to match pattern

	err := h.s.RegisterUser(ctx, fv.Login, fv.Password)
	if err != nil {
		handler.HandleAlreadyExistsErr(err, log, h.w, "account", "login")
		handleSyncRunningErr(err, h.w)
		handler.HandleUnexpectedErr(err, log, h.w)
	}

	h.printOutput("synchronization started")
}

func (h *SignupHandler) printOutput(formated string, args ...any) {
	fmt.Fprintf(h.w, formated, args...)
	fmt.Fprintln(h.w)
}

type SigninHandler struct {
	l logger.Logger
	s UserAuthorizer
	w io.Writer
}

func NewSignin(l logger.Logger, s UserAuthorizer, w io.Writer) *SigninHandler {
	return &SigninHandler{l, s, w}
}

func (h *SigninHandler) Handle(ctx context.Context, fv synccommand.AuthFlags) {
	const op = "SigninHandler.Handle"

	log := h.l.WithOp(op)

	// TODO: verify login and password to match pattern

	err := h.s.AuthorizeUser(ctx, fv.Login, fv.Password)
	if err != nil {
		h.handleCredentialsErr(err)
		handleSyncRunningErr(err, h.w)
		handler.HandleUnexpectedErr(err, log, h.w)
	}

	h.printOutput("synchronization started")
}

func (h *SigninHandler) handleCredentialsErr(err error) {
	if !errors.Is(err, authservice.ErrCredentials) {
		return
	}
	h.printOutput("invalid login or password")
	os.Exit(1)
}

func (h *SigninHandler) printOutput(formated string, args ...any) {
	fmt.Fprintf(h.w, formated, args...)
	fmt.Fprintln(h.w)
}

type LogoutHandler struct {
	l logger.Logger
	s SyncCloser
	w io.Writer
}

func NewLogout(l logger.Logger, s SyncCloser, w io.Writer) *LogoutHandler {
	return &LogoutHandler{l, s, w}
}

func (h *LogoutHandler) Handle(ctx context.Context) {
	const op = "LogoutHandler.Handle"

	log := h.l.WithOp(op)

	err := h.s.CloseSynchronization(ctx)
	if err != nil {
		h.handleNoSyncErr(err)
		handler.HandleUnexpectedErr(err, log, h.w)
	}

	h.printOutput("synchronization stopped")
}

func (h *LogoutHandler) handleNoSyncErr(err error) {
	if !errors.Is(err, syncservice.ErrNoSync) {
		return
	}
	h.printOutput("synchronization is not running")
	os.Exit(1)
}

func (h *LogoutHandler) printOutput(formated string, args ...any) {
	fmt.Fprintf(h.w, formated, args...)
	fmt.Fprintln(h.w)
}

func handleSyncRunningErr(err error, w io.Writer) {
	if !errors.Is(err, authservice.ErrSyncAlreadyRunning) {
		return
	}

	fmt.Fprintln(w, "synchronization is working, logout and login for restart")
	os.Exit(1)

}
