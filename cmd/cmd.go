package cmd

import (
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
	app.Start()
}
