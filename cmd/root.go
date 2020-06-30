package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "magnet",
	Short: "Magnet simplifies commands",
	Long:  `Created to simplify some tasks that Devops does day to day`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("welcome")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
