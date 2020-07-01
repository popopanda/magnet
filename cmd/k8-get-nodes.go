package cmd

import (
	"fmt"

	"github.com/popopanda/magnet/internal/k8helper"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(k8GetNodesCmd)
}

var k8GetNodesCmd = &cobra.Command{
	Use:   "k8-get-nodes",
	Short: "Get list of nodes in Kubernetes Cluster",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		_, _, combinedList := k8helper.K8GetNodesList()

		for _, i := range combinedList {
			fmt.Println(i)
		}
	},
}
