package data

import (
	"github.com/AdRoll/goamz/aws"
	"github.com/AdRoll/goamz/s3"
)

var (
	Bucket *s3.Bucket
)

func InitBucket(name string) error {
	auth, err := aws.EnvAuth()
	if err != nil {
		return err
	}
	s3 := s3.New(auth, aws.GetRegion("us-east-1"))
	Bucket = s3.Bucket(name)
	return nil
}
