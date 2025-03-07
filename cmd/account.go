package cmd

import (
	"github.com/spf13/cobra"
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

func init() {
	account.AddCommand(accountList)

	rootCmd.AddCommand(account)
}
