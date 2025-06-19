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

	filepathShorthand  = "f"
	filepathDefault    = ""
	readFilepathUsage  = "path to file (required)"
	writeFilepathUsage = "path to file for write data"
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
		FilepathFlag, filepathShorthand, filepathDefault, readFilepathUsage,
	)
	c.MarkFlagRequired(FilepathFlag)

	return &command.Command{Command: c}
}

func NewBinReadCommand(h command.Handler) *command.Command {
	c := &cobra.Command{
		Use: "read",
		Example: "'gophkeeper binary read -k key -e 7 -f /folder/to/file.ext'" +
			" - For write stored data to file.",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), cmd.Flags())
		},
	}
	c.Flags().StringP(
		MasterKeyFlag, masterKeyShorthand, masterKeyDefault, masterKeyUsage,
	)
	c.MarkFlagRequired(MasterKeyFlag)

	c.Flags().IntP(EntryNumFlag, entryNumShorthand, entryNumDefault, entryNumUsage)
	c.MarkFlagRequired(EntryNumFlag)

	c.Flags().StringP(
		FilepathFlag, filepathShorthand, filepathDefault, writeFilepathUsage,
	)

	return &command.Command{Command: c}
}
