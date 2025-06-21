package synccommand

import (
	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/spf13/cobra"
)

const (
	PasswordFlag = "password"
	LoginFlag    = "login"
)

const (
	loginShorthand = "l"
	loginDefault   = ""
	loginUsage     = "sync account login (required)"

	passwordShorthand = "p"
	passwordDefault   = ""
	passwordUsage     = "sync account password (required)"
)

func New() *command.Command {
	c := &cobra.Command{
		Use:   "sync",
		Short: "Use the sync command to start or stop data synchronization",
	}
	return &command.Command{Command: c}
}

func NewSignup(h command.Handler) *command.Command {
	c := &cobra.Command{
		Use:   "signup",
		Short: "Register new account and start synchronization",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), cmd.Flags())
		},
	}
	flagSet := c.Flags()

	flagSet.StringP(LoginFlag, loginShorthand, loginDefault, loginUsage)

	flagSet.StringP(
		PasswordFlag, passwordShorthand, passwordDefault, passwordUsage,
	)

	c.MarkFlagRequired(LoginFlag)
	c.MarkFlagRequired(PasswordFlag)
	return &command.Command{Command: c}
}

func NewSignin(h command.Handler) *command.Command {
	c := &cobra.Command{
		Use:   "signin",
		Short: "Login and start synchronization",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), cmd.Flags())
		},
	}
	flagSet := c.Flags()

	flagSet.StringP(LoginFlag, loginShorthand, loginDefault, loginUsage)

	flagSet.StringP(
		PasswordFlag, passwordShorthand, passwordDefault, passwordUsage,
	)

	c.MarkFlagRequired(LoginFlag)
	c.MarkFlagRequired(PasswordFlag)
	return &command.Command{Command: c}
}

func NewLogout(h command.Handler) *command.Command {
	c := &cobra.Command{
		Use:   "logout",
		Short: "Logout and stop synchronization",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), cmd.Flags())
		},
	}
	return &command.Command{Command: c}
}
