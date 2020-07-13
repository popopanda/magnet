package awshelper

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// set ASG to add 1 instance,
// terminate old instance
// wait new instance avaiable
// proceed to next

// AutoScaleRoll will expand the ASG, roll the instances
func AutoScaleRoll(instanceIDList []string, awsProfile []string) {
	sess := sessionHelper(awsProfile[0])
	autoScaleGroupName := ""
	for _, instanceID := range instanceIDList {
		autoScaleGroupName = asgLocator(instanceID, sess)
		break
	}

	fmt.Println("Autoscalegroup name: ", autoScaleGroupName)
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

func asgScaleUp(asgName string, sess session.Session) {
	svc := autoscaling.New()
	input := &autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: aws.String("my-auto-scaling-group"),
		DesiredCapacity:      aws.Int64(2),
		HonorCooldown:        aws.Bool(true),
	}

	result, err := svc.SetDesiredCapacity(input)
	if err != nil {
		log.Fatal(err)
	}
}

func ec2Terminate(asgName string) {
}
