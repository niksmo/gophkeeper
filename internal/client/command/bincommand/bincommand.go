package bincommand

import (
	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/spf13/cobra"
)

const (
	SecretKeyFlag = command.SecreKeyFlag
	NameFlag      = command.NameFlag
	EntryNumFlag  = command.EntryNumFlag
	FilepathFlag  = "file"
)

const (
	secretKeyShorthand = command.SecretKeyShorthand
	secretKeyDefault   = command.SecretKeyDefault
	secretKeyUsage     = command.SecretKeyUsage

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

type AddCmdFlags struct {
	Key, Name, Filepath string
}

func NewAdd(h command.GenCmdHandler[AddCmdFlags]) *command.Command {
	var fv AddCmdFlags

	c := &cobra.Command{
		Use: "add",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), fv)
		},
	}

	flagSet := c.Flags()

	flagSet.StringVarP(&fv.Key,
		SecretKeyFlag, secretKeyShorthand, secretKeyDefault, secretKeyUsage)

	flagSet.StringVarP(&fv.Name,
		NameFlag, nameShorthand, nameDefault, nameUsage)

	flagSet.StringVarP(&fv.Filepath,
		FilepathFlag, filepathShorthand, filepathDefault, readFilepathUsage)

	c.MarkFlagRequired(SecretKeyFlag)
	c.MarkFlagRequired(NameFlag)
	c.MarkFlagRequired(FilepathFlag)

	return &command.Command{Command: c}
}

type EditCmdFlags struct {
	Key, Name, Filepath string
	EntryNum            int
}

func NewEdit(h command.GenCmdHandler[EditCmdFlags]) *command.Command {
	var fv EditCmdFlags

	c := &cobra.Command{
		Use: "edit",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), fv)
		},
	}
	flagSet := c.Flags()

	flagSet.StringVarP(&fv.Key,
		SecretKeyFlag, secretKeyShorthand, secretKeyDefault, secretKeyUsage)

	flagSet.StringVarP(&fv.Name,
		NameFlag, nameShorthand, nameDefault, nameUsage)

	flagSet.StringVarP(&fv.Filepath,
		FilepathFlag, filepathShorthand, filepathDefault, readFilepathUsage)

	flagSet.IntVarP(&fv.EntryNum,
		EntryNumFlag, entryNumShorthand, entryNumDefault, entryNumUsage)

	c.MarkFlagRequired(SecretKeyFlag)
	c.MarkFlagRequired(NameFlag)
	c.MarkFlagRequired(EntryNumFlag)
	c.MarkFlagRequired(FilepathFlag)

	return &command.Command{Command: c}
}

type ReadCmdFlags struct {
	Key, Filepath string
	EntryNum      int
}

func NewRead(h command.GenCmdHandler[ReadCmdFlags]) *command.Command {
	var fv ReadCmdFlags

	c := &cobra.Command{
		Use: "read",
		Example: "gophkeeper binary read -k 'key' -e 7 -f '/folder/to/file.ext'" +
			" - For write stored data to file.",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), fv)
		},
	}

	flagSet := c.Flags()

	flagSet.StringVarP(&fv.Key,
		SecretKeyFlag, secretKeyShorthand, secretKeyDefault, secretKeyUsage)

	flagSet.IntVarP(&fv.EntryNum,
		EntryNumFlag, entryNumShorthand, entryNumDefault, entryNumUsage)

	flagSet.StringVarP(&fv.Filepath,
		FilepathFlag, filepathShorthand, filepathDefault, writeFilepathUsage)

	c.MarkFlagRequired(EntryNumFlag)
	c.MarkFlagRequired(SecretKeyFlag)

	return &command.Command{Command: c}
}

func NewList(h command.NoFlagsCmdHandler) *command.Command {
	c := &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context())
		},
	}
	return &command.Command{Command: c}
}

func NewRemove(h command.RemoveCmdHandler) *command.Command {
	var entryNum int

	c := &cobra.Command{
		Use: "remove",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), entryNum)
		},
	}

	c.Flags().IntVarP(&entryNum,
		EntryNumFlag, entryNumShorthand, entryNumDefault, entryNumUsage)

	c.MarkFlagRequired(EntryNumFlag)

	return &command.Command{Command: c}
}
