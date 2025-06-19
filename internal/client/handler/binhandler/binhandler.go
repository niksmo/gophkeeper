package binhandler

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/niksmo/gophkeeper/internal/client/command/bincommand"
	"github.com/niksmo/gophkeeper/internal/client/dto"
	"github.com/niksmo/gophkeeper/internal/client/handler"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type AddFlags struct {
	Key string
	dto.BIN
}

func NewAddHandler(
	l logger.Logger, s handler.AddService[dto.BIN], w io.Writer,
) *handler.AddHandler[AddFlags, dto.BIN] {
	h := &handler.AddHandler[AddFlags, dto.BIN]{
		Log: l, Service: s, Writer: w, Name: "binary",
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

		data, ext, err := getFileData(filepath)
		if err != nil {
			return AddFlags{}, fmt.Errorf("file error: %w", err)
		}

		dto := dto.BIN{Name: name, Data: data, Ext: ext}
		return AddFlags{key, dto}, err
	}

	h.GetServiceArgsHook = func(f AddFlags) (
		key string, name string, dto dto.BIN,
	) {
		return f.Key, f.Name, f.BIN
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

func getFileData(path string) ([]byte, string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, "", err
	}

	if err := verifyFile(path); err != nil {
		return nil, "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}
	ext := filepath.Ext(path)
	return data, ext, nil
}

func verifyFile(path string) error {
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
