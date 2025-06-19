package pwdcommand

import (
	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/spf13/cobra"
)

const (
	MasterKeyFlag = command.MasterKeyFlag
	NameFlag      = command.NameFlag
	EntryNumFlag  = command.EntryNumFlag
	PasswordFlag  = "password"
	LoginFlag     = "login"
)

const (
	masterKeyShorthand = command.MasterKeyShorthand
	masterKeyDefault   = command.MasterKeyDefault
	masterKeyUsage     = command.MasterKeyUsage

	nameShorthand = command.NameShorthand
	nameDefault   = command.NameDefault
	nameUsage     = "title for account (required)"

	entryNumShorthand = command.EntryNumShorthand
	entryNumDefault   = command.EntryNumDefault
	entryNumUsage     = "entry number of stored account (required)"

	passwordShorthand = "p"
	passwordDefault   = ""
	passwordUsage     = "account password (required)"

	loginShorthand = "l"
	loginDefault   = ""
	loginUsage     = "account login"
)

func New() *command.Command {
	c := &cobra.Command{
		Use:   "password",
		Short: "Use the password command to save your accounts",
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

	flagSet.StringP(NameFlag, nameShorthand, nameDefault, nameUsage)

	flagSet.StringP(
		PasswordFlag, passwordShorthand, passwordDefault, passwordUsage,
	)

	flagSet.StringP(LoginFlag, loginShorthand, loginDefault, loginUsage)

	c.MarkFlagRequired(MasterKeyFlag)
	c.MarkFlagRequired(NameFlag)
	c.MarkFlagRequired(PasswordFlag)
	return &command.Command{Command: c}
}

func NewRead(h command.Handler) *command.Command {
	c := &cobra.Command{
		Use: "read",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), cmd.Flags())
		},
	}
	flagSet := c.Flags()

	flagSet.StringP(
		MasterKeyFlag, masterKeyShorthand, masterKeyDefault, masterKeyUsage,
	)

	flagSet.IntP(
		EntryNumFlag, entryNumShorthand, entryNumDefault, entryNumUsage,
	)

	c.MarkFlagRequired(MasterKeyFlag)
	c.MarkFlagRequired(EntryNumFlag)

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

	flagSet.StringP(NameFlag, nameShorthand, nameDefault, nameUsage)

	flagSet.StringP(
		PasswordFlag, passwordShorthand, passwordDefault, passwordUsage,
	)

	flagSet.IntP(
		EntryNumFlag, entryNumShorthand, entryNumDefault, entryNumUsage,
	)

	flagSet.StringP(LoginFlag, loginShorthand, loginDefault, loginUsage)

	c.MarkFlagRequired(MasterKeyFlag)
	c.MarkFlagRequired(NameFlag)
	c.MarkFlagRequired(EntryNumFlag)
	c.MarkFlagRequired(PasswordFlag)

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
