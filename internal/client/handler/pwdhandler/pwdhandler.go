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
	"github.com/niksmo/gophkeeper/internal/client/handler"
	"github.com/niksmo/gophkeeper/internal/client/objects"
	"github.com/niksmo/gophkeeper/internal/client/service/pwdservice"
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

	key, name, password, login, err := h.getFlagValues(v)
	if err != nil {
		log.Debug().Err(err).Send()
		fmt.Fprintln(h.w, err.Error())
		os.Exit(1)
	}

	obj := objects.PWD{
		Name:     name,
		Login:    login,
		Password: password,
	}

	entryNum, err := h.s.Add(ctx, key, obj)
	if err != nil {
		log.Debug().Err(err).Msg("failed to save password")
		handler.InternalError(h.w, err)
		os.Exit(1)
	}

	log.Debug().Int("entry", entryNum).Msg("password saved")
	fmt.Fprintf(
		h.w,
		"the password is saved under the record number: %d\n",
		entryNum,
	)
}

func (h *PwdAddHandler) getFlagValues(
	v command.ValueGetter,
) (k, n, p, l string, err error) {
	var errs []error
	k, err = getMKey(v)
	if err != nil {
		errs = append(errs, err)
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

	err = handler.RequiredFlagsErr(errs)

	return k, n, p, l, err
}

type (
	pwdReadService interface {
		Read(ctx context.Context, key string, id int) (objects.PWD, error)
	}

	PwdReadHandler struct {
		l logger.Logger
		s pwdReadService
		w io.Writer
	}
)

func NewReadHandler(
	l logger.Logger, s pwdReadService, w io.Writer,
) *PwdReadHandler {
	return &PwdReadHandler{l, s, w}
}

func (h *PwdReadHandler) Handle(ctx context.Context, v command.ValueGetter) {
	const op = "pwdReadHandler.Handle"
	log := h.l.With().Str("op", op).Logger()

	k, e, err := h.getFlagValues(v)
	if err != nil {
		log.Debug().Err(err).Send()
		fmt.Fprintln(h.w, err.Error())
		os.Exit(1)
	}

	obj, err := h.s.Read(ctx, k, e)
	if err != nil {
		if errors.Is(err, pwdservice.ErrPwdNotExists) {
			log.Debug().Err(err).Int("id", e).Msg("not exists")
			fmt.Fprintln(h.w, err.Error())
			return
		}
		if errors.Is(err, pwdservice.ErrInvalidKey) {
			log.Debug().Err(err).Msg("invalid key")
			fmt.Fprintln(h.w, err.Error())
			os.Exit(1)
		}

		log.Debug().Err(err).Msg("failed to read password")
		handler.InternalError(h.w, err)
		os.Exit(1)
	}

	log.Debug().Int("entry", e).Msg("password readed")
	fmt.Fprintf(
		h.w,
		"the password with entry %d: name=%q, login=%q, password=%q\n",
		e, obj.Name, obj.Login, obj.Password,
	)

}

func (h *PwdReadHandler) getFlagValues(
	v command.ValueGetter,
) (k string, e int, err error) {
	var errs []error
	k, err = getMKey(v)
	if err != nil {
		errs = append(errs, err)
	}

	e, err = v.GetInt(pwdcommand.EntryNumFlag)
	if err != nil {
		errs = append(errs, fmt.Errorf("--%s", pwdcommand.EntryNumFlag))
	}

	err = handler.RequiredFlagsErr(errs)
	return k, e, err
}

type (
	pwdListService interface {
		List(context.Context) ([][2]string, error)
	}

	PwdListHandler struct {
		l logger.Logger
		s pwdListService
		w io.Writer
	}
)

func NewListHandler(
	l logger.Logger, s pwdListService, w io.Writer,
) *PwdListHandler {
	return &PwdListHandler{l, s, w}
}

func (h *PwdListHandler) Handle(ctx context.Context, v command.ValueGetter) {
	const op = "pwdListHandler.Handle"
	log := h.l.With().Str("op", op).Logger()

	entrysNames, err := h.s.List(ctx)
	if err != nil {
		if errors.Is(err, pwdservice.ErrEmptyList) {
			log.Debug().Msg("empty password names list")
			fmt.Fprintln(h.w, err.Error())
			return
		}
		log.Debug().Err(err).Msg("failed to list password names")
		handler.InternalError(h.w, err)
		os.Exit(1)
	}

	log.Debug().Msg("password names list printed")
	h.printNames(entrysNames)

}
func (h *PwdListHandler) printNames(data [][2]string) {
	var out strings.Builder
	for _, v := range data {
		entry, name := v[0], v[1]
		out.WriteString(fmt.Sprintf("\n%s: %s", entry, name))
	}
	fmt.Fprintf(h.w, "saved passwords names:%s\n", out.String())
}

func getMKey(v command.ValueGetter) (string, error) {
	k, err := v.GetString(pwdcommand.MasterKeyFlag)
	if err != nil || len(strings.TrimSpace(k)) == 0 {
		return "", fmt.Errorf("--%s", pwdcommand.MasterKeyFlag)
	}
	return k, nil
}
