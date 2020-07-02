package cmd

import (
	"fmt"

	"github.com/popopanda/magnet/internal/awshelper"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(eksVersionCmd)
}

var eksVersionCmd = &cobra.Command{
	Use:   "eks-version [profile] [cluster-name]",
	Short: "Print the version number of eks cluster",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		eksVersion, eksArn := awshelper.GetEKSVersion(args)

		fmt.Printf("%v\nversion: %v\n", eksArn, eksVersion)

	},
}
