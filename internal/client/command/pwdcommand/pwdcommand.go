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

	c.Flags().StringP(
		MasterKeyFlag, masterKeyShorthand, masterKeyDefault, masterKeyUsage,
	)
	c.MarkFlagRequired(MasterKeyFlag)

	c.Flags().StringP(NameFlag, nameShorthand, nameDefault, nameUsage)
	c.MarkFlagRequired(NameFlag)

	c.Flags().StringP(PasswordFlag, passwordShorthand, passwordDefault, passwordUsage)
	c.MarkFlagRequired(PasswordFlag)

	c.Flags().StringP(LoginFlag, loginShorthand, loginDefault, loginUsage)
	return &command.Command{Command: c}
}

func NewRead(h command.Handler) *command.Command {
	c := &cobra.Command{
		Use: "read",
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

	c.Flags().StringP(
		MasterKeyFlag, masterKeyShorthand, masterKeyDefault, masterKeyUsage,
	)
	c.MarkFlagRequired(MasterKeyFlag)

	c.Flags().StringP(NameFlag, nameShorthand, nameDefault, nameUsage)
	c.MarkFlagRequired(NameFlag)

	c.Flags().StringP(
		PasswordFlag, passwordShorthand, passwordDefault, passwordUsage,
	)
	c.MarkFlagRequired(PasswordFlag)

	c.Flags().IntP(
		EntryNumFlag, entryNumShorthand, entryNumDefault, entryNumUsage,
	)
	c.MarkFlagRequired(EntryNumFlag)

	c.Flags().StringP(LoginFlag, loginShorthand, loginDefault, loginUsage)
	return &command.Command{Command: c}
}

func NewRemove(h command.Handler) *command.Command {
	c := &cobra.Command{
		Use: "remove",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), cmd.Flags())
		},
	}
	c.Flags().IntP(EntryNumFlag, entryNumShorthand, entryNumDefault, entryNumUsage)
	c.MarkFlagRequired(EntryNumFlag)
	return &command.Command{Command: c}
}
