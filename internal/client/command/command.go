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

	Command struct {
		*cobra.Command
	}
)

func NewRootCommand() *Command {
	var c = &cobra.Command{
		Use:   "gophkeeper",
		Short: "gophkeeper is a reliable storage of your personal data",
	}
	c.CompletionOptions.DisableDefaultCmd = true
	return &Command{c}
}

func (c *Command) AddCommand(subCmds ...*Command) {
	for _, subCmd := range subCmds {
		c.Command.AddCommand(subCmd.Command)
	}
}
