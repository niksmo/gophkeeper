package command

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

const (
	MasterKeyFlag = "master-key"
	NameFlag      = "name"
	EntryNumFlag  = "entry"

	MasterKeyShorthand = "k"
	MasterKeyDefault   = ""
	MasterKeyUsage     = "key for encrypting, decrypting" +
		" and accessing to stored data (required)"

	NameShorthand = "n"
	NameDefault   = ""

	EntryNumShorthand = "e"
	EntryNumDefault   = 0
)

type (
	IntGetter interface {
		GetInt(name string) (int, error)
	}

	StringGetter interface {
		GetString(name string) (string, error)
	}

	ValueGetter interface {
		StringGetter
		IntGetter
	}

	Handler interface {
		Handle(context.Context, ValueGetter)
	}

	Command struct {
		*cobra.Command
	}
)

func NewRootCommand(version, buildDate string) *Command {
	var c = &cobra.Command{
		Use:     "gophkeeper",
		Short:   "gophkeeper is a reliable storage of your personal data",
		Version: fmt.Sprintf("v%s %s", version, buildDate),
	}
	c.CompletionOptions.DisableDefaultCmd = true
	return &Command{c}
}

func (c *Command) AddCommand(subCmds ...*Command) {
	for _, subCmd := range subCmds {
		c.Command.AddCommand(subCmd.Command)
	}
}
