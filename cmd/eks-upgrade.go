package cmd

import (
	"github.com/popopanda/magnet/internal/awshelper"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(eksUpgradeCmd)
}

var eksUpgradeCmd = &cobra.Command{
	Use:   "eks-upgrade [profile] [cluster-name] [version]",
	Short: "Upgrade EKS Version to targeted version",
	Args:  cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		awshelper.EKSUpgrade(args)
	},
}
