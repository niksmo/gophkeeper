package pwdhandler

import (
	"context"
	"errors"
	"fmt"
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

	pwdAddHandler struct {
		l logger.Logger
		s pwdAddService
	}
)

func NewAddHandler(l logger.Logger, s pwdAddService) *pwdAddHandler {
	return &pwdAddHandler{l, s}
}

func (h *pwdAddHandler) Handle(ctx context.Context, v command.ValueGetter) {
	const op = "pwdAddHandler.Handle"
	log := h.l.With().Str("op", op).Logger()

	key, name, password, login, err := h.getStrFlagValues(v)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	o := objects.PWD{
		Name:     name,
		Login:    login,
		Password: password,
	}

	entryNum, err := h.s.Add(ctx, key, o)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to save password")
	}

	fmt.Println("Password saved. Entry number:", entryNum)
}

func (h *pwdAddHandler) getStrFlagValues(
	v command.ValueGetter,
) (k, n, p, l string, errs error) {
	var err error
	k, err = v.GetString(pwdcommand.MasterKeyFlag)
	if err != nil || len(strings.TrimSpace(k)) == 0 {
		errs = errors.Join(
			fmt.Errorf("--%s flag is required", pwdcommand.MasterKeyFlag),
		)
	}

	n, err = v.GetString(pwdcommand.NameFlag)
	if err != nil || len(strings.TrimSpace(n)) == 0 {
		errs = errors.Join(
			fmt.Errorf("--%s flag is required", pwdcommand.NameFlag),
		)
	}

	p, err = v.GetString(pwdcommand.PasswordFlag)
	if err != nil || len(strings.TrimSpace(p)) == 0 {
		errs = errors.Join(
			fmt.Errorf("--%s flag is required", pwdcommand.PasswordFlag),
		)
	}

	l, _ = v.GetString(pwdcommand.LoginFlag)

	return
}
