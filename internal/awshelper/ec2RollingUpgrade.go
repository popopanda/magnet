package awshelper

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// AutoScaleRoll will expand the ASG, roll the instances
func AutoScaleRoll(instanceIDList []string, awsProfile []string) {
	sess := sessionHelper(awsProfile[0])
	var autoScaleGroupNames []string

	for _, instanceID := range instanceIDList {
		autoScaleGroupNames = append(autoScaleGroupNames, asgLocator(instanceID, sess))
	}

	for _, autoScaleGroup := range autoScaleGroupNames {
		fmt.Println("Scaling up: ", autoScaleGroup)
		asgScaler(autoScaleGroup, *sess)
	}
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

func asgScaler(asgName string, sess session.Session) {
	svc := autoscaling.New(&sess)

	currentCapacity := asgGetCurrentDesiredCap(asgName, svc)
	scaledUpCapacity := currentCapacity + currentCapacity

	input := &autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: aws.String(asgName),
		DesiredCapacity:      aws.Int64(scaledUpCapacity),
		HonorCooldown:        aws.Bool(true),
	}

	_, err := svc.SetDesiredCapacity(input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Adding %v extra nodes to %v.\n", currentCapacity, asgName)
}

func asgGetCurrentDesiredCap(asg string, service *autoscaling.AutoScaling) int64 {
	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{
			aws.String(asg),
		},
	}

	result, err := service.DescribeAutoScalingGroups(input)
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range result.AutoScalingGroups {
		return aws.Int64Value(i.DesiredCapacity)
	}

	return 1
}

// TerminateEC2 terminates a list of ec2 instance ids
func TerminateEC2(instanceIDList []string, awsProfile []string) {
	sess := sessionHelper(awsProfile[0])
	svc := ec2.New(sess)

	for _, instance := range instanceIDList {

		input := &ec2.TerminateInstancesInput{
			InstanceIds: []*string{
				aws.String(instance),
			},
		}

		_, err := svc.TerminateInstances(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return
		}

		fmt.Printf("Terminating: %v\n", instance)
	}

}
