package cmdpassword

import "github.com/spf13/cobra"

func newAddCMD() *cobra.Command {
	c := &cobra.Command{
		Use: "add",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	c.Flags().StringP("name", "n", "", "title for account (required)")
	c.Flags().StringP("login", "l", "", "account login")
	c.Flags().StringP("password", "p", "", "account password (required)")
	c.MarkFlagRequired("name")
	c.MarkFlagRequired("password")
	return c
}
