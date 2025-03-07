package cmd

import (
	"fmt"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/vigneshwaran-48/zshell/utils"
)

var folderCmd = &cobra.Command{
	Use:   "folder",
	Short: "Folder commands",
	Long:  "Folder commands",
}

var folderListCmd = &cobra.Command{
	Use:    "list",
	Short:  "List folders",
	Long:   "List folders of an account",
	PreRun: ResetPreviousOutput,
	Run: func(cmd *cobra.Command, args []string) {
		accountId, err := cmd.Flags().GetString("account")
		if err != nil {
			cobra.CheckErr(err)
		}

		client, ctx := getAuthDetails(cmd)

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
		foldersResponse, httpResp, err := client.FoldersAPI.GetAllFolders(ctx, accountId).Execute()
		if err != nil {
			handleClientReqError(httpResp, err)
		}

		var rows []map[string]string
		for _, folder := range foldersResponse.Data {
			previousFolderId := "-"
			if folder.HasPreviousFolderId() {
				previousFolderId = *folder.PreviousFolderId
			}
			rows = append(rows, map[string]string{
				"Folder Id":       *folder.FolderId,
				"Folder Name":     *folder.FolderName,
				"Type":            *folder.FolderType,
				"Path":            *folder.Path,
				"Previous Folder": previousFolderId,
			})
		}

		lastCmdResult = &CmdResult{
			header: []string{"Folder Id", "Folder Name", "Type", "Path", "Previous Folder"},
			rows:   rows,
		}
	},
}

func init() {
	folderCmd.AddCommand(folderListCmd)

	folderCmd.PersistentFlags().String("account", "", "Account Id")

	rootCmd.AddCommand(folderCmd)
}
