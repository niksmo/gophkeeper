package binhandler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/niksmo/gophkeeper/internal/client/command/bincommand"
	"github.com/niksmo/gophkeeper/internal/client/dto"
	"github.com/niksmo/gophkeeper/internal/client/handler"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

const entity = "binary"

type AddCmdHandler struct {
	l logger.Logger
	s handler.AddService[dto.BIN]
	w io.Writer
}

func NewAdd(
	l logger.Logger, s handler.AddService[dto.BIN], w io.Writer,
) *AddCmdHandler {
	return &AddCmdHandler{l, s, w}
}

func (h *AddCmdHandler) Handle(ctx context.Context, fv bincommand.AddCmdFlags) {
	const op = "binhandlerAdd.Handle"

	log := h.l.WithOp(op)

	o, err := getBinDto(fv.Name, fv.Filepath)
	if err != nil {
		handler.HandleUnexpectedErr(err, log, h.w)
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
	s handler.EditService[dto.BIN]
	w io.Writer
}

func NewEdit(
	l logger.Logger, s handler.EditService[dto.BIN], w io.Writer,
) *EditCmdHandler {
	return &EditCmdHandler{l, s, w}
}

func (h *EditCmdHandler) Handle(
	ctx context.Context, fv bincommand.EditCmdFlags,
) {
	const op = "binhandlerEdit.Handle"

	log := h.l.WithOp(op)

	o, err := getBinDto(fv.Name, fv.Filepath)
	if err != nil {
		handler.HandleUnexpectedErr(err, log, h.w)
	}

	err = h.s.Edit(ctx, fv.Key, fv.EntryNum, fv.Name, o)
	if err != nil {
		handler.HandleAlreadyExistsErr(err, log, h.w, entity, fv.Name)
		handler.HandleNotExistsErr(err, log, h.w, entity, fv.EntryNum)
		handler.HandleUnexpectedErr(err, log, h.w)
	}

	handler.PrintSaveEntryOutput(h.w, entity, fv.EntryNum)
}

type ReadCmdHandler struct {
	l logger.Logger
	s handler.ReadService[dto.BIN]
	w io.Writer
}

func NewRead(
	l logger.Logger, s handler.ReadService[dto.BIN], w io.Writer,
) *ReadCmdHandler {
	return &ReadCmdHandler{l, s, w}
}

func (h *ReadCmdHandler) Handle(
	ctx context.Context, fv bincommand.ReadCmdFlags,
) {
	const op = "binhandlerRead.Handle"

	log := h.l.WithOp(op)

	obj, err := h.s.Read(ctx, fv.Key, fv.EntryNum)
	if err != nil {
		handler.HandleInvalidKeyErr(err, log, h.w)
		handler.HandleNotExistsErr(err, log, h.w, entity, fv.EntryNum)
		handler.HandleUnexpectedErr(err, log, h.w)
	}

	output := h.buildOutput(fv.EntryNum, fv.Filepath, obj)
	h.printOutput(output)
}

func (h *ReadCmdHandler) buildOutput(
	entryNum int, filepath string, o dto.BIN,
) string {
	var b strings.Builder
	b.WriteString(
		fmt.Sprintf(
			"the binary data with entry %d: name=%q size=%d ext=%q \n",
			entryNum, o.Name, len(o.Data), o.Ext,
		))

	if h.writeToFile(filepath) {
		err := h.writeData(filepath, o.Data)
		if err != nil {
			b.WriteString(err.Error() + "\n")
		} else {
			b.WriteString(fmt.Sprintf("saved to filepath: %s\n", filepath))
		}
	}

	return b.String()
}

func (h *ReadCmdHandler) writeData(filepath string, data []byte) error {
	err := writeData(filepath, data)
	if err != nil {
		return fmt.Errorf("write file error: %w", err)
	}
	return nil
}

func (h *ReadCmdHandler) writeToFile(filepath string) bool {
	return filepath != ""
}

func (h *ReadCmdHandler) printOutput(out string) {
	fmt.Fprint(h.w, out)
}

func NewList(
	l logger.Logger, s handler.ListService, w io.Writer,
) *handler.ListCmdHandler {
	return &handler.ListCmdHandler{
		Log:        l,
		Service:    s,
		Writer:     w,
		Name:       entity,
		NamePlural: "binaries",
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

func getBinDto(name, path string) (dto.BIN, error) {
	data, ext, err := getFileData(path)
	if err != nil {
		return dto.BIN{}, fmt.Errorf("file error: %w", err)
	}
	return dto.BIN{Name: name, Data: data, Ext: ext}, nil
}

func getFileData(path string) ([]byte, string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, "", err
	}

	if err := verifyReadingFile(path); err != nil {
		return nil, "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}
	ext := filepath.Ext(path)
	return data, ext, nil
}

func verifyReadingFile(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return errors.New("filepath is directory")
	}

	if stat.Size() > 1024*1024*100 {
		return errors.New("the file size should be less or equal 100Mb")
	}
	return nil
}

func writeData(path string, data []byte) error {
	if path == "" {
		return nil
	}

	if fileExists(path) {
		return fmt.Errorf("the file is exists")
	}

	return os.WriteFile(path, data, 0o644)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
