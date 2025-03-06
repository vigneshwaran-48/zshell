package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
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
		foldersResponse, httpResp, err := client.FoldersAPI.GetAllFolders(ctx, accountId).Execute()
		if err != nil {
			bodyStr, err := io.ReadAll(httpResp.Body)
			if err != nil {
				cobra.CheckErr(err)
			}
			fmt.Println(string(bodyStr))
			cobra.CheckErr(err)
			return
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

	folderCmd.PersistentFlags().String("account", "-1", "Account Id")

	rootCmd.AddCommand(folderCmd)
}
