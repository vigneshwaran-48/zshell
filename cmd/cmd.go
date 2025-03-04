package cmd

import (
	"fmt"

	"github.com/reeflective/console"
	"github.com/spf13/cobra"
)

var app *console.Console

func GetCmds() *cobra.Command {
	return rootCmd
}

func StartInteractiveShell() {
	app = console.New("ZShell")
	app.ActiveMenu().SetCommands(GetCmds)
	app.Shell().Keymap.Register(map[string]func(){
		"completion-command": func() {
			fmt.Println("completion")
		},
	})
	app.ActiveMenu().PersistentPostRun = func(cmd *cobra.Command, args []string) {
		fmt.Println("Tess")
		if len(args) == 0 {
			return
		}
		if args[0] == "|" {
			fmt.Println("Piping request")
		}
	}
	app.Start()
}
