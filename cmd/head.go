package cmd

import "github.com/spf13/cobra"

var headCmd = &cobra.Command{
	Use:  "head",
	Long: "Fetch the first n number of previous result rows",
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
		lastCmdResult.rows = lastCmdResult.rows[0:lines]
	},
}

func init() {
	headCmd.PersistentFlags().Int16P("lines", "n", 0, "Number of rows")

	rootCmd.AddCommand(headCmd)
}
