package awshelper

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
)

// GetEKSVersion Obtain AWS EKS Version
func GetEKSVersion(awsProfile []string) (string, string) {

	sess := sessionHelper(awsProfile[0])

	svc := eks.New(sess)
	input := &eks.DescribeClusterInput{
		Name: aws.String(awsProfile[1]),
	}

	result, err := svc.DescribeCluster(input)
	if err != nil {
		log.Fatal(err)
	}

	resultVersion := aws.StringValue(result.Cluster.Version)
	resultArn := aws.StringValue(result.Cluster.Arn)
	return resultVersion, resultArn
}
