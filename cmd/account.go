package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	zmail "github.com/vigneshwaran-48/zmail-go-sdk"
	"github.com/vigneshwaran-48/zshell/utils"
	"golang.org/x/oauth2"
)

var account = &cobra.Command{
	Use:   "account",
	Short: "Account operations",
	Long:  "Account operations",
}

var accountList = &cobra.Command{
	Use:   "list",
	Short: "Lists all accounts",
	Long:  "Lists all accounts",
  PreRun: ResetPreviousOutput,
	Run: func(cmd *cobra.Command, args []string) {
		dcName, err := cmd.Flags().GetString("dc")
		if err != nil {
			cobra.CheckErr(err)
		}
		accessToken, err := utils.GetAccessToken(dcName)
		if err != nil {
			cobra.CheckErr(err)
		}

		config := zmail.NewConfiguration()
		client := zmail.NewAPIClient(config)

		token := &oauth2.Token{
			AccessToken: accessToken,
			TokenType:   "Bearer",
		}

		tokenSource := oauth2.StaticTokenSource(token)

		ctx := context.WithValue(context.Background(), zmail.ContextOAuth2, tokenSource)
		req := client.AccountsAPI.Getmailaccounts(ctx)
		accountsResp, httpResp, err := req.Execute()
		if err != nil {
			data, err := json.Marshal(httpResp)
			if err != nil {
				cobra.CheckErr(err)
			}
			fmt.Println(string(data))
			cobra.CheckErr(err)
		}
		var rows []map[string]string
		for _, account := range accountsResp.Data {
			row := map[string]string{
				"Account Id":   *account.AccountId,
				"Account Name": *account.AccountDisplayName,
				"Email":        *account.AccountDisplayName,
			}
			rows = append(rows, row)
		}
		lastCmdResult = &CmdResult{
			header: []string{"Account Id", "Account Name", "Email"},
			rows:   rows,
		}
	},
}

func init() {
	account.AddCommand(accountList)

	rootCmd.AddCommand(account)
}
