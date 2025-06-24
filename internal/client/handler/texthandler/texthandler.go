package texthandler

import (
	"context"
	"fmt"
	"io"

	"github.com/niksmo/gophkeeper/internal/client/command/textcommand"
	"github.com/niksmo/gophkeeper/internal/client/dto"
	"github.com/niksmo/gophkeeper/internal/client/handler"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

const entity = "text"

type AddCmdHandler struct {
	l logger.Logger
	s handler.AddService[dto.Text]
	w io.Writer
}

func NewAdd(
	l logger.Logger, s handler.AddService[dto.Text], w io.Writer,
) *AddCmdHandler {
	return &AddCmdHandler{l, s, w}
}

func (h *AddCmdHandler) Handle(
	ctx context.Context, fv textcommand.AddCmdFlags,
) {
	const op = "texthandlerAdd.Handle"

	log := h.l.WithOp(op)

	o := dto.Text{Name: fv.Name, Data: fv.Text}
	entryNum, err := h.s.Add(ctx, fv.Key, fv.Name, o)
	if err != nil {
		handler.HandleAlreadyExistsErr(err, log, h.w, entity, fv.Name)
		handler.HandleUnexpectedErr(err, log, h.w)
	}

	handler.PrintSaveEntryOutput(h.w, entity, entryNum)
}

type EditCmdHandler struct {
	l logger.Logger
	s handler.EditService[dto.Text]
	w io.Writer
}

func NewEdit(
	l logger.Logger, s handler.EditService[dto.Text], w io.Writer,
) *EditCmdHandler {
	return &EditCmdHandler{l, s, w}
}

func (h *EditCmdHandler) Handle(
	ctx context.Context, fv textcommand.EditCmdFlags,
) {
	const op = "texthandlerEdit.Handle"

	log := h.l.WithOp(op)

	o := dto.Text{Name: fv.Name, Data: fv.Text}
	err := h.s.Edit(ctx, fv.Key, fv.EntryNum, fv.Name, o)
	if err != nil {
		handler.HandleAlreadyExistsErr(err, log, h.w, entity, fv.Name)
		handler.HandleUnexpectedErr(err, log, h.w)
	}

	handler.PrintSaveEntryOutput(h.w, entity, fv.EntryNum)
}

type ReadCmdHandler struct {
	l logger.Logger
	s handler.ReadService[dto.Text]
	w io.Writer
}

func NewRead(
	l logger.Logger, s handler.ReadService[dto.Text], w io.Writer,
) *ReadCmdHandler {
	return &ReadCmdHandler{l, s, w}
}

func (h *ReadCmdHandler) Handle(
	ctx context.Context, key string, entryNum int,
) {
	const op = "texthandlerRead.Handle"

	log := h.l.WithOp(op)

	obj, err := h.s.Read(ctx, key, entryNum)
	if err != nil {
		handler.HandleInvalidKeyErr(err, log, h.w)
		handler.HandleNotExistsErr(err, log, h.w, entity, entryNum)
		handler.HandleUnexpectedErr(err, log, h.w)
	}
	h.printOutput(entryNum, obj)
}

func (h *ReadCmdHandler) printOutput(entryNum int, o dto.Text) {
	fmt.Fprintf(h.w,
		"the text with entry %d: name=%q\ntext:\n%q\n",
		entryNum, o.Name, o.Data)
}

func NewList(
	l logger.Logger, s handler.ListService, w io.Writer,
) *handler.ListCmdHandler {
	return &handler.ListCmdHandler{
		Log:        l,
		Service:    s,
		Writer:     w,
		Name:       entity,
		NamePlural: entity + "s",
	}
}

func NewRemove(
	l logger.Logger, s handler.RemoveService, w io.Writer,
) *handler.RemoveCmdHandler {
	return &handler.RemoveCmdHandler{
		Log:     l,
		Service: s,
		Writer:  w,
		Name:    entity,
	}
}
