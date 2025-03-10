package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vigneshwaran-48/zshell/utils"
)

var login = &cobra.Command{
	Use:   "login",
	Short: "Login to a DC",
	Long:  "Login to a DC",
	Run: func(cmd *cobra.Command, args []string) {
		dcName, err := cmd.Flags().GetString("dc")
		if err != nil {
			cobra.CheckErr(err)
		}
		_, err = utils.LoginToDC(dcName, password)
		if err != nil {
			cobra.CheckErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(login)
}
