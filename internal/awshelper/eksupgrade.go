package awshelper

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
)

// EKSUpgrade call API to upgrade EKS Cluster
func EKSUpgrade(awsProfile []string) {
	sess := sessionHelper(awsProfile[0])

	svc := eks.New(sess)

	upgradeClusterInput := &eks.UpdateClusterVersionInput{
		Name:    aws.String(awsProfile[1]),
		Version: aws.String(awsProfile[2]),
	}

	result, err := svc.UpdateClusterVersion(upgradeClusterInput)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("EKS Upgrade in progress...\n%v\n%v", result.String, result.Update.Status)
}
