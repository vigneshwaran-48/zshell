package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/vigneshwaran-48/zmail-go-sdk"
	"github.com/vigneshwaran-48/zshell/utils"
)

var account = &cobra.Command{
	Use:   "account",
	Short: "Account operations",
	Long:  "Account operations",
}

var accountList = &cobra.Command{
	Use:    "list",
	Short:  "Lists all accounts",
	Long:   "Lists all accounts",
	PreRun: ResetPreviousOutput,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getAuthDetails(cmd)

		req := client.AccountsAPI.Getmailaccounts(ctx)
		accountsResp, httpResp, err := req.Execute()
		if err != nil {
			handleClientReqError(httpResp, err)
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

func getAccountId(accountId string, client *zmail.APIClient, ctx context.Context) string {
	if accountId == "" || !utils.IsNumber(accountId) {
		req := client.AccountsAPI.Getmailaccounts(ctx)
		accountsResp, httpResp, err := req.Execute()
		if err != nil {
			handleClientReqError(httpResp, err)
		}
		if accountId == "" {
			options := []string{}
			for _, account := range accountsResp.Data {
				options = append(options, fmt.Sprintf("%s (%s)", *account.AccountDisplayName, *account.AccountId))
			}
			selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(options).Show("Please select an account")
			seletedOptionSplit := strings.Split(selectedOption, " ")
			accountIdStr := seletedOptionSplit[len(seletedOptionSplit)-1]
			accountId = accountIdStr[1 : len(accountIdStr)-1]
		} else {
			for _, account := range accountsResp.Data {
				if *account.AccountDisplayName == accountId {
					accountId = *account.AccountId
					break
				}
			}
		}
	}
	return accountId
}

func init() {
	account.AddCommand(accountList)

	rootCmd.AddCommand(account)
}
