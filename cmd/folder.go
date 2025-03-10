package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/vigneshwaran-48/zmail-go-sdk"
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

		accountId = getAccountId(accountId, client, ctx)

		foldersResponse, httpResp, err := client.FoldersAPI.GetAllFolders(ctx, accountId).Execute()
		if err != nil {
			handleClientReqError(httpResp, err)
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

		SetLastCommandResult(&CmdResult{
			header: []string{"Folder Id", "Folder Name", "Type", "Path", "Previous Folder"},
			rows:   rows,
		})
	},
}

var folderMoveCmd = &cobra.Command{
	Use:    "move",
	Short:  "Move folder",
	Long:   "Move folder",
	PreRun: ResetPreviousOutput,
	Run: func(cmd *cobra.Command, args []string) {
		accountId, err := cmd.Flags().GetString("account")
		if err != nil {
			cobra.CheckErr(err)
		}

		client, ctx := getAuthDetails(cmd)

		accountId = getAccountId(accountId, client, ctx)

		folderId, err := cmd.Flags().GetString("folder")
		if err != nil {
			cobra.CheckErr(err)
		}
		folderId = getFolderId(accountId, folderId, client, ctx)

		parentFolderId, err := cmd.Flags().GetString("parent-folder")
		if err != nil {
			cobra.CheckErr(err)
		}

		previousFolderId, err := cmd.Flags().GetString("previous-folder")
		if err != nil {
			cobra.CheckErr(err)
		}

		payload := zmail.NewFolderUpdatePayload(zmail.FOLDERUPDATEMODE_MOVE)

		if parentFolderId != "" {
			parentFolderId = getFolderId(accountId, parentFolderId, client, ctx)
			if parentFolderId != "" {
				payload.SetParentFolderId(parentFolderId)
			} else {
				cobra.CheckErr(errors.New("Invalid parent folder given"))
			}
		}

		if previousFolderId != "" {
			previousFolderId = getFolderId(accountId, previousFolderId, client, ctx)
			if previousFolderId != "" {
				payload.SetPreviousFolderId(previousFolderId)
			} else {
				cobra.CheckErr(errors.New("Invalid previous folder given"))
			}
		}

		if parentFolderId == "" && previousFolderId == "" {
			cobra.CheckErr(errors.New("--parent-folder or --previous-folder is required"))
		}

		_, httpResp, err := client.FoldersAPI.UpdateFolder(ctx, accountId, folderId).FolderUpdatePayload(*payload).Execute()
		if err != nil {
			handleClientReqError(httpResp, err)
			return
		}
	},
}

var folderRenameCmd = &cobra.Command{
	Use:    "rename",
	Short:  "Rename a folder",
	Long:   "Rename a folder",
	PreRun: ResetPreviousOutput,
	Run: func(cmd *cobra.Command, args []string) {
		accountId, err := cmd.Flags().GetString("account")
		if err != nil {
			cobra.CheckErr(err)
		}

		client, ctx := getAuthDetails(cmd)

		accountId = getAccountId(accountId, client, ctx)

		folderId, err := cmd.Flags().GetString("folder")
		if err != nil {
			cobra.CheckErr(err)
		}
		folderId = getFolderId(accountId, folderId, client, ctx)

		folderName, err := cmd.Flags().GetString("name")
		if err != nil {
			cobra.CheckErr(err)
		}
		if folderName == "" {
			folderName, err = pterm.DefaultInteractiveTextInput.WithDefaultValue("Renamed").Show()
			if err != nil {
				cobra.CheckErr(err)
			}
		}

		payload := zmail.NewFolderUpdatePayload(zmail.FOLDERUPDATEMODE_RENAME)
		payload.SetFolderName(folderName)

		_, httpResp, err := client.FoldersAPI.UpdateFolder(ctx, accountId, folderId).FolderUpdatePayload(*payload).Execute()
		if err != nil {
			handleClientReqError(httpResp, err)
			return
		}
	},
}

var folderEnableImapCmd = &cobra.Command{
	Use:    "enable-imap",
	Short:  "Enables imap view",
	Long:   "Enables imap view",
	PreRun: ResetPreviousOutput,
	Run: func(cmd *cobra.Command, args []string) {
		accountId, err := cmd.Flags().GetString("account")
		if err != nil {
			cobra.CheckErr(err)
		}

		client, ctx := getAuthDetails(cmd)

		accountId = getAccountId(accountId, client, ctx)

		folderId, err := cmd.Flags().GetString("folder")
		if err != nil {
			cobra.CheckErr(err)
		}
		folderId = getFolderId(accountId, folderId, client, ctx)

		payload := zmail.NewFolderUpdatePayload(zmail.FOLDERUPDATEMODE_SET_VIEWABLE_IN_IMAP)
		_, httpResp, err := client.FoldersAPI.UpdateFolder(ctx, accountId, folderId).FolderUpdatePayload(*payload).Execute()
		if err != nil {
			handleClientReqError(httpResp, err)
			return
		}
	},
}

var folderDisableImapCmd = &cobra.Command{
	Use:    "disable-imap",
	Short:  "Disable imap view",
	Long:   "Disable imap view",
	PreRun: ResetPreviousOutput,
	Run: func(cmd *cobra.Command, args []string) {
		accountId, err := cmd.Flags().GetString("account")
		if err != nil {
			cobra.CheckErr(err)
		}

		client, ctx := getAuthDetails(cmd)

		accountId = getAccountId(accountId, client, ctx)

		folderId, err := cmd.Flags().GetString("folder")
		if err != nil {
			cobra.CheckErr(err)
		}
		folderId = getFolderId(accountId, folderId, client, ctx)

		payload := zmail.NewFolderUpdatePayload(zmail.FOLDERUPDATEMODE_REMOVE_VIEWABLE_IN_IMAP)
		_, httpResp, err := client.FoldersAPI.UpdateFolder(ctx, accountId, folderId).FolderUpdatePayload(*payload).Execute()
		if err != nil {
			handleClientReqError(httpResp, err)
			return
		}
	},
}

var folderReadCmd = &cobra.Command{
	Use:    "read",
	Short:  "Mark as Read folder",
	Long:   "Mark as Read folder",
	PreRun: ResetPreviousOutput,
	Run: func(cmd *cobra.Command, args []string) {
		accountId, err := cmd.Flags().GetString("account")
		if err != nil {
			cobra.CheckErr(err)
		}

		client, ctx := getAuthDetails(cmd)

		accountId = getAccountId(accountId, client, ctx)

		folderId, err := cmd.Flags().GetString("folder")
		if err != nil {
			cobra.CheckErr(err)
		}
		folderId = getFolderId(accountId, folderId, client, ctx)

		payload := zmail.NewFolderUpdatePayload(zmail.FOLDERUPDATEMODE_MARK_AS_READ)
		_, httpResp, err := client.FoldersAPI.UpdateFolder(ctx, accountId, folderId).FolderUpdatePayload(*payload).Execute()
		if err != nil {
			handleClientReqError(httpResp, err)
			return
		}
	},
}

var folderEmptyCmd = &cobra.Command{
	Use:    "empty",
	Short:  "Empty folder",
	Long:   "Empty folder",
	PreRun: ResetPreviousOutput,
	Run: func(cmd *cobra.Command, args []string) {
		accountId, err := cmd.Flags().GetString("account")
		if err != nil {
			cobra.CheckErr(err)
		}

		client, ctx := getAuthDetails(cmd)

		accountId = getAccountId(accountId, client, ctx)

		folderId, err := cmd.Flags().GetString("folder")
		if err != nil {
			cobra.CheckErr(err)
		}
		folderId = getFolderId(accountId, folderId, client, ctx)

		payload := zmail.NewFolderUpdatePayload(zmail.FOLDERUPDATEMODE_EMPTY_FOLDER)
		_, httpResp, err := client.FoldersAPI.UpdateFolder(ctx, accountId, folderId).FolderUpdatePayload(*payload).Execute()
		if err != nil {
			handleClientReqError(httpResp, err)
			return
		}
	},
}

var folderDeleteCmd = &cobra.Command{
	Use:    "delete",
	Short:  "Delete folder",
	Long:   "Delete folder",
	PreRun: ResetPreviousOutput,
	Run: func(cmd *cobra.Command, args []string) {
		accountId, err := cmd.Flags().GetString("account")
		if err != nil {
			cobra.CheckErr(err)
		}

		client, ctx := getAuthDetails(cmd)

		accountId = getAccountId(accountId, client, ctx)

		folderId, err := cmd.Flags().GetString("folder")
		if err != nil {
			cobra.CheckErr(err)
		}
		folderId = getFolderId(accountId, folderId, client, ctx)

		_, httpResp, err := client.FoldersAPI.DeleteFolder(ctx, accountId, folderId).Execute()
		if err != nil {
			handleClientReqError(httpResp, err)
			return
		}
	},
}

var folderCreateCmd = &cobra.Command{
	Use:    "create",
	Short:  "Create a folder",
	Long:   "Create a folder",
	PreRun: ResetPreviousOutput,
	Run: func(cmd *cobra.Command, args []string) {
		accountId, err := cmd.Flags().GetString("account")
		if err != nil {
			cobra.CheckErr(err)
		}

		client, ctx := getAuthDetails(cmd)

		accountId = getAccountId(accountId, client, ctx)

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			cobra.CheckErr(err)
		}
		if name == "" {
			name, err = pterm.DefaultInteractiveTextInput.WithDefaultText("Folder Name").Show()
			if err != nil {
				cobra.CheckErr(err)
			}
			if name == "" {
				cobra.CheckErr(errors.New("'name' is required"))
			}
		}
		payload := zmail.NewCreateFolderRequest(name)

		parentFolderId, err := cmd.Flags().GetString("parent-folder")
		if err != nil {
			cobra.CheckErr(err)
		}

		if parentFolderId != "" {
			parentFolderId = getFolderId(accountId, parentFolderId, client, ctx)
			if parentFolderId != "" {
				payload.SetParentFolderId(parentFolderId)
			} else {
				cobra.CheckErr(errors.New("Invalid parent folder given"))
			}
		}

		parentFolderPath, err := cmd.Flags().GetString("parent-path")
		if err != nil {
			cobra.CheckErr(err)
		}

		if parentFolderPath != "" {
			payload.SetParentFolderPath(parentFolderPath)
		}

		_, httpResp, err := client.FoldersAPI.CreateFolder(ctx, accountId).CreateFolderRequest(*payload).Execute()
		if err != nil {
			handleClientReqError(httpResp, err)
		}
	},
}

func getFolderId(accountId string, folderId string, client *zmail.APIClient, ctx context.Context) string {
	newFolderId := ""
	if folderId == "" || !utils.IsNumber(folderId) {
		req := client.FoldersAPI.GetAllFolders(ctx, accountId)
		foldersResp, httpResp, err := req.Execute()
		if err != nil {
			handleClientReqError(httpResp, err)
		}
		if folderId == "" {
			options := []string{}
			for _, folder := range foldersResp.Data {
				options = append(options, fmt.Sprintf("%s (%s)", folder.GetFolderName(), folder.GetFolderId()))
			}
			selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(options).Show("Please select an folder")
			seletedOptionSplit := strings.Split(selectedOption, " ")
			folderIdStr := seletedOptionSplit[len(seletedOptionSplit)-1]
			newFolderId = folderIdStr[1 : len(folderIdStr)-1]
		} else {
			for _, folder := range foldersResp.Data {
				if folder.GetPath() == folderId {
					newFolderId = folder.GetFolderId()
					break
				}
			}
		}
	}
	return newFolderId
}

func init() {
	folderCmd.AddCommand(folderListCmd)
	folderCmd.AddCommand(folderMoveCmd)
	folderCmd.AddCommand(folderRenameCmd)
	folderCmd.AddCommand(folderEnableImapCmd)
	folderCmd.AddCommand(folderDisableImapCmd)
	folderCmd.AddCommand(folderReadCmd)
	folderCmd.AddCommand(folderEmptyCmd)
	folderCmd.AddCommand(folderDeleteCmd)
	folderCmd.AddCommand(folderCreateCmd)

	folderCmd.PersistentFlags().String("account", "", "Account Id (Can be id or the account's name)")

	folderMoveCmd.PersistentFlags().String("folder", "", "Folder (Can be id or the folder's path)")
	folderMoveCmd.PersistentFlags().String("parent-folder", "", "Parent Folder (Can be id or the folder's path)")
	folderMoveCmd.PersistentFlags().String("previous-folder", "", "Previous Folder (Can be id or the folder's path)")

	folderRenameCmd.PersistentFlags().String("folder", "", "Folder (Can be id or the folder's path)")
	folderRenameCmd.PersistentFlags().String("name", "", "New folder name")

	folderEnableImapCmd.PersistentFlags().String("folder", "", "Folder (Can be id or the folder's path)")

	folderDisableImapCmd.PersistentFlags().String("folder", "", "Folder (Can be id or the folder's path)")

	folderReadCmd.PersistentFlags().String("folder", "", "Folder (Can be id or the folder's path)")

	folderEmptyCmd.PersistentFlags().String("folder", "", "Folder (Can be id or the folder's path)")

	folderDeleteCmd.PersistentFlags().String("folder", "", "Folder (Can be id or the folder's path)")

	folderCreateCmd.PersistentFlags().String("name", "", "Folder name")
	folderCreateCmd.PersistentFlags().String("parent-folder", "", "Parent folder (Can be id or the folder's path)")
	folderCreateCmd.PersistentFlags().String("parent-path", "", "Parent folder path")

	rootCmd.AddCommand(folderCmd)
}
