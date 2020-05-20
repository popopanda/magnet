package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var eksUpgradeCmd = &cobra.Command{
	Use:   "aws eks-upgrade-check [profile] [cluster-name]",
	Short: "Print the version number of eks cluster",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

	},
}
