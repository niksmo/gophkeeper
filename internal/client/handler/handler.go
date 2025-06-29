package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

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

type ListCmdHandler struct {
	Log              logger.Logger
	Service          ListService
	Writer           io.Writer
	Name, NamePlural string
}

func (h *ListCmdHandler) Handle(ctx context.Context) {
	const op = "ListCmdHandler.Handle"

	log := h.Log.WithOp(op)

	idNamePairs, err := h.Service.List(ctx)
	if err != nil {
		HandleUnexpectedErr(err, log, h.Writer)
	}

	h.printOutput(idNamePairs)
}

func (h *ListCmdHandler) printOutput(data [][2]string) {
	if len(data) == 0 {
		fmt.Fprintf(h.Writer, "there are no saved %s\n", h.NamePlural)
		return
	}

	list := h.buildList(data)

	fmt.Fprintf(h.Writer, "saved %s names:%s\n", h.NamePlural, list)
}

func (h *ListCmdHandler) buildList(data [][2]string) string {
	var b strings.Builder
	for _, v := range data {
		entry, name := v[0], v[1]
		b.WriteString(fmt.Sprintf("\n%s: %s", entry, name))
	}
	return b.String()
}

type RemoveCmdHandler struct {
	Log     logger.Logger
	Service RemoveService
	Writer  io.Writer
	Name    string
}

func (h *RemoveCmdHandler) Handle(ctx context.Context, entryNum int) {
	const op = "RemoveCmdHandler.Handle"

	log := h.Log.WithOp(op)

	err := h.Service.Remove(ctx, entryNum)
	if err != nil {
		HandleNotExistsErr(err, log, h.Writer, h.Name, entryNum)
		HandleUnexpectedErr(err, log, h.Writer)
	}

	h.printOutput(entryNum)
}

func (h *RemoveCmdHandler) printOutput(entryNum int) {
	fmt.Fprintf(h.Writer,
		"the %s under the record number %d was removed\n", h.Name, entryNum)
}

func PrintSaveEntryOutput(w io.Writer, entity string, entryNum int) {
	fmt.Fprintf(w,
		"the %s is saved under the record number %d\n", entity, entryNum)

}

func HandleAlreadyExistsErr(
	err error, l logger.Logger, w io.Writer, entity, objName string,
) {
	if !errors.Is(err, service.ErrAlreadyExists) {
		return
	}

	l.Debug().Err(err).Msg("object already exists")

	fmt.Fprintf(w,
		"%s with name '%s' already exists\n", entity, objName)
	os.Exit(1)
}

func HandleNotExistsErr(
	err error, l logger.Logger, w io.Writer, entity string, entryNum int,
) {
	if !errors.Is(err, service.ErrNotExists) {
		return
	}

	l.Debug().Err(err).Int("id", entryNum).Msg("not exists")

	fmt.Fprintf(w,
		"the %s with entry number %d is not exists\n", entity, entryNum)
	os.Exit(1)
}

func HandleInvalidKeyErr(err error, l logger.Logger, w io.Writer) {
	if !errors.Is(err, service.ErrInvalidKey) {
		return
	}

	l.Debug().Err(err).Msg("invalid key")

	fmt.Fprintln(w, "invalid key provided")
	os.Exit(1)
}

func HandleUnexpectedErr(err error, log logger.Logger, w io.Writer) {
	if err == nil {
		return
	}

	log.Debug().Err(err).Msg("unexpected error")

	fmt.Fprintf(w,
		"the application completed with an error: %s\n", err.Error())
	os.Exit(1)
}
