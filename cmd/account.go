package cmd

import "github.com/spf13/cobra"

var account = &cobra.Command{
	Use:   "account",
	Short: "Account operations",
	Long:  "Account operations",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
