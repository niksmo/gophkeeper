package bincommand

import (
	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/spf13/cobra"
)

const (
	MasterKeyFlag = command.MasterKeyFlag
	NameFlag      = command.NameFlag
	EntryNumFlag  = command.EntryNumFlag
	FilepathFlag  = "file"
)

const (
	masterKeyShorthand = command.MasterKeyShorthand
	masterKeyDefault   = command.MasterKeyDefault
	masterKeyUsage     = command.MasterKeyUsage

	nameShorthand = command.NameShorthand
	nameDefault   = command.NameDefault
	nameUsage     = "title for stored file data (required)"

	entryNumShorthand = command.EntryNumShorthand
	entryNumDefault   = command.EntryNumDefault
	entryNumUsage     = "entry number of stored binary data (required)"

	filepathShorthand = "f"
	filepathDefault   = ""
	filepathUsage     = "path to file (required)"
)

func NewBinCommand() *command.Command {
	c := &cobra.Command{
		Use:   "binary",
		Short: "Use the binary command to save your files data",
	}
	return &command.Command{Command: c}
}

func NewBinAddCommand(h command.Handler) *command.Command {
	c := &cobra.Command{
		Use: "add",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), cmd.Flags())
		},
	}

	c.Flags().StringP(
		MasterKeyFlag, masterKeyShorthand, masterKeyDefault, masterKeyUsage,
	)
	c.MarkFlagRequired(MasterKeyFlag)

	c.Flags().StringP(
		NameFlag, nameShorthand, nameDefault, nameUsage,
	)
	c.MarkFlagRequired(NameFlag)

	c.Flags().StringP(
		FilepathFlag, filepathShorthand, filepathDefault, filepathUsage,
	)
	c.MarkFlagRequired(FilepathFlag)

	return &command.Command{Command: c}
}
