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
	masterKeyShorthand = "k"
	masterKeyDefault   = ""
	masterKeyUsage     = "key for encrypting, decrypting" +
		" and accessing to stored data (required)"

	nameShorthand = "n"
	nameDefault   = ""
	nameUsage     = "title for account (required)"

	passwordShorthand = "p"
	passwordDefault   = ""
	passwordUsage     = "account password (required)"

	loginShorthand = "l"
	loginDefault   = ""
	loginUsage     = "account login"

	entryNumShorthand = "e"
	entryNumDefault   = 0
	entryNumUsage     = "stored account entry number"
)

func NewPwdCommand() *command.Command {
	c := &cobra.Command{
		Use:   "password",
		Short: "Save accounts data here",
	}
	return &command.Command{Command: c}
}

func NewPwdAddCommand(h command.Handler) *command.Command {
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

func NewPwdReadCommand(h command.Handler) *command.Command {
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

func NewPwdListCommand(h command.Handler) *command.Command {
	c := &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), cmd.Flags())
		},
	}
	return &command.Command{Command: c}
}

func NewPwdEditCommand(h command.Handler) *command.Command {
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

	c.Flags().StringP(PasswordFlag, passwordShorthand, passwordDefault, passwordUsage)
	c.MarkFlagRequired(PasswordFlag)

	c.Flags().IntP(EntryNumFlag, entryNumShorthand, entryNumDefault, entryNumUsage)
	c.MarkFlagRequired(EntryNumFlag)

	c.Flags().StringP(LoginFlag, loginShorthand, loginDefault, loginUsage)
	return &command.Command{Command: c}
}

func NewPwdRemoveCommand(h command.Handler) *command.Command {
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
