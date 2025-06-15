package pwdhandler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/niksmo/gophkeeper/internal/client/command/pwdcommand"
	"github.com/niksmo/gophkeeper/internal/client/objects"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type (
	pwdAddService interface {
		Add(ctx context.Context, key string, obj objects.PWD) (int, error)
	}

	PwdAddHandler struct {
		l logger.Logger
		s pwdAddService
		w io.Writer
	}
)

func NewAddHandler(
	l logger.Logger, s pwdAddService, w io.Writer,
) *PwdAddHandler {
	return &PwdAddHandler{l, s, w}
}

func (h *PwdAddHandler) Handle(ctx context.Context, v command.ValueGetter) {
	const op = "pwdAddHandler.Handle"
	log := h.l.With().Str("op", op).Logger()

	key, name, password, login, err := h.getStrFlagValues(v)
	if err != nil {
		log.Debug().Err(err).Send()
		fmt.Fprintln(h.w, err.Error())
		os.Exit(1)
	}

	o := objects.PWD{
		Name:     name,
		Login:    login,
		Password: password,
	}

	entryNum, err := h.s.Add(ctx, key, o)
	if err != nil {
		log.Debug().Err(err).Msg("failed to save password")
		fmt.Fprintf(
			h.w,
			"the password is not saved, application is crashed: %s\n",
			err.Error(),
		)
		os.Exit(1)
	}

	log.Debug().Int("entry", entryNum).Msg("password saved")
	fmt.Fprintf(
		h.w,
		"the password is saved under the record number: %d\n",
		entryNum,
	)
}

func (h *PwdAddHandler) getStrFlagValues(
	v command.ValueGetter,
) (k, n, p, l string, err error) {
	var errs []error
	k, err = v.GetString(pwdcommand.MasterKeyFlag)
	if err != nil || len(strings.TrimSpace(k)) == 0 {
		errs = append(errs, fmt.Errorf("--%s", pwdcommand.MasterKeyFlag))

	}

	n, err = v.GetString(pwdcommand.NameFlag)
	if err != nil || len(strings.TrimSpace(n)) == 0 {
		errs = append(errs, fmt.Errorf("--%s", pwdcommand.NameFlag))

	}

	p, err = v.GetString(pwdcommand.PasswordFlag)
	if err != nil || len(strings.TrimSpace(p)) == 0 {
		errs = append(errs, fmt.Errorf("--%s", pwdcommand.PasswordFlag))
	}

	l, _ = v.GetString(pwdcommand.LoginFlag)

	if len(errs) != 0 {
		err = errors.Join(errs...)
		err = fmt.Errorf("required flags are not specified:\n%w", err)
	}

	return
}
