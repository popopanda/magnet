package cmd

import (
	"fmt"

	"github.com/popopanda/magnet/internal/k8helper"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(k8UpgradeNodesCmd)
}

var k8UpgradeNodesCmd = &cobra.Command{
	Use:   "eks-upgrade-nodes",
	Short: "Upgrade the k8 nodes in the kubernetes cluster by draining each node, then contract and expand the ASGs",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Upgrading...")
		fmt.Println("This will take a moment...")

		nodeNameList, nodeIDList, _ := k8helper.K8GetNodesList()
		fmt.Println(nodeNameList, nodeIDList)

		k8helper.K8NodeDrain(nodeNameList)
	},
}
