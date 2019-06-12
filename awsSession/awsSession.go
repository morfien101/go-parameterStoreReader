package awsSession

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
)

// New will create and return an AWS Session
func New() (*session.Session, error) {
	return session.NewSession()
}

// SetRegion publish the supplied region if there is one given
func SetRegion(region string) error {
	// Set what was passed in
	if region != "" {
		err := os.Setenv("AWS_REGION", region)
		if err != nil {
			return err
		}
		return nil
	}
	// If it was blank then look for EC2_REGION
	if value, ok := os.LookupEnv("EC2_REGION"); ok {
		os.Setenv("AWS_REGION", value)
	}
	return nil
}

// SetAccessKey publish AWS_ACCESS_KEY_ID if one is given
func SetAccessKey(accessKey string) error {
	return setEnv(accessKey, "AWS_ACCESS_KEY_ID")
}

// SetSecretKey publish AWS_SECRET_ACCESS_KEY if one is given
func SetSecretKey(secretKey string) error {
	return setEnv(secretKey, "AWS_SECRET_ACCESS_KEY")
}

func setEnv(input string, envName string) error {
	if input != "" {
		if err := os.Setenv(envName, input); err != nil {
			return err
		}
	}
	return nil
}
