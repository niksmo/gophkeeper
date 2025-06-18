package pwdhandler

import (
	"fmt"
	"io"

	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/niksmo/gophkeeper/internal/client/command/pwdcommand"
	"github.com/niksmo/gophkeeper/internal/client/dto"
	"github.com/niksmo/gophkeeper/internal/client/handler"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type AddFlags struct {
	Key string
	dto.PWD
}

func NewAddHandler(
	l logger.Logger, s handler.AddService[dto.PWD], w io.Writer,
) *handler.AddHandler[AddFlags, dto.PWD] {
	h := &handler.AddHandler[AddFlags, dto.PWD]{
		Log: l, Service: s, Writer: w, Name: "password",
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

		password, err := getPasswordValue(v)
		if err != nil {
			errs = append(errs, err)
		}

		login := getLoginValue(v)

		err = handler.RequiredFlagsErr(errs)
		dto := dto.PWD{Name: name, Login: login, Password: password}
		return AddFlags{key, dto}, err
	}

	h.GetServiceArgsHook = func(f AddFlags) (key string, name string, dto dto.PWD) {
		return f.Key, f.Name, f.PWD
	}

	return h
}

type ReadFlags struct {
	Key      string
	EntryNum int
}

func NewReadHandler(
	l logger.Logger, s handler.ReadService[dto.PWD], w io.Writer,
) *handler.ReadHandler[ReadFlags, dto.PWD] {
	h := &handler.ReadHandler[ReadFlags, dto.PWD]{
		Log:     l,
		Service: s,
		Writer:  w,
		Name:    "password",
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

	h.GetOutStr = func(entryNum int, dto dto.PWD) string {
		return fmt.Sprintf(
			"the password with entry %d: name=%q, login=%q, password=%q",
			entryNum, dto.Name, dto.Login, dto.Password,
		)
	}

	return h
}

func NewListHandler(
	l logger.Logger, s handler.ListService, w io.Writer,
) *handler.ListHandler {
	return &handler.ListHandler{
		Log:        l,
		Service:    s,
		Writer:     w,
		Name:       "password",
		NamePlural: "passwords",
	}
}

type EditFlags struct {
	Key      string
	EntryNum int
	dto.PWD
}

func NewEditHandler(
	l logger.Logger, s handler.EditService[dto.PWD], w io.Writer,
) *handler.EditHandler[EditFlags, dto.PWD] {
	h := &handler.EditHandler[EditFlags, dto.PWD]{
		Log:     l,
		Service: s,
		Writer:  w,
		Name:    "password",
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

		password, err := getPasswordValue(v)
		if err != nil {
			errs = append(errs, err)
		}

		login := getLoginValue(v)

		err = handler.RequiredFlagsErr(errs)
		dto := dto.PWD{
			Name:     name,
			Login:    login,
			Password: password,
		}
		return EditFlags{key, entryNum, dto}, err
	}

	h.GetServiceArgsHook = func(
		f EditFlags,
	) (key string, entryNum int, name string, dto dto.PWD) {
		return f.Key, f.EntryNum, f.Name, f.PWD
	}

	return h
}

type RemoveFlags struct {
	EntryNum int
}

func NewRemoveHandler(
	l logger.Logger, s handler.RemoveService, w io.Writer,
) *handler.RemoveHandler[RemoveFlags] {
	h := &handler.RemoveHandler[RemoveFlags]{
		Log:     l,
		Service: s,
		Writer:  w,
		Name:    "password",
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

func getPasswordValue(v command.ValueGetter) (string, error) {
	p, err := v.GetString(pwdcommand.PasswordFlag)
	if err != nil || handler.IsZeroStr(p) {
		return "", fmt.Errorf("--%s", pwdcommand.PasswordFlag)
	}
	return p, nil
}

func getLoginValue(v command.ValueGetter) string {
	l, _ := v.GetString(pwdcommand.LoginFlag)
	return l
}
