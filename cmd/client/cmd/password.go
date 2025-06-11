/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// passwordCmd represents the password command
var passwordCmd = &cobra.Command{
	Use:   "password",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("password called", "args:", args)
	// },
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "short description",
	Long:  "Long description",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {
	addCmd.Flags().StringP("login", "l", "", "account login")
	addCmd.Flags().StringP("password", "p", "", "account password")
	addCmd.Flags().StringP("description", "d", "", "description for account")
	// addCmd.MarkFlagRequired("login")
	addCmd.MarkFlagsRequiredTogether("login", "password", "description")
	passwordCmd.AddCommand(addCmd)
	rootCmd.AddCommand(passwordCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// passwordCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// passwordCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
