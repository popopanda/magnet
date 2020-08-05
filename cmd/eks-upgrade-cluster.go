package cmd

import (
	"fmt"
	"os"

	"github.com/popopanda/magnet/internal/awshelper"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(eksUpgradeCmd)
}

var eksUpgradeCmd = &cobra.Command{
	Use:   "eks-upgrade-cluster [profile] [cluster-name] [version]",
	Short: "Upgrade EKS Cluster version to targeted version",
	Args:  cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Proceed with upgrade?")
		if yesNo() {
			awshelper.EKSUpgrade(args)
		} else {
			os.Exit(1)
		}
	},
}
