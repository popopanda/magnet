package awshelper

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
)

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

// AsgScaleDown scale down the ASG by reducing by 50%
func AsgScaleDown(asgNameList []string, awsProfile []string) {
	sess := sessionHelper(awsProfile[0])
	svc := autoscaling.New(sess)

	scaledDownCapacity := int64(0)

	for _, i := range asgNameList {
		setASGMinSize(i, svc)
		setASGDesiredCap(i, svc, scaledDownCapacity)
		fmt.Printf("Scaling down the AutoScalingGroup: %v to %v.\n", i, scaledDownCapacity)
	}
}

func setASGMinSize(asgName string, svc *autoscaling.AutoScaling) {
	updateASGInput := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(asgName),
		MinSize:              aws.Int64(0),
	}

	_, err := svc.UpdateAutoScalingGroup(updateASGInput)

	if err != nil {
		log.Fatal(err)
	}
}

func setASGDesiredCap(asgName string, svc *autoscaling.AutoScaling, scaledDownInt int64) {
	// currentCapacity := asgGetCurrentDesiredCap(i, svc)

	desiredCapInput := &autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: aws.String(asgName),
		DesiredCapacity:      aws.Int64(scaledDownInt),
		HonorCooldown:        aws.Bool(true),
	}

	_, err := svc.SetDesiredCapacity(desiredCapInput)
	if err != nil {
		log.Fatal(err)
	}
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
