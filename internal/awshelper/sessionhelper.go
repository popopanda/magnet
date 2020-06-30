package awshelper

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func sessionHelper(awsProfile string) *session.Session {
	sess, err := session.NewSessionWithOptions(session.Options{
		// Specify profile to load for the session's config
		Profile: awsProfile,

		// Provide SDK Config options, such as Region.
		Config: aws.Config{
			Region: aws.String("us-east-1"),
		},

		// Force enable Shared Config support
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		log.Fatal(err)
	}
	return sess
}
