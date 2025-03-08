package cmd

import (
	"github.com/spf13/cobra"
)

var tailCmd = &cobra.Command{
	Use:  "tail",
	Long: "Fetch the last n number of previous result rows",
	Run: func(cmd *cobra.Command, args []string) {
		lastCmdResult := GetLastCmdResult()
		if lastCmdResult == nil {
			return
		}
		lines, err := cmd.Flags().GetInt16("lines")
		if err != nil {
			cobra.CheckErr(err)
		}
		if int(lines) > len(lastCmdResult.rows) {
			lines = int16(len(lastCmdResult.rows))
		}
		startIndex := len(lastCmdResult.rows) - int(lines)
		lastCmdResult.rows = lastCmdResult.rows[startIndex:len(lastCmdResult.rows)]
	},
}

func init() {
	tailCmd.PersistentFlags().Int16P("lines", "n", 0, "Number of rows")

	rootCmd.AddCommand(tailCmd)
}
