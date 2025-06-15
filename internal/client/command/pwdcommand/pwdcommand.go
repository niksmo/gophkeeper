package pwdcommand

import (
	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/spf13/cobra"
)

const (
	MasterKeyFlag = "master-key"
	NameFlag      = "name"
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
)

func NewPwdCommand() *command.Command {
	c := &cobra.Command{
		Use:   "password",
		Short: "Save accounts data here",
	}
	return &command.Command{c}
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
	return &command.Command{c}
}
