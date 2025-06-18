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
	"github.com/niksmo/gophkeeper/internal/client/dto"
	"github.com/niksmo/gophkeeper/internal/client/handler"
	"github.com/niksmo/gophkeeper/internal/client/service"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type (
	addService interface {
		Add(ctx context.Context, key, name string, obj dto.PWD) (int, error)
	}

	readService interface {
		Read(ctx context.Context, key string, id int) (dto.PWD, error)
	}

	listService interface {
		List(context.Context) ([][2]string, error)
	}

	editService interface {
		Edit(
			ctx context.Context,
			key string,
			entryNum int,
			name string,
			obj dto.PWD,
		) error
	}

	removeService interface {
		Remove(ctx context.Context, entryNum int) error
	}
)

type PwdAddHandler struct {
	l logger.Logger
	s addService
	w io.Writer
}

func NewAddHandler(
	l logger.Logger, s addService, w io.Writer,
) *PwdAddHandler {
	return &PwdAddHandler{l, s, w}
}

func (h *PwdAddHandler) Handle(ctx context.Context, v command.ValueGetter) {
	const op = "PwdAddHandler.Handle"
	log := h.l.With().Str("op", op).Logger()

	key, name, password, login, err := h.getFlagValues(v)
	if err != nil {
		log.Debug().Err(err).Send()
		fmt.Fprintln(h.w, err.Error())
		os.Exit(1)
	}

	obj := dto.PWD{
		Name:     name,
		Login:    login,
		Password: password,
	}

	entryNum, err := h.s.Add(ctx, key, name, obj)
	if err != nil {
		log.Debug().Err(err).Msg("failed to save password")
		handler.InternalError(h.w, err)
		os.Exit(1)
	}

	h.printOut(entryNum)
}

func (h *PwdAddHandler) getFlagValues(
	v command.ValueGetter,
) (k, n, p, l string, err error) {
	var errs []error
	k, err = handler.GetMasterKeyValue(v)
	if err != nil {
		errs = append(errs, err)
	}

	n, err = handler.GetNameValue(v)
	if err != nil {
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

func (h *PwdAddHandler) printOut(entryNum int) {
	fmt.Fprintf(
		h.w,
		"the password is saved under the record number %d\n",
		entryNum,
	)
}

type PwdReadHandler struct {
	l logger.Logger
	s readService
	w io.Writer
}

func NewReadHandler(
	l logger.Logger, s readService, w io.Writer,
) *PwdReadHandler {
	return &PwdReadHandler{l, s, w}
}

func (h *PwdReadHandler) Handle(ctx context.Context, v command.ValueGetter) {
	const op = "PwdReadHandler.Handle"
	log := h.l.With().Str("op", op).Logger()

	key, entryNum, err := h.getFlagValues(v)
	if err != nil {
		log.Debug().Err(err).Send()
		fmt.Fprintln(h.w, err.Error())
		os.Exit(1)
	}

	obj, err := h.s.Read(ctx, key, entryNum)
	if err != nil {
		if errors.Is(err, service.ErrNotExists) {
			log.Debug().Err(err).Int("id", entryNum).Msg("not exists")
			fmt.Fprintf(
				h.w,
				"the password with entry number %d is not exists\n",
				entryNum,
			)
			return
		}
		if errors.Is(err, service.ErrInvalidKey) {
			log.Debug().Err(err).Msg("invalid key")
			fmt.Fprintln(h.w, err.Error())
			os.Exit(1)
		}

		log.Debug().Err(err).Msg("failed to read password")
		handler.InternalError(h.w, err)
		os.Exit(1)
	}

	h.printOut(entryNum, obj)
}

func (h *PwdReadHandler) getFlagValues(
	v command.ValueGetter,
) (k string, e int, err error) {
	var errs []error
	k, err = handler.GetMasterKeyValue(v)
	if err != nil {
		errs = append(errs, err)
	}

	e, err = handler.GetEnryNumValue(v)
	if err != nil {
		errs = append(errs, err)
	}

	err = handler.RequiredFlagsErr(errs)
	return k, e, err
}

func (h *PwdReadHandler) printOut(entryNum int, obj dto.PWD) {
	fmt.Fprintf(
		h.w,
		"the password with entry %d: name=%q, login=%q, password=%q\n",
		entryNum, obj.Name, obj.Login, obj.Password,
	)
}

type PwdListHandler struct {
	l logger.Logger
	s listService
	w io.Writer
}

func NewListHandler(
	l logger.Logger, s listService, w io.Writer,
) *PwdListHandler {
	return &PwdListHandler{l, s, w}
}

func (h *PwdListHandler) Handle(ctx context.Context, v command.ValueGetter) {
	const op = "PwdListHandler.Handle"
	log := h.l.With().Str("op", op).Logger()

	idNamePairs, err := h.s.List(ctx)
	if err != nil {
		log.Debug().Err(err).Msg("failed to list password names")
		handler.InternalError(h.w, err)
		os.Exit(1)
	}
	h.printOut(idNamePairs)

}
func (h *PwdListHandler) printOut(data [][2]string) {
	if len(data) == 0 {
		fmt.Fprintln(h.w, "there are no saved passwords")
		return
	}

	var out strings.Builder
	for _, v := range data {
		entry, name := v[0], v[1]
		out.WriteString(fmt.Sprintf("\n%s: %s", entry, name))
	}
	fmt.Fprintf(h.w, "saved passwords names:%s\n", out.String())
}

type PwdEditHandler struct {
	l logger.Logger
	s editService
	w io.Writer
}

func NewEditHandler(
	l logger.Logger, s editService, w io.Writer,
) *PwdEditHandler {
	return &PwdEditHandler{l, s, w}
}

func (h *PwdEditHandler) Handle(ctx context.Context, v command.ValueGetter) {
	const op = "PwdEditHandler.Handle"
	log := h.l.With().Str("op", op).Logger()

	key, name, password, login, entryNum, err := h.getFlagValues(v)
	if err != nil {
		log.Debug().Err(err).Send()
		fmt.Fprintln(h.w, err.Error())
		os.Exit(1)
	}

	obj := dto.PWD{
		Name:     name,
		Login:    login,
		Password: password,
	}

	err = h.s.Edit(ctx, key, entryNum, name, obj)
	if err != nil {
		if errors.Is(err, service.ErrNotExists) {
			log.Debug().Err(err).Int("id", entryNum).Msg("not exists")
			fmt.Fprintf(
				h.w,
				"the password with entry number %d is not exists\n",
				entryNum,
			)
			return
		}
		if errors.Is(err, service.ErrInvalidKey) {
			log.Debug().Err(err).Msg("invalid key")
			fmt.Fprintln(h.w, err.Error())
			os.Exit(1)
		}
		log.Debug().Err(err).Msg("failed to edit password")
		handler.InternalError(h.w, err)
		os.Exit(1)
	}

	h.printOut(entryNum)
}

func (h *PwdEditHandler) getFlagValues(
	v command.ValueGetter,
) (k, n, p, l string, e int, err error) {
	var errs []error
	k, err = handler.GetMasterKeyValue(v)
	if err != nil {
		errs = append(errs, err)
	}

	n, err = handler.GetNameValue(v)
	if err != nil {
		errs = append(errs, err)

	}

	e, err = handler.GetEnryNumValue(v)
	if err != nil {
		errs = append(errs, err)
	}

	p, err = v.GetString(pwdcommand.PasswordFlag)
	if err != nil || len(strings.TrimSpace(p)) == 0 {
		errs = append(errs, fmt.Errorf("--%s", pwdcommand.PasswordFlag))
	}

	l, _ = v.GetString(pwdcommand.LoginFlag)

	err = handler.RequiredFlagsErr(errs)

	return k, n, p, l, e, err
}

func (h *PwdEditHandler) printOut(entryNum int) {
	fmt.Fprintf(
		h.w,
		"the password under the record number %d was edited\n",
		entryNum,
	)
}

type PwdRemoveHandler struct {
	l logger.Logger
	s removeService
	w io.Writer
}

func NewRemoveHandler(
	l logger.Logger, s removeService, w io.Writer,
) *PwdRemoveHandler {
	return &PwdRemoveHandler{l, s, w}
}

func (h *PwdRemoveHandler) Handle(ctx context.Context, v command.ValueGetter) {
	const op = "PwdRemoveHandler.Handle"
	log := h.l.With().Str("op", op).Logger()

	entryNum, err := h.getFlagValues(v)
	if err != nil {
		log.Debug().Err(err).Send()
		fmt.Fprintln(h.w, err.Error())
		os.Exit(1)
	}

	err = h.s.Remove(ctx, entryNum)
	if err != nil {
		if errors.Is(err, service.ErrNotExists) {
			log.Debug().Err(err).Int("id", entryNum).Msg("not exists")
			fmt.Fprintf(
				h.w,
				"the password with entry number %d is not exists\n",
				entryNum,
			)
			return
		}
		log.Debug().Err(err).Msg("failed to remove password")
		handler.InternalError(h.w, err)
		os.Exit(1)
	}

	h.printOut(entryNum)
}

func (h *PwdRemoveHandler) getFlagValues(
	v command.ValueGetter,
) (e int, err error) {
	var errs []error
	e, err = handler.GetEnryNumValue(v)
	if err != nil {
		errs = append(errs, err)
	}

	err = handler.RequiredFlagsErr(errs)
	return e, err
}

func (h *PwdRemoveHandler) printOut(entryNum int) {
	fmt.Fprintf(
		h.w,
		"the password under the record number %d was removed\n",
		entryNum,
	)
}
