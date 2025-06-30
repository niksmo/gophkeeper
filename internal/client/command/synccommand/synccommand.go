package synccommand

import (
	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/spf13/cobra"
)

const (
	PasswordFlag = "password"
	LoginFlag    = "login"
	TokenFlag    = "token"
)

const (
	loginShorthand = "l"
	loginDefault   = ""
	loginUsage     = "sync account login (required)"

	passwordShorthand = "p"
	passwordDefault   = ""
	passwordUsage     = "sync account password (required)"

	tokenShorthand = "t"
	tokenDefault   = ""
	tokenUsage     = ""
)

func New() *command.Command {
	c := &cobra.Command{
		Use:   "sync",
		Short: "Use the sync command to start or stop data synchronization",
	}
	return &command.Command{Command: c}
}

type AuthFlags struct {
	Login, Password string
}

func NewSignup(h command.GenCmdHandler[AuthFlags]) *command.Command {
	var fv AuthFlags

	c := &cobra.Command{
		Use:   "signup",
		Short: "Register new account and start synchronization",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), fv)
		},
	}
	flagSet := c.Flags()

	flagSet.StringVarP(&fv.Login,
		LoginFlag, loginShorthand, loginDefault, loginUsage)

	flagSet.StringVarP(&fv.Password,
		PasswordFlag, passwordShorthand, passwordDefault, passwordUsage)

	c.MarkFlagRequired(LoginFlag)
	c.MarkFlagRequired(PasswordFlag)
	return &command.Command{Command: c}
}

func NewSignin(h command.GenCmdHandler[AuthFlags]) *command.Command {
	var fv AuthFlags

	c := &cobra.Command{
		Use:   "signin",
		Short: "Login and start synchronization",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), fv)
		},
	}
	flagSet := c.Flags()

	flagSet.StringVarP(&fv.Login,
		LoginFlag, loginShorthand, loginDefault, loginUsage)

	flagSet.StringVarP(&fv.Password,
		PasswordFlag, passwordShorthand, passwordDefault, passwordUsage)

	c.MarkFlagRequired(LoginFlag)
	c.MarkFlagRequired(PasswordFlag)
	return &command.Command{Command: c}
}

func NewLogout(h command.NoFlagsCmdHandler) *command.Command {
	c := &cobra.Command{
		Use:   "logout",
		Short: "Logout and stop synchronization",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context())
		},
	}
	return &command.Command{Command: c}
}

type StartCmdFlags struct {
	Token string
}

func NewStart(h command.GenCmdHandler[StartCmdFlags]) *command.Command {
	var fv StartCmdFlags
	c := &cobra.Command{
		Hidden: true,
		Use:    "start",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), fv)
		},
	}
	c.Flags().StringVarP(&fv.Token,
		TokenFlag, tokenShorthand, tokenDefault, tokenUsage,
	)
	c.MarkFlagRequired(TokenFlag)

	return &command.Command{Command: c}
}
