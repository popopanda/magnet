package cmd

import (
	"fmt"
	"os"

	"github.com/popopanda/magnet/internal/awshelper"
	"github.com/popopanda/magnet/internal/k8helper"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(k8UpgradeNodesCmd)
}

var k8UpgradeNodesCmd = &cobra.Command{
	Use:   "eks-upgrade-nodes [profile]",
	Short: "Upgrade the k8 nodes in the kubernetes cluster by draining each node, then rolling the EC2 instances",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This will take a moment...")

		nodeNameList, nodeIDList, _ := k8helper.K8GetNodesList()

		fmt.Println("Rolling the following nodes: ")

		for _, i := range nodeNameList {
			fmt.Printf("%v\n", i)
		}

		if yesNo() {
			awshelper.AutoScaleRoll(nodeIDList, args)
			k8helper.K8NodeDrain(nodeNameList)
			awshelper.TerminateEC2(nodeIDList, args)
		} else {
			os.Exit(1)
		}

	},
}
