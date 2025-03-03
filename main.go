package main

import (
	"github.com/reeflective/console"
	"github.com/vigneshwaran-48/zshell/cmd"
)

func main() {
	app := console.New("ZShell")
	app.ActiveMenu().SetCommands(cmd.GetCmds)
	app.Start()
}
