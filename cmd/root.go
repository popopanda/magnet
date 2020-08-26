package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "magnet",
	Short: "Magnet simplifies commands",
	Long:  `Created to simplify some tasks that Devops does day to day`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run magnet -h to get list of commands")
	},
}

// Execute entry
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func yesNo() bool {
	prompt := promptui.Select{
		Label: "Continue? [Yes/No]",
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result == "Yes"
}
