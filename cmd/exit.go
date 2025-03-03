package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var exit = &cobra.Command{
	Use:   "exit",
	Short: "Exits the application",
	Long:  "Exits the application",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(exit)
}
