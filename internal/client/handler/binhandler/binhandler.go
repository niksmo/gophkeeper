package binhandler

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/niksmo/gophkeeper/internal/client/command/bincommand"
	"github.com/niksmo/gophkeeper/internal/client/dto"
	"github.com/niksmo/gophkeeper/internal/client/handler"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

const handlerName = "binary"

type AddFlags struct {
	Key string
	dto.BIN
}

func NewAdd(
	l logger.Logger, s handler.AddService[dto.BIN], w io.Writer,
) *handler.AddHandler[AddFlags, dto.BIN] {
	h := &handler.AddHandler[AddFlags, dto.BIN]{
		Log: l, Service: s, Writer: w, Name: handlerName,
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

		filepath, err := getFilePathValue(v)
		if err != nil {
			errs = append(errs, err)
		}

		err = handler.RequiredFlagsErr(errs)

		if err != nil {
			return AddFlags{}, err
		}

		dto, err := getBinDto(name, filepath)
		if err != nil {
			return AddFlags{}, err
		}

		return AddFlags{key, dto}, err
	}

	h.GetServiceArgsHook = func(f AddFlags) (
		key string, name string, dto dto.BIN,
	) {
		return f.Key, f.Name, f.BIN
	}

	return h
}

type ReadFlags struct {
	Key, Filepath string
	EntryNum      int
}

func NewRead(
	l logger.Logger, s handler.ReadService[dto.BIN], w io.Writer,
) *handler.ReadHandler[ReadFlags, dto.BIN] {
	h := &handler.ReadHandler[ReadFlags, dto.BIN]{
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

		filepath, _ := getFilePathValue(v)

		err = handler.RequiredFlagsErr(errs)
		return ReadFlags{key, filepath, entryNum}, err
	}

	h.GetServiceArgsHook = func(f ReadFlags) (key string, entryNum int) {
		return f.Key, f.EntryNum
	}

	h.GetOutputHook = func(flags ReadFlags, entryNum int, dto dto.BIN) string {
		var strBuilder strings.Builder
		strBuilder.WriteString(
			fmt.Sprintf(
				"the binary data with entry %d: name=%q size=%d ext=%q \n",
				entryNum, dto.Name, len(dto.Data), dto.Ext,
			))

		err := writeData(flags.Filepath, dto.Data)
		if err != nil {
			strBuilder.WriteString(
				fmt.Sprintf("write file error: %s\n", err.Error()),
			)
		}

		return strBuilder.String()
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
		NamePlural: "binaries",
	}
}

type EditFlags struct {
	Key      string
	EntryNum int
	dto.BIN
}

func NewEdit(
	l logger.Logger, s handler.EditService[dto.BIN], w io.Writer,
) *handler.EditHandler[EditFlags, dto.BIN] {
	h := &handler.EditHandler[EditFlags, dto.BIN]{
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

		filepath, err := getFilePathValue(v)
		if err != nil {
			errs = append(errs, err)
		}

		err = handler.RequiredFlagsErr(errs)

		if err != nil {
			return EditFlags{}, err
		}

		dto, err := getBinDto(name, filepath)
		if err != nil {
			return EditFlags{}, err
		}

		return EditFlags{key, entryNum, dto}, err
	}

	h.GetServiceArgsHook = func(
		f EditFlags,
	) (key string, entryNum int, name string, dto dto.BIN) {
		return f.Key, f.EntryNum, f.Name, f.BIN
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

func getFilePathValue(v command.ValueGetter) (string, error) {
	filepath, err := v.GetString(bincommand.FilepathFlag)
	if err != nil || handler.IsZeroStr(filepath) {
		return "", fmt.Errorf("--%s", command.MasterKeyFlag)
	}
	return filepath, nil
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
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
