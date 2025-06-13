package command

import (
	"context"

	"github.com/spf13/cobra"
)

type (
	ValueGetter interface {
		GetString(name string) (string, error)
	}

	Handler interface {
		Handle(context.Context, ValueGetter)
	}
)

func NewRootCommand() *cobra.Command {
	var c = &cobra.Command{
		Use:   "gophkeeper",
		Short: "gophkeeper is a reliable storage of your personal data",
	}
	c.CompletionOptions.DisableDefaultCmd = true
	return c
}
