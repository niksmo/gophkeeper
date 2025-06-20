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
	nameUsage     = "name for stored file data (required)"

	entryNumShorthand = command.EntryNumShorthand
	entryNumDefault   = command.EntryNumDefault
	entryNumUsage     = "entry number of stored binary data (required)"

	filepathShorthand  = "f"
	filepathDefault    = ""
	readFilepathUsage  = "path to file (required)"
	writeFilepathUsage = "path to file for write data"
)

func New() *command.Command {
	c := &cobra.Command{
		Use:   "binary",
		Short: "Use the binary command to save your files data",
	}
	return &command.Command{Command: c}
}

func NewAdd(h command.Handler) *command.Command {
	c := &cobra.Command{
		Use: "add",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), cmd.Flags())
		},
	}

	flagSet := c.Flags()

	flagSet.StringP(
		MasterKeyFlag, masterKeyShorthand, masterKeyDefault, masterKeyUsage,
	)

	flagSet.StringP(
		NameFlag, nameShorthand, nameDefault, nameUsage,
	)

	flagSet.StringP(
		FilepathFlag, filepathShorthand, filepathDefault, readFilepathUsage,
	)

	c.MarkFlagRequired(MasterKeyFlag)
	c.MarkFlagRequired(NameFlag)
	c.MarkFlagRequired(FilepathFlag)

	return &command.Command{Command: c}
}

func NewRead(h command.Handler) *command.Command {
	c := &cobra.Command{
		Use: "read",
		Example: "gophkeeper binary read -k 'key' -e 7 -f '/folder/to/file.ext'" +
			" - For write stored data to file.",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), cmd.Flags())
		},
	}
	flagSet := c.Flags()

	flagSet.StringP(
		MasterKeyFlag, masterKeyShorthand, masterKeyDefault, masterKeyUsage,
	)

	flagSet.IntP(EntryNumFlag, entryNumShorthand, entryNumDefault, entryNumUsage)

	flagSet.StringP(
		FilepathFlag, filepathShorthand, filepathDefault, writeFilepathUsage,
	)

	c.MarkFlagRequired(EntryNumFlag)
	c.MarkFlagRequired(MasterKeyFlag)

	return &command.Command{Command: c}
}

func NewList(h command.Handler) *command.Command {
	c := &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), cmd.Flags())
		},
	}
	return &command.Command{Command: c}
}

func NewEdit(h command.Handler) *command.Command {
	c := &cobra.Command{
		Use: "edit",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), cmd.Flags())
		},
	}
	flagSet := c.Flags()

	flagSet.StringP(
		MasterKeyFlag, masterKeyShorthand, masterKeyDefault, masterKeyUsage,
	)

	flagSet.StringP(
		NameFlag, nameShorthand, nameDefault, nameUsage,
	)

	flagSet.StringP(
		FilepathFlag, filepathShorthand, filepathDefault, readFilepathUsage,
	)

	flagSet.IntP(
		EntryNumFlag, entryNumShorthand, entryNumDefault, entryNumUsage,
	)

	c.MarkFlagRequired(MasterKeyFlag)
	c.MarkFlagRequired(NameFlag)
	c.MarkFlagRequired(EntryNumFlag)
	c.MarkFlagRequired(FilepathFlag)

	return &command.Command{Command: c}
}

func NewRemove(h command.Handler) *command.Command {
	c := &cobra.Command{
		Use: "remove",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), cmd.Flags())
		},
	}

	c.Flags().IntP(
		EntryNumFlag, entryNumShorthand, entryNumDefault, entryNumUsage,
	)

	c.MarkFlagRequired(EntryNumFlag)

	return &command.Command{Command: c}
}
