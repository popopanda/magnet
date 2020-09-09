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
	Short: "Blue/Green migration for eks nodegroups",
	Long: `Run this before spinning up the blue/green nodegroups in Terraform. This module assumes the current nodes in the K8 cluster will be depreciated.
	This will grab the current instances, grab their Autoscalegroup IDs, and mark the instances as unschedulable. 
	It will then delete each node of pods, excluding kube-system pods. Then scale down the Autoscalegroups to 0.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This will take a moment...")

		nodeNameList, nodeIDList, _ := k8helper.K8GetNodesList()

		fmt.Println("List of original EC2 workers:")

		for _, i := range nodeNameList {
			fmt.Printf("- %v\n", i)
		}
		fmt.Println("\nBefore continuing, deploy the the blue/green groups via Terraform")
		fmt.Println("\nProceed with the migration?")
		if yesNo() {
			k8helper.K8NodeDrain(nodeNameList)
			fmt.Println("\nProceed with scaling down the the original nodegroups?")
			if yesNo() {
				asgNameList := awshelper.GetAutoScaleGroupList(nodeIDList, args)

				for _, asgName := range asgNameList {
					fmt.Printf("%v will be scaled down.\n", asgName)
				}
				fmt.Println("Proceed?")
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

		fmt.Print("Completed...\nRemember to clean up Terraform code.\n")

	},
}
