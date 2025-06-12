package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

type CMD struct {
	*cobra.Command
}

func New() *CMD {
	const mKey = "master-key"
	var c = &cobra.Command{
		Use:   "gophkeeper",
		Short: "gophkeeper is a reliable storage of your personal data",
	}
	c.CompletionOptions.DisableDefaultCmd = true
	c.PersistentFlags().StringP(
		mKey, "k", "",
		"required key for encrypting/decrypting and accessing to stored data",
	)
	c.MarkPersistentFlagRequired(mKey)
	return &CMD{c}
}

func (cmd *CMD) Execute() {
	err := cmd.Command.Execute()
	if err != nil {
		os.Exit(1)
	}
}
