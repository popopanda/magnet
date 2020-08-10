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
		fmt.Println("This will take a moment...\n")

		nodeNameList, nodeIDList, _ := k8helper.K8GetNodesList()

		fmt.Println("Rolling the following nodes: ")

		for _, i := range nodeNameList {
			fmt.Printf("- %v\n", i)
		}
		fmt.Println("\nProceed with Expanding the AutoScalingGroups?")
		if yesNo() {
			asgNameList := awshelper.AutoScaleUp(nodeIDList, args)
			fmt.Println("\nEnsure new instances have been deployed and registered with the K8 Cluster.\nProceed with draining old nodes?")
			if yesNo() {
				k8helper.K8NodeDrain(nodeNameList)
				fmt.Println("\nProceed with Compacting the AutoScalingGroups?")
				if yesNo() {
					awshelper.AsgScaleDown(asgNameList, args)
				} else {
					os.Exit(1)
				}
			} else {
				os.Exit(1)
			}
		} else {
			os.Exit(1)
		}

	},
}
