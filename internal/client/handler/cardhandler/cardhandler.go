package cardhandler

import (
	"context"
	"fmt"
	"io"

	"github.com/niksmo/gophkeeper/internal/client/command/cardcommand"
	"github.com/niksmo/gophkeeper/internal/client/dto"
	"github.com/niksmo/gophkeeper/internal/client/handler"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

const entity = "bank card"

type AddCmdHandler struct {
	l logger.Logger
	s handler.AddService[dto.BankCard]
	w io.Writer
}

func NewAdd(
	l logger.Logger, s handler.AddService[dto.BankCard], w io.Writer,
) *AddCmdHandler {
	return &AddCmdHandler{l, s, w}
}

func (h *AddCmdHandler) Handle(
	ctx context.Context, fv cardcommand.AddCmdFlags,
) {
	const op = "cardhandlerAdd.Handle"

	log := h.l.WithOp(op)

	o := dto.BankCard{
		Name:       fv.Name,
		Number:     fv.CardNum,
		ExpDate:    fv.Exp,
		HolderName: fv.Holder,
	}
	entryNum, err := h.s.Add(ctx, fv.Key, fv.Name, o)
	if err != nil {
		handler.HandleAlreadyExistsErr(err, log, h.w, entity, fv.Name)
		handler.HandleUnexpectedErr(err, log, h.w)
	}

	handler.PrintSaveEntryOutput(h.w, entity, entryNum)
}

type EditCmdHandler struct {
	l logger.Logger
	s handler.EditService[dto.BankCard]
	w io.Writer
}

func NewEdit(
	l logger.Logger, s handler.EditService[dto.BankCard], w io.Writer,
) *EditCmdHandler {
	return &EditCmdHandler{l, s, w}
}

func (h *EditCmdHandler) Handle(
	ctx context.Context, fv cardcommand.EditCmdFlags,
) {
	const op = "cardhandlerEdit.Handle"

	log := h.l.WithOp(op)

	o := dto.BankCard{
		Name:       fv.Name,
		Number:     fv.CardNum,
		ExpDate:    fv.Exp,
		HolderName: fv.Holder,
	}
	err := h.s.Edit(ctx, fv.Key, fv.EntryNum, fv.Name, o)
	if err != nil {
		handler.HandleAlreadyExistsErr(err, log, h.w, entity, fv.Name)
		handler.HandleUnexpectedErr(err, log, h.w)
	}

	handler.PrintSaveEntryOutput(h.w, entity, fv.EntryNum)
}

type ReadCmdHandler struct {
	l logger.Logger
	s handler.ReadService[dto.BankCard]
	w io.Writer
}

func NewRead(
	l logger.Logger, s handler.ReadService[dto.BankCard], w io.Writer,
) *ReadCmdHandler {
	return &ReadCmdHandler{l, s, w}
}

func (h *ReadCmdHandler) Handle(
	ctx context.Context, key string, entryNum int,
) {
	const op = "cardhandlerRead.Handle"

	log := h.l.WithOp(op)

	obj, err := h.s.Read(ctx, key, entryNum)
	if err != nil {
		handler.HandleInvalidKeyErr(err, log, h.w)
		handler.HandleNotExistsErr(err, log, h.w, entity, entryNum)
		handler.HandleUnexpectedErr(err, log, h.w)
	}
	h.printOutput(entryNum, obj)
}

func (h *ReadCmdHandler) printOutput(entryNum int, o dto.BankCard) {
	fmt.Fprintf(h.w,
		"the bank card with entry %d: name=%q number=%q exp=%q holder=%q\n",
		entryNum, o.Name, o.Number, o.ExpDate, o.HolderName)
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
