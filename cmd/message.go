package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/vigneshwaran-48/zmail-go-sdk"
	"github.com/vigneshwaran-48/zshell/utils"
)

var messageCmd = &cobra.Command{
	Use:   "message",
	Short: "Message commands",
	Long:  "Message commands",
}

var messageListCmd = &cobra.Command{
	Use:    "list",
	Short:  "List messages",
	Long:   "List messages",
	PreRun: ResetPreviousOutput,
	Run: func(cmd *cobra.Command, args []string) {
		accountId, err := cmd.Flags().GetString("account")
		if err != nil {
			cobra.CheckErr(err)
		}

		client, ctx := getAuthDetails(cmd)

		accountId = getAccountId(accountId, client, ctx)

		listEmailReq := client.EmailAPI.ListEmails(ctx, accountId)

		folder, err := cmd.Flags().GetString("folder")
		cobra.CheckErr(err)
		if folder != "" {
			folder = getFolderId(accountId, folder, client, ctx)
			folderId, err := strconv.Atoi(folder)
			if err != nil {
				cobra.CheckErr(err)
			}
			listEmailReq = listEmailReq.FolderId(int64(folderId))
		}

		start, err := cmd.Flags().GetInt("start")
		cobra.CheckErr(err)
		listEmailReq = listEmailReq.Start(int32(start))

		limit, err := cmd.Flags().GetInt("limit")
		cobra.CheckErr(err)
		listEmailReq = listEmailReq.Limit(int32(limit))

		status, err := cmd.Flags().GetString("status")
		cobra.CheckErr(err)
		if status != "" {
			listEmailReq = listEmailReq.Status(status)
		}

		flagid, err := cmd.Flags().GetInt("flagid")
		cobra.CheckErr(err)
		if flagid != -1 {
			listEmailReq = listEmailReq.Flagid(int32(flagid))
		}

		labelid, err := cmd.Flags().GetInt64("labelid")
		cobra.CheckErr(err)
		if labelid != -1 {
			listEmailReq = listEmailReq.Labelid(labelid)
		}

		threadId, err := cmd.Flags().GetInt64("threadId")
		cobra.CheckErr(err)
		if threadId != -1 {
			listEmailReq = listEmailReq.ThreadId(threadId)
		}

		sortBy, err := cmd.Flags().GetString("sortBy")
		cobra.CheckErr(err)
		if sortBy != "" {
			listEmailReq = listEmailReq.SortBy(sortBy)
		}

		sortorder, err := cmd.Flags().GetBool("sortorder")
		cobra.CheckErr(err)
		listEmailReq = listEmailReq.Sortorder(sortorder)

		includeto, err := cmd.Flags().GetBool("includeto")
		cobra.CheckErr(err)
		listEmailReq = listEmailReq.Includeto(includeto)

		includesent, err := cmd.Flags().GetBool("includesent")
		cobra.CheckErr(err)
		listEmailReq = listEmailReq.Includesent(includesent)

		includearchive, err := cmd.Flags().GetBool("includearchive")
		cobra.CheckErr(err)
		listEmailReq = listEmailReq.Includearchive(includearchive)

		attachedMails, err := cmd.Flags().GetBool("attachedMails")
		cobra.CheckErr(err)
		listEmailReq = listEmailReq.AttachedMails(attachedMails)

		inlinedMails, err := cmd.Flags().GetBool("inlinedMails")
		cobra.CheckErr(err)
		listEmailReq = listEmailReq.InlinedMails(inlinedMails)

		flaggedMails, err := cmd.Flags().GetBool("flaggedMails")
		cobra.CheckErr(err)
		listEmailReq = listEmailReq.FlaggedMails(flaggedMails)

		respondedMails, err := cmd.Flags().GetBool("respondedMails")
		cobra.CheckErr(err)
		listEmailReq = listEmailReq.RespondedMails(respondedMails)

		threadedMails, err := cmd.Flags().GetBool("threadedMails")
		cobra.CheckErr(err)
		listEmailReq = listEmailReq.ThreadedMails(threadedMails)

		messagesResp, httpResp, err := listEmailReq.Execute()
		if err != nil {
			handleClientReqError(httpResp, err)
			return
		}

		header := []string{"Message Id", "Subject", "From Address", "To Address"}
		var rows []map[string]string

		for _, message := range messagesResp.Data {
			rows = append(rows, map[string]string{
				"Message Id":   message.GetMessageId(),
				"Subject":      utils.BreakStringIntoLines(message.GetSubject(), 60),
				"From Address": utils.BreakStringIntoLines(message.GetFromAddress(), 40),
				"To Address":   utils.BreakStringIntoLines(message.GetToAddress(), 40),
			})
		}

		SetLastCommandResult(&CmdResult{
			header: header,
			rows:   rows,
		})
	},
}

var messageReadCmd = &cobra.Command{
	Use:   "read",
	Short: "Mark as read message",
	Long:  "Mark as read message",
	Run: func(cmd *cobra.Command, args []string) {
		accountId, err := cmd.Flags().GetString("account")
		if err != nil {
			cobra.CheckErr(err)
		}

		client, ctx := getAuthDetails(cmd)

		accountId = getAccountId(accountId, client, ctx)

		messageIds, err := cmd.Flags().GetStringSlice("message")
		if err != nil {
			cobra.CheckErr(err)
		}
		if len(messageIds) == 0 {
			cobra.CheckErr(errors.New("'message' is required"))
		}

		payload := zmail.MessageUpdatePayload{
			Mode:      zmail.MESSAGEUPDATEMODE_MARK_AS_READ,
			MessageId: messageIds,
		}
		_, httpResp, err := client.EmailAPI.UpdateMessage(ctx, accountId).MessageUpdatePayload(payload).Execute()
		if err != nil {
			handleClientReqError(httpResp, err)
		}
	},
}

var messageUnReadCmd = &cobra.Command{
	Use:   "unread",
	Short: "Mark as unread message",
	Long:  "Mark as unread message",
	Run: func(cmd *cobra.Command, args []string) {
		accountId, err := cmd.Flags().GetString("account")
		if err != nil {
			cobra.CheckErr(err)
		}

		client, ctx := getAuthDetails(cmd)

		accountId = getAccountId(accountId, client, ctx)

		messageIds, err := cmd.Flags().GetStringSlice("message")
		if err != nil {
			cobra.CheckErr(err)
		}
		if len(messageIds) == 0 {
			cobra.CheckErr(errors.New("'message' is required"))
		}

		payload := zmail.MessageUpdatePayload{
			Mode:      zmail.MESSAGEUPDATEMODE_MARK_AS_UNREAD,
			MessageId: messageIds,
		}
		_, httpResp, err := client.EmailAPI.UpdateMessage(ctx, accountId).MessageUpdatePayload(payload).Execute()
		if err != nil {
			handleClientReqError(httpResp, err)
		}
	},
}

var messageMoveCmd = &cobra.Command{
	Use:   "move",
	Short: "Move email",
	Long:  "Move email",
	Run: func(cmd *cobra.Command, args []string) {
		accountId, err := cmd.Flags().GetString("account")
		if err != nil {
			cobra.CheckErr(err)
		}

		client, ctx := getAuthDetails(cmd)

		accountId = getAccountId(accountId, client, ctx)

		messageIds, err := cmd.Flags().GetStringSlice("message")
		if err != nil {
			cobra.CheckErr(err)
		}
		if len(messageIds) == 0 {
			cobra.CheckErr(errors.New("'message' is required"))
		}
		destinationFolderId, err := cmd.Flags().GetString("destination-folder")
		if err != nil {
			cobra.CheckErr(err)
		}
		destinationFolderId = getFolderId(accountId, destinationFolderId, client, ctx)

		payload := zmail.MessageUpdatePayload{
			Mode:         zmail.MESSAGEUPDATEMODE_MOVE_MESSAGE,
			MessageId:    messageIds,
			DestfolderId: &destinationFolderId,
		}
		_, httpResp, err := client.EmailAPI.UpdateMessage(ctx, accountId).MessageUpdatePayload(payload).Execute()
		if err != nil {
			handleClientReqError(httpResp, err)
		}
	},
}

func init() {
	messageCmd.PersistentFlags().String("account", "", "Account Id (Can be id or the account's name)")

	messageListCmd.Flags().String("folder", "", "Folder (Can be id or the folder's path)")
	messageListCmd.Flags().Int("start", 1, "The starting sequence number of the emails.")
	messageListCmd.Flags().Int("limit", 10, "The number of emails to retrieve.")
	messageListCmd.Flags().String("status", "", "Retrieve emails by read or unread status.")
	messageListCmd.Flags().Int("flagid", -1, "The unique identifier for the flag.")
	messageListCmd.Flags().Int64("labelid", -1, "The unique identifier for the label.")
	messageListCmd.Flags().Int64("threadId", -1, "The unique identifier for the thread.")
	messageListCmd.Flags().String("sortBy", "", "The basis on which the sorting should be done.")
	messageListCmd.Flags().Bool("sortorder", false, "The order in which the sorting should be done.")
	messageListCmd.Flags().Bool("includeto", false, "Whether to details need to be included.")
	messageListCmd.Flags().Bool("includesent", false, "Whether sent emails need to be included.")
	messageListCmd.Flags().Bool("includearchive", false, "Whether archived emails need to be included.")
	messageListCmd.Flags().Bool("attachedMails", false, "Retrieve only the emails with attachments.")
	messageListCmd.Flags().Bool("inlinedMails", false, "Retrieve only the emails with inline images.")
	messageListCmd.Flags().Bool("flaggedMails", false, "Retrieve only flagged emails.")
	messageListCmd.Flags().Bool("respondedMails", false, "Retrieve only emails with replies.")
	messageListCmd.Flags().Bool("threadedMails", false, "Retrieve emails that are a part of conversations.")

	messageReadCmd.Flags().StringSlice("message", nil, "Message ids to read")

	messageUnReadCmd.Flags().StringSlice("message", nil, "Message ids to read")

	messageMoveCmd.Flags().StringSlice("message", nil, "Message ids to read")
	messageMoveCmd.Flags().String("destination-folder", "", "Destination folder (Can be id or the folder's path)")

	messageCmd.AddCommand(messageListCmd)
	messageCmd.AddCommand(messageReadCmd)
	messageCmd.AddCommand(messageUnReadCmd)
	messageCmd.AddCommand(messageMoveCmd)

	rootCmd.AddCommand(messageCmd)
}
