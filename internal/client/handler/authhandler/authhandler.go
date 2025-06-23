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
	"github.com/niksmo/gophkeeper/internal/client/service/syncservice"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

var ErrServerConn = errors.New("synchronization server connection error")

type (
	UserRegistrar interface {
		RegisterUser(ctx context.Context, login, password string) error
	}

	UserAuthorizer interface {
		AuthorizeUser(ctx context.Context, login, password string) error
	}

	SyncStopper interface {
		StopSynchronization(ctx context.Context) error
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

	err = h.service.RegisterUser(ctx, login, password)
	if err != nil {
		// switch {
		// case errors.Is(err, service.ErrAlreadyExists):
		// 	fmt.Fprintf(h.Writer, "the login '%s' already exists\n", l)
		// 	os.Exit(1)
		// case errors.Is(err, syncservice.ErrPIDConflict):
		// 	fmt.Fprintln(
		// 		h.Writer,
		// 		"synchronization start error, please logout and then login",
		// 	)
		// 	os.Exit(1)
		// default:
		// 	log.Debug().Err(err).Msg("failed to register new user account")
		// 	handler.InternalError(
		// 		h.Writer, fmt.Errorf("%w: %w", ErrServerConn, err),
		// 	)
		// 	os.Exit(1)
		// }
	}

	writeOutput(h.writer)
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

	err = h.service.AuthorizeUser(ctx, login, password)
	if err != nil {
		// switch {
		// case errors.Is(err, service.ErrCredentials):
		// 	fmt.Fprintln(h.Writer, "invalid login or password")
		// 	os.Exit(1)
		// case errors.Is(err, syncservice.ErrPIDConflict):
		// 	fmt.Fprintln(
		// 		h.Writer,
		// 		"synchronization is working, logout and login for restart",
		// 	)
		// 	return
		// default:
		// 	log.Debug().Err(err).Msg("failed to login")
		// 	handler.InternalError(
		// 		h.Writer, fmt.Errorf("%w: %w", ErrServerConn, err),
		// 	)
		// 	os.Exit(1)
		// }
	}

	writeOutput(h.writer)
}

type LogoutHandler struct {
	logger  logger.Logger
	service SyncStopper
	writer  io.Writer
}

func NewLogout(l logger.Logger, s SyncStopper, w io.Writer) *LogoutHandler {
	return &LogoutHandler{l, s, w}
}

func (h *LogoutHandler) Handle(ctx context.Context, v command.ValueGetter) {
	const op = "LogoutHandler.Handle"
	log := h.logger.WithOp(op)

	err := h.service.StopSynchronization(ctx)
	if err != nil {
		if errors.Is(err, syncservice.ErrNoSync) {
			fmt.Fprintln(h.writer, "synchronization is not running")
			return
		}
		log.Debug().Err(err).Msg("failed to stop synchronization")
		handler.InternalError(h.writer, err)
		os.Exit(1)
	}

	fmt.Fprintln(h.writer, "synchronization stopped")
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

func writeOutput(w io.Writer) {
	fmt.Fprintln(w, "synchronization started")
}
