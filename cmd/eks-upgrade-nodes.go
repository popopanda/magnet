package cmd

import (
	"fmt"
	"os"

	"github.com/popopanda/magnet/internal/k8helper"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(k8UpgradeNodesCmd)
}

var k8UpgradeNodesCmd = &cobra.Command{
	Use:   "eks-upgrade-nodes",
	Short: "Upgrade the k8 nodes in the kubernetes cluster by draining each node, then rolling the EC2 instances",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This will take a moment...")

		nodeNameList, _, _ := k8helper.K8GetNodesList()

		fmt.Println("Rolling the following nodes: ")

		for _, i := range nodeNameList {
			fmt.Printf("%v\n", i)
		}

		if yesNo() {
			k8helper.K8NodeDrain(nodeNameList)
			//aws terminate instance, and wait for new instance
		} else {
			os.Exit(1)
		}

	},
}
