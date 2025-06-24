package cardhandler

import (
	"fmt"
	"io"

	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/niksmo/gophkeeper/internal/client/command/cardcommand"
	"github.com/niksmo/gophkeeper/internal/client/dto"
	"github.com/niksmo/gophkeeper/internal/client/handler"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

const handlerName = "bank card"

type AddFlags struct {
	Key string
	dto.BankCard
}

func NewAdd(
	l logger.Logger, s handler.AddService[dto.BankCard], w io.Writer,
) *handler.AddHandler[AddFlags, dto.BankCard] {
	h := &handler.AddHandler[AddFlags, dto.BankCard]{
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

		cardNum, err := getCardNumberValue(v)
		if err != nil {
			errs = append(errs, err)
		}

		expDate, err := getExpDateValue(v)
		if err != nil {
			errs = append(errs, err)
		}

		holderName, err := getHolderNameValue(v)
		if err != nil {
			errs = append(errs, err)
		}

		err = handler.RequiredFlagsErr(errs)

		dto := dto.BankCard{
			Name:       name,
			Number:     cardNum,
			ExpDate:    expDate,
			HolderName: holderName,
		}

		return AddFlags{key, dto}, err
	}

	h.GetServiceArgsHook = func(
		f AddFlags,
	) (key string, name string, dto dto.BankCard) {
		return f.Key, f.Name, f.BankCard
	}

	return h
}

type ReadFlags struct {
	Key      string
	EntryNum int
}

func NewRead(
	l logger.Logger, s handler.ReadService[dto.BankCard], w io.Writer,
) *handler.ReadHandler[ReadFlags, dto.BankCard] {
	h := &handler.ReadHandler[ReadFlags, dto.BankCard]{
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

	h.GetOutputHook = func(_ ReadFlags, entryNum int, dto dto.BankCard) string {
		return fmt.Sprintf(
			"the bank card with entry %d: name=%q number=%q exp=%q holder=%q\n",
			entryNum, dto.Name, dto.Number, dto.ExpDate, dto.HolderName,
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
		NamePlural: "bank cards",
	}
}

type EditFlags struct {
	Key      string
	EntryNum int
	dto.BankCard
}

func NewEdit(
	l logger.Logger, s handler.EditService[dto.BankCard], w io.Writer,
) *handler.EditHandler[EditFlags, dto.BankCard] {
	h := &handler.EditHandler[EditFlags, dto.BankCard]{
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

		entryNum, err := handler.GetEnryNumValue(v)
		if err != nil {
			errs = append(errs, err)
		}

		name, err := handler.GetNameValue(v)
		if err != nil {
			errs = append(errs, err)
		}

		cardNum, err := getCardNumberValue(v)
		if err != nil {
			errs = append(errs, err)
		}

		expDate, err := getExpDateValue(v)
		if err != nil {
			errs = append(errs, err)
		}

		holderName, err := getHolderNameValue(v)
		if err != nil {
			errs = append(errs, err)
		}

		err = handler.RequiredFlagsErr(errs)

		dto := dto.BankCard{
			Name:       name,
			Number:     cardNum,
			ExpDate:    expDate,
			HolderName: holderName,
		}

		return EditFlags{key, entryNum, dto}, err
	}

	h.GetServiceArgsHook = func(
		f EditFlags,
	) (key string, entryNum int, name string, dto dto.BankCard) {
		return f.Key, f.EntryNum, f.Name, f.BankCard
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

func getCardNumberValue(v command.ValueGetter) (string, error) {
	cardNum, err := v.GetString(cardcommand.CardNumFlag)
	if err != nil || handler.IsZeroStr(cardNum) {
		return "", fmt.Errorf("--%s", cardcommand.CardNumFlag)
	}
	return cardNum, nil
}

func getExpDateValue(v command.ValueGetter) (string, error) {
	expDate, err := v.GetString(cardcommand.ExpDateFlag)
	if err != nil || handler.IsZeroStr(expDate) {
		return "", fmt.Errorf("--%s", cardcommand.ExpDateFlag)
	}
	return expDate, nil
}

func getHolderNameValue(v command.ValueGetter) (string, error) {
	holderName, err := v.GetString(cardcommand.HolderNameFlag)
	if err != nil || handler.IsZeroStr(holderName) {
		return "", fmt.Errorf("--%s", cardcommand.HolderNameFlag)
	}
	return holderName, nil
}
