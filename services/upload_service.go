package services

import (
	"github.com/mitchellh/goamz/s3"
)

type UploadService struct {
	bucket *s3.Bucket
}

func NewUploadService(bucket *s3.Bucket) *UploadService {
	return &UploadService{
		bucket,
	}
}

// Uploads file to S3
func (uploader *UploadService) Upload(filename string, file []byte, enctype string, acl s3.ACL) {

	go func(filename string, file []byte, enctype string, acl s3.ACL) {
		err := uploader.bucket.Put(filename, file, enctype, acl)
		if err != nil {
			panic(err)
		}
	}(filename, file, enctype, acl)
}

// Get - Fetches a single resource form S3
func (uploader *UploadService) Get(filepath string) ([]byte, error) {
	return uploader.bucket.Get(filepath)
}
