package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var displayCmd = &cobra.Command{
	Use:   "display",
	Short: "Display last command output",
	Long:  "Display last command output",
	Run: func(cmd *cobra.Command, args []string) {
		if lastCmdResult == nil {
			return
		}
		table := pterm.TableData{
			lastCmdResult.header,
		}
		for _, row := range lastCmdResult.rows {
			cols := []string{}
			for _, header := range lastCmdResult.header {
				cols = append(cols, row[header])
			}
			table = append(table, cols)
		}
		pterm.DefaultTable.WithHasHeader().WithBoxed().WithRowSeparator("-").WithHeaderRowSeparator("-").WithData(table).Render()
	},
}

func init() {
	rootCmd.AddCommand(displayCmd)
}
