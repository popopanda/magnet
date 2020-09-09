package awshelper

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

// AsgScaleDown scale down the ASG by reducing by 50%
func AsgScaleDown(asgNameList []string, awsProfile []string) {
	sess := sessionHelper(awsProfile[0])
	svc := autoscaling.New(sess)

	scaledDownCapacity := int64(0)
	minSizeInt := int64(0)

	for _, i := range asgNameList {
		fmt.Printf("Scaling down the AutoScalingGroup: %v to %v.\n", i, scaledDownCapacity)
		_, errorSetMin := setASGMinSize(i, svc, minSizeInt)

		if errorSetMin != nil {
			fmt.Printf("Unexpeced Error: %v\n", errorSetMin)
			continue
		}

		_, errorSetDesired := setASGDesiredCap(i, svc, scaledDownCapacity)

		if errorSetDesired != nil {
			fmt.Printf("Unexpeced Error: %v\n", errorSetDesired)
			continue
		}
	}
}

func setASGMinSize(asgName string, svc *autoscaling.AutoScaling, minSizeInt int64) (*autoscaling.UpdateAutoScalingGroupOutput, error) {
	updateASGInput := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(asgName),
		MinSize:              aws.Int64(minSizeInt),
	}

	result, err := svc.UpdateAutoScalingGroup(updateASGInput)

	return result, err
}

func setASGDesiredCap(asgName string, svc *autoscaling.AutoScaling, scaledDownInt int64) (*autoscaling.SetDesiredCapacityOutput, error) {
	desiredCapInput := &autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: aws.String(asgName),
		DesiredCapacity:      aws.Int64(scaledDownInt),
		HonorCooldown:        aws.Bool(true),
	}

	result, err := svc.SetDesiredCapacity(desiredCapInput)

	return result, err
}
