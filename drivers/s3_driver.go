package drivers

import (
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"log"
	"os"
)

// GetBucket
// Authenticates AWS // S3 connection from
// Env variables, then returns the correct bucket.
func GetBucket(bucketName string) *s3.Bucket {
	auth, err := aws.EnvAuth()

	if err != nil {
		log.Panic(err)
	}

	client := s3.New(auth, aws.EUWest)

	log.Println(client.ListBuckets())

	// If not bucket name is given, chekc the env vars
	if bucketName == "" {
		bucketName = os.Getenv("AWS_BUCKET_NAME")
	}

	// If still no bucket name, note we can do really
	if bucketName == "" {
		panic("No bucket detected")
	}

	// Return bucket
	return client.Bucket(bucketName)
}
