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
		fmt.Printf("Scaling down the AutoScalingGroup: %v to %v.\n", i, scaledDownCapacity)
		_, errorSetMin := setASGMinSize(i, svc)

		if errorSetMin != nil {
			log.Printf("Unexpeced Error: %v\n", errorSetMin)
			continue
		}

		_, errorSetDesired := setASGDesiredCap(i, svc, scaledDownCapacity)

		if errorSetDesired != nil {
			log.Printf("Unexpeced Error: %v\n", errorSetDesired)
			continue
		}

	}
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

func setASGMinSize(asgName string, svc *autoscaling.AutoScaling) (*autoscaling.UpdateAutoScalingGroupOutput, error) {
	updateASGInput := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(asgName),
		MinSize:              aws.Int64(0),
	}

	result, err := svc.UpdateAutoScalingGroup(updateASGInput)

	return result, err
}

func setASGDesiredCap(asgName string, svc *autoscaling.AutoScaling, scaledDownInt int64) (*autoscaling.SetDesiredCapacityOutput, error) {
	// currentCapacity := asgGetCurrentDesiredCap(i, svc)

	desiredCapInput := &autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: aws.String(asgName),
		DesiredCapacity:      aws.Int64(scaledDownInt),
		HonorCooldown:        aws.Bool(true),
	}

	result, err := svc.SetDesiredCapacity(desiredCapInput)

	return result, err
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
