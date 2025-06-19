package texthandler

import (
	"fmt"
	"io"

	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/niksmo/gophkeeper/internal/client/command/textcommand"
	"github.com/niksmo/gophkeeper/internal/client/dto"
	"github.com/niksmo/gophkeeper/internal/client/handler"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

const handlerName = "text"

type AddFlags struct {
	Key string
	dto.Text
}

func NewAdd(
	l logger.Logger, s handler.AddService[dto.Text], w io.Writer,
) *handler.AddHandler[AddFlags, dto.Text] {
	h := &handler.AddHandler[AddFlags, dto.Text]{
		Log:     l,
		Service: s,
		Writer:  w,
		Name:    handlerName,
	}

	h.GetFlagsHook = func(v command.ValueGetter) (AddFlags, error) {
		var errs []error
		key, err := handler.GetMasterKeyValue(v)
		if err != nil {
			errs = append(errs, err)
		}

		name, err := handler.GetNameValue(v)
		if err != nil {
			errs = append(errs, err)

		}

		text, err := getTextValue(v)
		if err != nil {
			errs = append(errs, err)
		}

		err = handler.RequiredFlagsErr(errs)
		dto := dto.Text{Name: name, Data: text}
		return AddFlags{key, dto}, err
	}

	h.GetServiceArgsHook = func(
		f AddFlags,
	) (key string, name string, dto dto.Text) {
		return f.Key, f.Name, f.Text
	}

	return h
}

type ReadFlags struct {
	Key      string
	EntryNum int
}

func NewRead(
	l logger.Logger, s handler.ReadService[dto.Text], w io.Writer,
) *handler.ReadHandler[ReadFlags, dto.Text] {
	h := &handler.ReadHandler[ReadFlags, dto.Text]{
		Log:     l,
		Service: s,
		Writer:  w,
		Name:    handlerName,
	}

	h.GetFlagsHook = func(v command.ValueGetter) (ReadFlags, error) {
		var errs []error
		key, err := handler.GetMasterKeyValue(v)
		if err != nil {
			errs = append(errs, err)
		}

		entryNum, err := handler.GetEnryNumValue(v)
		if err != nil {
			errs = append(errs, err)
		}

		err = handler.RequiredFlagsErr(errs)
		return ReadFlags{key, entryNum}, err
	}

	h.GetServiceArgsHook = func(f ReadFlags) (key string, entryNum int) {
		return f.Key, f.EntryNum
	}

	h.GetOutputHook = func(_ ReadFlags, entryNum int, dto dto.Text) string {
		return fmt.Sprintf(
			"the text with entry %d: name=%q\ntext:\n%q\n",
			entryNum, dto.Name, dto.Data,
		)
	}

	return h
}

func NewList(
	l logger.Logger, s handler.ListService, w io.Writer,
) *handler.ListHandler {
	return &handler.ListHandler{
		Log:        l,
		Service:    s,
		Writer:     w,
		Name:       handlerName,
		NamePlural: "texts",
	}
}

type EditFlags struct {
	Key      string
	EntryNum int
	dto.Text
}

func NewEdit(
	l logger.Logger, s handler.EditService[dto.Text], w io.Writer,
) *handler.EditHandler[EditFlags, dto.Text] {
	h := &handler.EditHandler[EditFlags, dto.Text]{
		Log:     l,
		Service: s,
		Writer:  w,
		Name:    handlerName,
	}

	h.GetFlagsHook = func(v command.ValueGetter) (EditFlags, error) {
		var errs []error
		key, err := handler.GetMasterKeyValue(v)
		if err != nil {
			errs = append(errs, err)
		}

		name, err := handler.GetNameValue(v)
		if err != nil {
			errs = append(errs, err)

		}

		entryNum, err := handler.GetEnryNumValue(v)
		if err != nil {
			errs = append(errs, err)
		}

		text, err := getTextValue(v)
		if err != nil {
			errs = append(errs, err)
		}

		err = handler.RequiredFlagsErr(errs)
		dto := dto.Text{Name: name, Data: text}
		return EditFlags{key, entryNum, dto}, err
	}

	h.GetServiceArgsHook = func(
		f EditFlags,
	) (key string, entryNum int, name string, dto dto.Text) {
		return f.Key, f.EntryNum, f.Name, f.Text
	}

	return h
}

type RemoveFlags struct {
	EntryNum int
}

func NewRemove(
	l logger.Logger, s handler.RemoveService, w io.Writer,
) *handler.RemoveHandler[RemoveFlags] {
	h := &handler.RemoveHandler[RemoveFlags]{
		Log:     l,
		Service: s,
		Writer:  w,
		Name:    handlerName,
	}
	h.GetFlagsHook = func(v command.ValueGetter) (RemoveFlags, error) {
		var errs []error
		entryNum, err := handler.GetEnryNumValue(v)
		if err != nil {
			errs = append(errs, err)
		}

		err = handler.RequiredFlagsErr(errs)
		return RemoveFlags{entryNum}, err
	}

	h.GetServiceArgsHook = func(f RemoveFlags) (entryNum int) {
		return f.EntryNum
	}

	return h
}

func getTextValue(v command.ValueGetter) (string, error) {
	t, err := v.GetString(textcommand.TextFlag)
	if err != nil || handler.IsZeroStr(t) {
		return "", fmt.Errorf("--%s", textcommand.TextFlag)
	}
	return t, nil
}
