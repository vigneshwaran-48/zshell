package cmd

import (
	"strings"

	"github.com/reeflective/console"
	"github.com/spf13/cobra"
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

				if lastCmdResult != nil {
					for i := range cmdToRun {
						token := &cmdToRun[i]
						if strings.HasPrefix(*token, "$") {
							variableName := (*token)[1:]
							if !isVariableExists(variableName) {
								continue
							}
							*token = lastCmdResult.rows[0][variableName] // Need to iterate all rows
						}
					}
				}
				previousCmd = cmdToRun

				app.ActiveMenu().RunCommandArgs(app.ActiveMenu().Context(), cmdToRun)
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

func isVariableExists(variable string) bool {
	for _, header := range lastCmdResult.header {
		if header == variable {
			return true
		}
	}
	return false
}
