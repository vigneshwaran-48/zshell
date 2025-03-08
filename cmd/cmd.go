package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kballard/go-shellquote"
	"github.com/reeflective/console"
	"github.com/spf13/cobra"
	"github.com/vigneshwaran-48/zmail-go-sdk"
	"github.com/vigneshwaran-48/zshell/utils"
	"golang.org/x/oauth2"
)

type ctxKey int

var lastCmdResultKey ctxKey = 0

var remainingCmdsKey ctxKey = 1

var isAliasCmdExecutionKey ctxKey = 2

type CmdResult struct {
	header []string
	rows   []map[string]string
}

var app *console.Console

var previousCmd []string

func GetCmds() *cobra.Command {
	return rootCmd
}

// All commands which is not dependant on previous command result can use this as a PreRun hook in the cobra command definition.
func ResetPreviousOutput(cmd *cobra.Command, args []string) {
	SetLastCommandResult(nil)
}

func SetLastCommandResult(lastCmdResult *CmdResult) {
	ctx := context.WithValue(app.ActiveMenu().Context(), lastCmdResultKey, lastCmdResult)
	app.ActiveMenu().SetContext(ctx)
}

func GetLastCmdResult() *CmdResult {
	value := app.ActiveMenu().Context().Value(lastCmdResultKey)
	if value != nil {
		return value.(*CmdResult)
	}
	return nil
}

func StartInteractiveShell() {
	app = console.New("ZShell")
	// This hook will run every time when command is executed that includes command executed with ActiveMenu().RunCommandArgs
	app.PostCmdRunHooks = []func() error{
		postHook,
	}
	// This hook run only once when the user enters a command
	app.PreCmdRunLineHooks = []func(args []string) ([]string, error){
		preHook,
	}
	app.ActiveMenu().SetCommands(GetCmds)
	app.Start()
}

func RunCustomCommand(command string) error {
	args, err := shellquote.Split(command)
	if err != nil {
		return err
	}
	previousCtx := app.ActiveMenu().Context()
	formattedCmds := formatCommand(args)

	var remainingCmds []string = nil
	if len(formattedCmds) > 1 {
		remainingCmds = formattedCmds[1]
	}

	ctx := context.WithValue(context.Background(), remainingCmdsKey, remainingCmds)

	app.ActiveMenu().RunCommandArgs(context.WithValue(ctx, isAliasCmdExecutionKey, true), formattedCmds[0])

	app.ActiveMenu().SetContext(context.WithValue(previousCtx, lastCmdResultKey, app.ActiveMenu().Context().Value(lastCmdResultKey)))

	return nil
}

func postHook() error {
	remainingCmds := getRemainingCmds()
	ctx := context.WithValue(context.Background(), isAliasCmdExecutionKey, app.ActiveMenu().Context().Value(isAliasCmdExecutionKey))
	if len(remainingCmds) > 0 {
		formattedCmds := formatCommand(remainingCmds)
		if len(formattedCmds) > 1 {
			remainingCmds = formattedCmds[1]
		} else {
			remainingCmds = []string{} // Reached the end, Reset the remaining command
		}
		ctx = context.WithValue(ctx, remainingCmdsKey, remainingCmds)

		cmdToRun := formattedCmds[0]
		previousCmd = cmdToRun

		lastCmdResult := GetLastCmdResult()

		if lastCmdResult != nil && hasVariables(cmdToRun) {
			var headers []string
			headers = append(headers, lastCmdResult.header...)
			for _, row := range lastCmdResult.rows {
				var currCmd []string
				currCmd = append(currCmd, cmdToRun...)

				for i := range currCmd {
					token := &currCmd[i]
					if strings.HasPrefix(*token, "$") && isVariableExists(headers, (*token)[1:]) {
						*token = row[(*token)[1:]]
					}
				}
				ctx = context.WithValue(ctx, lastCmdResultKey, lastCmdResult)
				app.ActiveMenu().RunCommandArgs(ctx, currCmd)
			}
		} else {
			ctx = context.WithValue(ctx, lastCmdResultKey, lastCmdResult)
			app.ActiveMenu().RunCommandArgs(ctx, cmdToRun)
		}
	} else {
		// If it is a alias command execution let the display command to be ran in the alias's context instead of here.
		if app.ActiveMenu().Context().Value(isAliasCmdExecutionKey) == nil && previousCmd[0] != "display" && previousCmd[0] != "help" {
			displayCmd.Run(rootCmd, []string{})
		}
	}
	return nil
}

func preHook(args []string) ([]string, error) {
	ctx := context.Background()
	formattedCmds := formatCommand(args)
	if len(formattedCmds) > 1 {
		remainingCmds := formattedCmds[1]
		ctx = context.WithValue(ctx, remainingCmdsKey, remainingCmds)
		app.ActiveMenu().RunCommandArgs(ctx, formattedCmds[0])
		return []string{"noop"}, nil
	} else {
		previousCmd = formattedCmds[0]
		return formattedCmds[0], nil
	}
}

// Get the first command as one group and rest of them as another group
func formatCommand(cmds []string) [][]string {
	var result [][]string
	var buffer []string

	grouped := false

	for _, cmd := range cmds {
		if !grouped && cmd == "|" {
			grouped = true
			result = append(result, buffer)
			buffer = []string{}
			continue
		}
		buffer = append(buffer, cmd)
	}

	result = append(result, buffer)

	return result
}

func isVariableExists(headers []string, variable string) bool {
	for _, header := range headers {
		if header == variable {
			return true
		}
	}
	return false
}

func getAuthDetails(cmd *cobra.Command) (*zmail.APIClient, context.Context) {
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
	return client, ctx
}

func hasVariables(command []string) bool {
	lastCmdResult := GetLastCmdResult()
	if lastCmdResult == nil || len(lastCmdResult.rows) == 0 {
		return false
	}
	for _, token := range command {
		if strings.HasPrefix(token, "$") && isVariableExists(lastCmdResult.header, token[1:]) {
			return true
		}
	}
	return false
}

func handleClientReqError(httpResp *http.Response, err error) {
	bodyStr, err := io.ReadAll(httpResp.Body)
	if err != nil {
		cobra.CheckErr(err)
	}
	fmt.Println(string(bodyStr))
	cobra.CheckErr(err)
}

func getRemainingCmds() []string {
	remainingCmds := app.ActiveMenu().Context().Value(remainingCmdsKey)
	if remainingCmds != nil {
		return remainingCmds.([]string)
	}
	return nil
}
