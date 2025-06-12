package cmdpassword

import "github.com/spf13/cobra"

func New() *cobra.Command {
	p := &cobra.Command{
		Use:   "password",
		Short: "Save accounts data here",
	}
	p.AddCommand(
		newAddCMD(),
	)

	return p
}
