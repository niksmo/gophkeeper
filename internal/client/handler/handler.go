package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/niksmo/gophkeeper/internal/client/service"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type (
	AddService[T any] interface {
		Add(ctx context.Context, key, name string, dto T) (int, error)
	}

	ReadService[T any] interface {
		Read(ctx context.Context, key string, id int) (T, error)
	}

	ListService interface {
		List(context.Context) ([][2]string, error)
	}

	EditService[T any] interface {
		Edit(
			ctx context.Context,
			key string, entryNum int, name string, obj T,
		) error
	}

	RemoveService interface {
		Remove(ctx context.Context, entryNum int) error
	}
)

// AddHandler
type AddHandler[F any, O any] struct {
	Log                logger.Logger
	Service            AddService[O]
	Writer             io.Writer
	Name               string
	GetFlagsHook       func(command.ValueGetter) (F, error)
	GetServiceArgsHook func(F) (key, name string, dto O)
}

func (h *AddHandler[F, O]) Handle(ctx context.Context, v command.ValueGetter) {
	const op = "AddHandler.Handle"

	log := h.Log.WithOp(op)

	flags, err := h.GetFlagsHook(v)
	if err != nil {
		log.Debug().Err(err).Send()
		fmt.Fprintln(h.Writer, err.Error())
		os.Exit(1)
	}

	key, name, dto := h.GetServiceArgsHook(flags)

	entryNum, err := h.Service.Add(ctx, key, name, dto)
	if err != nil {
		switch {
		case handleAlreadyExists(log, h.Writer, h.Name, name, err):
			return
		default:
			log.Debug().Err(err).Msg(fmt.Sprintf("failed to add %s", h.Name))
			InternalError(h.Writer, err)
			os.Exit(1)
		}
	}

	h.printOut(entryNum)
}

func (h *AddHandler[F, O]) printOut(entryNum int) {
	fmt.Fprintf(
		h.Writer,
		"the %s is saved under the record number %d\n",
		h.Name, entryNum,
	)
}

// ReadHandler
type ReadHandler[F any, O any] struct {
	Log                logger.Logger
	Service            ReadService[O]
	Writer             io.Writer
	Name               string
	GetFlagsHook       func(command.ValueGetter) (F, error)
	GetServiceArgsHook func(F) (key string, entryNum int)
	GetOutputHook      func(flags F, entryNum int, dto O) string
}

func (h *ReadHandler[F, O]) Handle(ctx context.Context, v command.ValueGetter) {
	const op = "ReadHandler.Handle"
	log := h.Log.WithOp(op)

	flags, err := h.GetFlagsHook(v)
	if err != nil {
		log.Debug().Err(err).Send()
		fmt.Fprintln(h.Writer, err.Error())
		os.Exit(1)
	}

	key, entryNum := h.GetServiceArgsHook(flags)

	dto, err := h.Service.Read(ctx, key, entryNum)
	if err != nil {
		switch {
		case handleNotExists(log, h.Writer, h.Name, entryNum, err):
			return
		case handleInvalidKey(log, h.Writer, err):
			os.Exit(1)
		default:
			log.Debug().Err(err).Msg(fmt.Sprintf("failed to read %s", h.Name))
			InternalError(h.Writer, err)
			os.Exit(1)
		}
	}

	h.printOut(flags, entryNum, dto)
}

func (h *ReadHandler[F, O]) printOut(flags F, entryNum int, obj O) {
	fmt.Fprint(h.Writer, h.GetOutputHook(flags, entryNum, obj))
}

// ListHandler
type ListHandler struct {
	Log              logger.Logger
	Service          ListService
	Writer           io.Writer
	Name, NamePlural string
}

func (h *ListHandler) Handle(ctx context.Context, v command.ValueGetter) {
	const op = "ListHandler.Handle"
	log := h.Log.With().Str("op", op).Logger()

	idNamePairs, err := h.Service.List(ctx)
	if err != nil {
		log.Debug().Str("handler", h.Name).Err(err).Msg("failed to list names")
		InternalError(h.Writer, err)
		os.Exit(1)
	}
	h.printOut(idNamePairs)
}

func (h *ListHandler) printOut(data [][2]string) {
	if len(data) == 0 {
		fmt.Fprintf(h.Writer, "there are no saved %s\n", h.NamePlural)
		return
	}

	var b strings.Builder
	for _, v := range data {
		entry, name := v[0], v[1]
		b.WriteString(fmt.Sprintf("\n%s: %s", entry, name))
	}
	fmt.Fprintf(h.Writer, "saved %s names:%s\n", h.NamePlural, b.String())
}

// EditHandler
type EditHandler[F any, O any] struct {
	Log                logger.Logger
	Service            EditService[O]
	Writer             io.Writer
	Name               string
	GetFlagsHook       func(command.ValueGetter) (F, error)
	GetServiceArgsHook func(F) (key string, entryNum int, name string, dto O)
}

func (h *EditHandler[F, O]) Handle(ctx context.Context, v command.ValueGetter) {
	const op = "EditHandler.Handle"
	log := h.Log.WithOp(op)

	flags, err := h.GetFlagsHook(v)
	if err != nil {
		log.Debug().Err(err).Send()
		fmt.Fprintln(h.Writer, err.Error())
		os.Exit(1)
	}

	key, entryNum, name, dto := h.GetServiceArgsHook(flags)

	err = h.Service.Edit(ctx, key, entryNum, name, dto)
	if err != nil {
		switch {
		case handleAlreadyExists(log, h.Writer, h.Name, name, err):
			return
		case handleNotExists(log, h.Writer, h.Name, entryNum, err):
			return
		default:
			log.Debug().Err(err).Msg(fmt.Sprintf("failed to edit %s", h.Name))
			InternalError(h.Writer, err)
			os.Exit(1)
		}
	}

	h.printOut(entryNum)
}

func (h *EditHandler[F, O]) printOut(entryNum int) {
	fmt.Fprintf(
		h.Writer,
		"the %s under the record number %d was edited\n",
		h.Name, entryNum,
	)
}

// RemoveHandler
type RemoveHandler[F any] struct {
	Log                logger.Logger
	Service            RemoveService
	Writer             io.Writer
	Name               string
	GetFlagsHook       func(command.ValueGetter) (F, error)
	GetServiceArgsHook func(F) (entryNum int)
}

func (h *RemoveHandler[F]) Handle(ctx context.Context, v command.ValueGetter) {
	const op = "RemoveHandler.Handle"
	log := h.Log.WithOp(op)

	flags, err := h.GetFlagsHook(v)
	if err != nil {
		log.Debug().Err(err).Send()
		fmt.Fprintln(h.Writer, err.Error())
		os.Exit(1)
	}

	entryNum := h.GetServiceArgsHook(flags)
	err = h.Service.Remove(ctx, entryNum)
	if err != nil {
		switch {
		case handleNotExists(log, h.Writer, h.Name, entryNum, err):
			return
		default:
			log.Debug().Err(err).Msg(fmt.Sprintf("failed to remove %s", h.Name))
			InternalError(h.Writer, err)
			os.Exit(1)
		}
	}

	h.printOut(entryNum)
}

func (h *RemoveHandler[F]) printOut(entryNum int) {
	fmt.Fprintf(
		h.Writer,
		"the %s under the record number %d was removed\n",
		h.Name, entryNum,
	)
}

// *Helpers*

func RequiredFlagsErr(errs []error) error {
	if len(errs) != 0 {
		return fmt.Errorf(
			"required flags are not specified:\n%w",
			errors.Join(errs...),
		)
	}
	return nil
}

func GetMasterKeyValue(v command.ValueGetter) (string, error) {
	k, err := v.GetString(command.MasterKeyFlag)
	if err != nil || IsZeroStr(k) {
		return "", fmt.Errorf("--%s", command.MasterKeyFlag)
	}
	return k, nil
}

func GetNameValue(v command.ValueGetter) (string, error) {
	n, err := v.GetString(command.NameFlag)
	if err != nil || IsZeroStr(n) {
		return "", fmt.Errorf("--%s", command.NameFlag)
	}
	return n, nil
}

func GetEnryNumValue(v command.ValueGetter) (int, error) {
	e, err := v.GetInt(command.EntryNumFlag)
	if err != nil {
		return 0, fmt.Errorf("--%s", command.EntryNumFlag)
	}
	return e, nil
}

func IsZeroStr(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func InternalError(w io.Writer, err error) {
	fmt.Fprintf(
		w,
		"the application completed with an error: %s\n",
		err.Error(),
	)
}

func handleAlreadyExists(
	l logger.Logger, w io.Writer, hName, oName string, err error,
) bool {
	if !errors.Is(err, service.ErrAlreadyExists) {
		return false
	}

	l.Debug().Err(err).Msg("object already exists")
	fmt.Fprintf(
		w, "%s with name '%s' already exists\n",
		hName, oName,
	)
	return true
}

func handleNotExists(
	l logger.Logger, w io.Writer, hName string, entryNum int, err error,
) bool {
	if !errors.Is(err, service.ErrNotExists) {
		return false
	}

	l.Debug().Err(err).Int("id", entryNum).Msg("not exists")
	fmt.Fprintf(
		w,
		"the %s with entry number %d is not exists\n",
		hName, entryNum,
	)
	return true
}

func handleInvalidKey(l logger.Logger, w io.Writer, err error) bool {
	if !errors.Is(err, service.ErrInvalidKey) {
		return false
	}
	l.Debug().Err(err).Msg("invalid key")
	fmt.Fprintln(w, err.Error())
	return true
}
