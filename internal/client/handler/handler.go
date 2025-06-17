package handler

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/niksmo/gophkeeper/internal/client/command"
)

func RequiredFlagsErr(errs []error) error {
	if len(errs) != 0 {
		return fmt.Errorf(
			"required flags are not specified:\n%w",
			errors.Join(errs...),
		)
	}
	return nil
}

func InternalError(w io.Writer, err error) {
	fmt.Fprintf(
		w,
		"the application completed with an error: %s\n",
		err.Error(),
	)
}

func GetMasterKeyValue(v command.ValueGetter) (string, error) {
	k, err := v.GetString(command.MasterKeyFlag)
	if err != nil || len(strings.TrimSpace(k)) == 0 {
		return "", fmt.Errorf("--%s", command.MasterKeyFlag)
	}
	return k, nil
}

func GetNameValue(v command.ValueGetter) (string, error) {
	n, err := v.GetString(command.NameFlag)
	if err != nil || len(strings.TrimSpace(n)) == 0 {
		return "", fmt.Errorf("--%s", command.NameFlag)
	}
	return n, nil
}
