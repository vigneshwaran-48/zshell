package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	zmail "github.com/vigneshwaran-48/zmail-go-sdk"
	"github.com/vigneshwaran-48/zshell/utils"
	"golang.org/x/oauth2"
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
		for _, account := range accountsResp.Data {
			data, err := json.Marshal(account)
			if err != nil {
				cobra.CheckErr(err)
			}
			fmt.Println(string(data))
		}
	},
}

func init() {
	account.AddCommand(accountList)

	rootCmd.AddCommand(account)
}
