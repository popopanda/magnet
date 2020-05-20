package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of magnet",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Magnet v1.0.0")
	},
}
