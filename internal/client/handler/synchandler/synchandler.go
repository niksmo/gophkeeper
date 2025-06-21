package synchandler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/niksmo/gophkeeper/internal/client/command/synccommand"
	"github.com/niksmo/gophkeeper/internal/client/handler"
	"github.com/niksmo/gophkeeper/internal/client/service"
	"github.com/niksmo/gophkeeper/internal/client/service/syncservice"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type (
	SignupService interface {
		Signup(ctx context.Context, login, password string) error
	}

	SigninService interface {
		Signin(ctx context.Context, login, password string) error
	}

	LogoutService interface {
		Logout(ctx context.Context) error
	}
)

type SignupHandler struct {
	Log     logger.Logger
	Service SignupService
	Writer  io.Writer
}

func NewSignup(l logger.Logger, s SignupService, w io.Writer) *SignupHandler {
	return &SignupHandler{l, s, w}
}

func (h *SignupHandler) Handle(ctx context.Context, v command.ValueGetter) {
	const op = "SignupHandler.Handle"
	log := h.Log.With().Str("op", op).Logger()

	l, p, err := getAuthFlags(v)
	if err != nil {
		log.Debug().Err(err).Send()
		fmt.Fprintln(h.Writer, err.Error())
		os.Exit(1)
	}

	err = h.Service.Signup(ctx, l, p)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAlreadyExists):
			fmt.Fprintf(h.Writer, "the login '%s' already exists\n", l)
			os.Exit(1)
		case errors.Is(err, syncservice.ErrPIDConflict):
			fmt.Fprintln(
				h.Writer,
				"synchronization start error, please logout and then login",
			)
			os.Exit(1)
		default:
			log.Debug().Err(err).Msg("failed to register new user account")
			handler.InternalError(h.Writer, err)
			os.Exit(1)
		}
	}

	writeSyncStart(h.Writer)
}

type SigninHandler struct {
	Log     logger.Logger
	Service SigninService
	Writer  io.Writer
}

func NewSignin(l logger.Logger, s SigninService, w io.Writer) *SigninHandler {
	return &SigninHandler{l, s, w}
}

func (h *SigninHandler) Handle(ctx context.Context, v command.ValueGetter) {
	const op = "SigninHandler.Handle"
	log := h.Log.With().Str("op", op).Logger()

	l, p, err := getAuthFlags(v)
	if err != nil {
		log.Debug().Err(err).Send()
		fmt.Fprintln(h.Writer, err.Error())
		os.Exit(1)
	}

	err = h.Service.Signin(ctx, l, p)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCredentials):
			fmt.Fprintln(h.Writer, "invalid login or password")
			os.Exit(1)
		default:
			log.Debug().Err(err).Msg("failed to login")
			handler.InternalError(h.Writer, err)
			os.Exit(1)
		}
	}

	writeSyncStart(h.Writer)
}

type LogoutHandler struct {
	Log     logger.Logger
	Service LogoutService
	Writer  io.Writer
}

func NewLogout(l logger.Logger, s LogoutService, w io.Writer) *LogoutHandler {
	return &LogoutHandler{l, s, w}
}

func (h *LogoutHandler) Handle(ctx context.Context, v command.ValueGetter) {
	const op = "LogoutHandler.Handle"
	log := h.Log.With().Str("op", op).Logger()

	err := h.Service.Logout(ctx)
	if err != nil {
		if errors.Is(err, syncservice.ErrNoSync) {
			fmt.Fprintln(h.Writer, "synchronization is not running")
			return
		}
		log.Debug().Err(err).Msg("failed to logout")
		handler.InternalError(h.Writer, err)
		os.Exit(1)
	}

	fmt.Fprintln(h.Writer, "synchronization stopped")

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

func writeSyncStart(w io.Writer) {
	fmt.Fprintf(w, "synchronization started\n")
}
