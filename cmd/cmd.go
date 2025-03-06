package cmd

import (
	"context"
	"strings"

	"github.com/reeflective/console"
	"github.com/spf13/cobra"
	"github.com/vigneshwaran-48/zmail-go-sdk"
	"github.com/vigneshwaran-48/zshell/utils"
	"golang.org/x/oauth2"
)

type CmdResult struct {
	header []string
	rows   []map[string]string
}

var app *console.Console

var remainingCmds []string

var previousCmd []string

var lastCmdResult *CmdResult

func GetCmds() *cobra.Command {
	return rootCmd
}

// All commands which is not dependant on previous command result can use this as a PreRun hook in the cobra command definition.
func ResetPreviousOutput(cmd *cobra.Command, args []string) {
	lastCmdResult = nil
}

func StartInteractiveShell() {
	app = console.New("ZShell")
	// This hook will run every time when command is executed that includes command executed with ActiveMenu().RunCommandArgs
	app.PostCmdRunHooks = []func() error{
		func() error {
			if len(remainingCmds) > 0 {
				formattedCmds := formatCommand(remainingCmds)
				if len(formattedCmds) > 1 {
					remainingCmds = formattedCmds[1]
				} else {
					remainingCmds = []string{} // Reached the end, Reset the remaining command
				}

				cmdToRun := formattedCmds[0]
				previousCmd = cmdToRun

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
						app.ActiveMenu().RunCommandArgs(context.Background(), currCmd)
					}
				} else {
					app.ActiveMenu().RunCommandArgs(context.Background(), cmdToRun)
				}
			} else {
				if previousCmd[0] != "display" && previousCmd[0] != "help" {
					displayCmd.Run(rootCmd, []string{})
				}
			}
			return nil
		},
	}
	// This hook run only once when the user enters a command
	app.PreCmdRunLineHooks = []func(args []string) ([]string, error){
		func(args []string) ([]string, error) {
			lastCmdResult = nil
			formattedCmds := formatCommand(args)
			if len(formattedCmds) > 1 {
				remainingCmds = formattedCmds[1]
			}
			previousCmd = formattedCmds[0]
			return formattedCmds[0], nil
		},
	}
	app.ActiveMenu().SetCommands(GetCmds)
	app.Start()
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
