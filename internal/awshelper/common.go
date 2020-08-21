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

	for _, i := range asgNameList {

		// currentCapacity := asgGetCurrentDesiredCap(i, svc)
		scaledDownCapacity := int64(0)

		input := &autoscaling.SetDesiredCapacityInput{
			AutoScalingGroupName: aws.String(i),
			DesiredCapacity:      aws.Int64(scaledDownCapacity),
			HonorCooldown:        aws.Bool(true),
		}

		_, err := svc.SetDesiredCapacity(input)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Scaling down the AutoScalingGroup: %v to %v.\n", i, scaledDownCapacity)
	}

}
