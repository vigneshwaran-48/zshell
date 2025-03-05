package cmd

import (
	"fmt"

	"github.com/reeflective/console"
	"github.com/spf13/cobra"
)

var app *console.Console

var remainingCmds []string

func GetCmds() *cobra.Command {
	return rootCmd
}

func StartInteractiveShell() {
	app = console.New("ZShell")
	// This hook will run every time when command is executed that includes command executed with ActiveMenu().RunCommandArgs
	app.PostCmdRunHooks = []func() error{
		func() error {
			if len(remainingCmds) > 0 {
				formattedCmds := formatCommand(remainingCmds)
				if len(formattedCmds) > 1 {
					fmt.Println("More commands present")
					remainingCmds = formattedCmds[1]
				} else {
					remainingCmds = []string{} // Reached the end, Reset the remaining command
				}
				app.ActiveMenu().RunCommandArgs(app.ActiveMenu().Context(), formattedCmds[0])
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
