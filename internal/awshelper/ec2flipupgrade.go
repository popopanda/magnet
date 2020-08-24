package awshelper

// // AutoScaleUp will expand the ASG, roll the instances
// func AutoScaleUp(instanceIDList []string, awsProfile []string) []string {
// 	sess := sessionHelper(awsProfile[0])
// 	var autoScaleGroupNames []string

// 	for _, instanceID := range instanceIDList {
// 		autoScaleGroupNames = append(autoScaleGroupNames, asgLocator(instanceID, sess))
// 	}

// 	cleanedAutoScaleGroupNames := removeDuplicateFromStringSlice(autoScaleGroupNames)

// 	for _, autoScaleGroup := range cleanedAutoScaleGroupNames {
// 		fmt.Println("Scaling up: ", autoScaleGroup)
// 		asgScaleUp(autoScaleGroup, *sess)
// 	}

// 	return cleanedAutoScaleGroupNames
// }
