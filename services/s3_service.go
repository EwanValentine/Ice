package services

import (
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"log"
)

func S3(bucket string) *s3.Bucket {
	auth, err := aws.EnvAuth()

	if err != nil {
		log.Fatal(err)
	}

	client := s3.New(auth, aws.EUWest)

	bucket, err := client.Bucket(bucket)

	if err != nil {
		log.Fatal(err)
	}

	return bucket
}
