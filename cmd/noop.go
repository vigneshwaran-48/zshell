package cmd

import "github.com/spf13/cobra"

var noopCmd = &cobra.Command{
	Use:   "noop",
	Short: "Does nothing",
	Long:  "Does nothing",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(noopCmd)
}
