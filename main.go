package main

import (
	"os"

	"github.com/reeflective/console"
	"github.com/spf13/cobra"
	"github.com/vigneshwaran-48/zshell/cmd"
)

func main() {
	if len(os.Args) > 1 {
		err := cmd.GetCmds().Execute()
		if err != nil {
			cobra.CheckErr(err)
		}
		return
	}
	app := console.New("ZShell")
	app.ActiveMenu().SetCommands(cmd.GetCmds)
	app.Start()
}
