package authhandler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/niksmo/gophkeeper/internal/client/command"
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
	logger  logger.Logger
	service UserRegistrar
	writer  io.Writer
}

func NewSignup(l logger.Logger, s UserRegistrar, w io.Writer) *SignupHandler {
	return &SignupHandler{l, s, w}
}

func (h *SignupHandler) Handle(ctx context.Context, v command.ValueGetter) {
	const op = "SignupHandler.Handle"
	log := h.logger.WithOp(op)

	login, password, err := getAuthFlags(v)
	if err != nil {
		log.Debug().Err(err).Send()
		fmt.Fprintln(h.writer, err.Error())
		os.Exit(1)
	}

	// TODO: verify login and password to match pattern

	err = h.service.RegisterUser(ctx, login, password)
	if err != nil {
		switch {
		case errors.Is(err, authservice.ErrAlreadyExists):
			h.printOutput("the login '%s' already exists", login)
			os.Exit(1)
		case errors.Is(err, authservice.ErrSyncAlreadyRunning):
			h.printOutput(
				"synchronization is working, logout and login for restart",
			)
			os.Exit(1)
		default:
			log.Debug().Err(err).Msg("failed to register new user account")
			handler.InternalError(
				h.writer, fmt.Errorf("%s: %w", op, err),
			)
			os.Exit(1)
		}
	}

	h.printOutput("synchronization started")
}

func (h *SignupHandler) printOutput(formated string, args ...any) {
	fmt.Fprintf(h.writer, formated, args...)
	fmt.Fprintln(h.writer)
}

type SigninHandler struct {
	logger  logger.Logger
	service UserAuthorizer
	writer  io.Writer
}

func NewSignin(l logger.Logger, s UserAuthorizer, w io.Writer) *SigninHandler {
	return &SigninHandler{l, s, w}
}

func (h *SigninHandler) Handle(ctx context.Context, v command.ValueGetter) {
	const op = "SigninHandler.Handle"
	log := h.logger.WithOp(op)

	login, password, err := getAuthFlags(v)
	if err != nil {
		log.Debug().Err(err).Send()
		fmt.Fprintln(h.writer, err.Error())
		os.Exit(1)
	}

	// TODO: verify login and password to match pattern

	err = h.service.AuthorizeUser(ctx, login, password)
	if err != nil {
		switch {
		case errors.Is(err, authservice.ErrCredentials):
			h.printOutput("invalid login or password")
			os.Exit(1)
		case errors.Is(err, authservice.ErrSyncAlreadyRunning):
			h.printOutput("synchronization is working, logout and login for restart")
			os.Exit(1)
		default:
			log.Debug().Err(err).Msg("failed to login")
			handler.InternalError(
				h.writer, fmt.Errorf("%s: %w", op, err),
			)
			os.Exit(1)
		}
	}

	h.printOutput("synchronization started")
}

func (h *SigninHandler) printOutput(formated string, args ...any) {
	fmt.Fprintf(h.writer, formated, args...)
	fmt.Fprintln(h.writer)
}

type LogoutHandler struct {
	logger  logger.Logger
	service SyncCloser
	writer  io.Writer
}

func NewLogout(l logger.Logger, s SyncCloser, w io.Writer) *LogoutHandler {
	return &LogoutHandler{l, s, w}
}

func (h *LogoutHandler) Handle(ctx context.Context, v command.ValueGetter) {
	const op = "LogoutHandler.Handle"
	log := h.logger.WithOp(op)

	err := h.service.CloseSynchronization(ctx)
	if err != nil {
		if errors.Is(err, syncservice.ErrNoSync) {
			h.printOutput("synchronization is not running")
			os.Exit(1)
		}
		log.Debug().Err(err).Msg("failed to close synchronization")
		handler.InternalError(h.writer, fmt.Errorf("%s: %w", op, err))
		os.Exit(1)
	}

	h.printOutput("synchronization stopped")
}

func (h *LogoutHandler) printOutput(formated string, args ...any) {
	fmt.Fprintf(h.writer, formated, args...)
	fmt.Fprintln(h.writer)
}

func getAuthFlags(v command.ValueGetter) (l, p string, err error) {
	var errs []error
	l, err = v.GetString(synccommand.LoginFlag)
	if err != nil {
		errs = append(errs, fmt.Errorf("--%s", synccommand.LoginFlag))
	}

	p, err = v.GetString(synccommand.PasswordFlag)
	if err != nil {
		errs = append(errs, fmt.Errorf("--%s", synccommand.PasswordFlag))
	}

	err = handler.RequiredFlagsErr(errs)
	return
}
