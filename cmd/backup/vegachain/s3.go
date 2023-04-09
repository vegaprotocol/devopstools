package vegachain

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Credentials struct {
	Endpoint     string
	Region       string
	AccessKey    string
	AccessSecret string
}

func NewSession(cred S3Credentials) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      &cred.Region,
		Endpoint:    &cred.Endpoint,
		Credentials: credentials.NewStaticCredentials(cred.AccessKey, cred.AccessSecret, ""),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create aws session: %w", err)
	}

	client := s3.New(sess)
	if _, err = client.ListBuckets(&s3.ListBucketsInput{}); err != nil {
		return nil, fmt.Errorf("given credentials \"%s\" do not have access to s3", cred.AccessKey)
	}

	return sess, nil
}
