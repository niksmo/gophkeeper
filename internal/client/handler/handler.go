package handler

import (
	"errors"
	"fmt"
	"io"
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
