package cmd

import (
	"fmt"

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

var aliasCreateCmd = &cobra.Command{
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
		existingAliasWithCmd, err := service.FindByCommand(command)
		if err != nil {
			cobra.CheckErr(err)
		}
		if existingAliasWithCmd != nil {
			confirm, err := pterm.DefaultInteractiveConfirm.WithDefaultText(fmt.Sprintf("Alias '%s' already exists with the same command, Do you want to continue?", existingAliasWithCmd.Name)).Show()
			if err != nil {
				cobra.CheckErr(err)
			}
			if !confirm {
				return
			}
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
		addAliasToRootCmd(alias)
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
		lastCmdResult := &CmdResult{
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
		SetLastCommandResult(lastCmdResult)
	},
}

var aliasRemoveCmd = &cobra.Command{
	Use:    "remove",
	Short:  "Remove alias command",
	Long:   "Remove alias command",
	PreRun: ResetPreviousOutput,
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			cobra.CheckErr(err)
		}
		alias, err := service.FindAliasByName(name)
		if err != nil {
			cobra.CheckErr(err)
		}
		if alias == nil {
			cobra.CheckErr(fmt.Errorf("Given alias command '%s' not exists", name))
		}
		err = service.RemoveAlias(name)
		if err != nil {
			cobra.CheckErr(err)
		}
		customCmd, _, err := rootCmd.Find([]string{alias.Name})
		if err != nil {
			cobra.CheckErr(err)
		}
		rootCmd.RemoveCommand(customCmd)
	},
}

func addAliasToRootCmd(alias *models.Alias) {
	customAliasCmd := &cobra.Command{
		Use:   alias.Name,
		Short: alias.Description,
		Long:  alias.Description,
		Run: func(cmd *cobra.Command, args []string) {
			err := RunCustomCommand(alias.Command)
			if err != nil {
				cobra.CheckErr(err)
			}
		},
	}
	rootCmd.AddCommand(customAliasCmd)
}

func init() {
	aliasCreateCmd.PersistentFlags().String("name", "", "Name of the alias")
	aliasCreateCmd.PersistentFlags().String("command", "", "Command to be aliased")
	aliasCreateCmd.PersistentFlags().String("description", "", "Description of the alias")

	aliasRemoveCmd.PersistentFlags().String("name", "", "Name of the alias to be removed")

	aliasCmd.AddCommand(aliasCreateCmd)
	aliasCmd.AddCommand(aliasListCmd)
	aliasCmd.AddCommand(aliasRemoveCmd)

	aliases, err := service.FindAllAlias()
	if err != nil {
		cobra.CheckErr(err)
	}
	for _, alias := range aliases {
		addAliasToRootCmd(&alias)
	}

	rootCmd.AddCommand(aliasCmd)
}
