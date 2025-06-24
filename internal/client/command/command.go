package command

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

const (
	SecreKeyFlag = "key"
	NameFlag     = "name"
	EntryNumFlag = "entry"

	SecretKeyShorthand = "k"
	SecretKeyDefault   = ""
	SecretKeyUsage     = "key for encrypting, decrypting" +
		" and accessing to stored data (required)"

	NameShorthand = "n"
	NameDefault   = ""

	EntryNumShorthand = "e"
	EntryNumDefault   = 0
)

type (
	GenCmdHandler[F any] interface {
		Handle(context.Context, F)
	}

	ReadCmdHandler interface {
		Handle(ctx context.Context, masterKey string, entryNum int)
	}

	NoFlagsCmdHandler interface {
		Handle(context.Context)
	}

	RemoveCmdHandler interface {
		Handle(ctx context.Context, entryNum int)
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
