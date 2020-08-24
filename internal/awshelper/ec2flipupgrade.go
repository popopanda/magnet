package awshelper

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// AutoScaleUp will expand the ASG, roll the instances
func AutoScaleUp(instanceIDList []string, awsProfile []string) []string {
	sess := sessionHelper(awsProfile[0])
	var autoScaleGroupNames []string

	for _, instanceID := range instanceIDList {
		autoScaleGroupNames = append(autoScaleGroupNames, asgLocator(instanceID, sess))
	}

	cleanedAutoScaleGroupNames := removeDuplicateFromStringSlice(autoScaleGroupNames)

	for _, autoScaleGroup := range cleanedAutoScaleGroupNames {
		fmt.Println("Scaling up: ", autoScaleGroup)
		asgScaleUp(autoScaleGroup, *sess)
	}

	return cleanedAutoScaleGroupNames
}

// GetAutoScaleGroupList get list of autoscalegroups used by the existing nodegroups
func GetAutoScaleGroupList(instanceIDList []string, awsProfile []string) []string {
	sess := sessionHelper(awsProfile[0])
	var autoScaleGroupNames []string

	for _, instanceID := range instanceIDList {
		autoScaleGroupNames = append(autoScaleGroupNames, asgLocator(instanceID, sess))
	}

	cleanedAutoScaleGroupNames := removeDuplicateFromStringSlice(autoScaleGroupNames)

	return cleanedAutoScaleGroupNames
}

func asgLocator(instanceID string, sess *session.Session) string {
	svc := ec2.New(sess)

	ec2Input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	}

	result, err := svc.DescribeInstances(ec2Input)
	if err != nil {
		log.Fatal(err)

	}

	for _, reservations := range result.Reservations {
		for _, tags := range reservations.Instances {
			for _, kv := range tags.Tags {
				if aws.StringValue(kv.Key) == "aws:autoscaling:groupName" {
					return aws.StringValue(kv.Value)
				}
			}
		}
	}
	return ""
}

func removeDuplicateFromStringSlice(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, i := range stringSlice {
		if _, value := keys[i]; !value {
			keys[i] = true
			list = append(list, i)
		}
	}
	return list
}

func asgScaleUp(asgName string, sess session.Session) {
	svc := autoscaling.New(&sess)

	currentCapacity := asgGetCurrentDesiredCap(asgName, svc)
	scaledUpCapacity := currentCapacity * 2

	input := &autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: aws.String(asgName),
		DesiredCapacity:      aws.Int64(scaledUpCapacity),
		HonorCooldown:        aws.Bool(true),
	}

	_, err := svc.SetDesiredCapacity(input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Expanding (Doubling) the AutoScalingGroup: %v.\n", asgName)
}
