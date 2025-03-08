package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/vigneshwaran-48/zshell/models"
	"github.com/vigneshwaran-48/zshell/service"
)

var aliasCmd = &cobra.Command{
	Use:    "alias",
	Short:  "Alias commands",
	Long:   "Alias commands",
	PreRun: ResetPreviousOutput,
}

var aliasCreatCmd = &cobra.Command{
	Use:    "create",
	Short:  "Create a alias",
	Long:   "Store custom commands as alias to the given name",
	PreRun: ResetPreviousOutput,
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			cobra.CheckErr(err)
		}
		if name == "" {
			name, err = pterm.DefaultInteractiveTextInput.WithDefaultText("Alias Name").Show()
			if err != nil {
				cobra.CheckErr(err)
			}
		}
		command, err := cmd.Flags().GetString("command")
		if err != nil {
			cobra.CheckErr(err)
		}
		if command == "" {
			command, err = pterm.DefaultInteractiveTextInput.WithDefaultText("Alias Command").Show()
			if err != nil {
				cobra.CheckErr(err)
			}
		}
		description, err := cmd.Flags().GetString("description")
		if err != nil {
			cobra.CheckErr(err)
		}
		alias := &models.Alias{
			Name:        name,
			Command:     command,
			Description: description,
		}
		err = service.AddAlias(alias)
		if err != nil {
			cobra.CheckErr(err)
		}
	},
}

var aliasListCmd = &cobra.Command{
	Use:    "list",
	Short:  "Lists alias commands",
	Long:   "List alias commands",
	PreRun: ResetPreviousOutput,
	Run: func(cmd *cobra.Command, args []string) {
		aliases, err := service.FindAllAlias()
		if err != nil {
			cobra.CheckErr(err)
		}
		lastCmdResult = &CmdResult{
			header: []string{"Name", "Command", "Description"},
		}
		var rows []map[string]string
		for _, alias := range aliases {
			rows = append(rows, map[string]string{
				"Name":        alias.Name,
				"Command":     alias.Command,
				"Description": alias.Description,
			})
		}
		lastCmdResult.rows = rows
	},
}

func init() {
	aliasCreatCmd.PersistentFlags().String("name", "", "Name of the alias")
	aliasCreatCmd.PersistentFlags().String("command", "", "Command to be aliased")
	aliasCreatCmd.PersistentFlags().String("description", "", "Description of the alias")

	aliasCmd.AddCommand(aliasCreatCmd)
	aliasCmd.AddCommand(aliasListCmd)

	rootCmd.AddCommand(aliasCmd)
}
