package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vigneshwaran-48/zshell/utils"
)

var account = &cobra.Command{
	Use:   "account",
	Short: "Account operations",
	Long:  "Account operations",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var accountList = &cobra.Command{
	Use:   "list",
	Short: "Lists all accounts",
	Long:  "Lists all accounts",
	Run: func(cmd *cobra.Command, args []string) {
		dcName, err := cmd.Flags().GetString("dc")
		if err != nil {
			cobra.CheckErr(err)
		}
		accessToken, err := utils.GetAccessToken(dcName)
		if err != nil {
			cobra.CheckErr(err)
		}
		fmt.Println(accessToken)
	},
}

func init() {
	account.AddCommand(accountList)

	rootCmd.AddCommand(account)
}
