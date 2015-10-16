package controllers

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/goamz/s3"
	"github.com/nfnt/resize"
	"image/jpeg"
	"log"
	"mime/multipart"
	"strconv"
)

type ResizeController struct {
	bucket *s3.Bucket
}

func NewResizeController(bucket *s3.Bucket) *ResizeController {
	return &ResizeController{
		bucket,
	}
}

// Resize - function for taking images and meta data and resizing
func (rc *ResizeController) Resize(c *gin.Context) {

	// Don't really know what this bit does to be fair
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		log.Fatal(err)
	}

	// Get file
	file, header, err := c.Request.FormFile("file")

	fmt.Println(file)

	// Get original filename
	filename := header.Filename
	if err != nil {
		log.Fatal(err)
	}

	// Foreach set of dimensions given
	for _, element := range c.Request.MultipartForm.Value {

		// Get height and width
		height := element[0]
		width := element[0]

		// Crop file
		finalFile := rc.Crop(height, width, file)

		fmt.Println(finalFile)
		fmt.Println(height)
		fmt.Println(width)

		// Upload file
		// go rc.Upload(filename, finalFile, "image/jpeg", s3.BucketOwnerFull)
	}

	c.JSON(200, gin.H{"filename": filename})
}

// Uploads file to S3
func (rc *ResizeController) Upload(filename string, file []byte, enctype string, acl s3.ACL) (error, string) {
	err := rc.bucket.Put(filename, file, enctype, acl)
	return err, filename
}

// Crops image and returns []byte of file
func (rc *ResizeController) Crop(height string, width string, file multipart.File) []byte {

	// Decode file to image file
	image, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	// Convert string height and width into 64int
	h64, err := strconv.ParseUint(height, 10, 32)
	w64, err := strconv.ParseUint(width, 10, 32)

	// Convert 64int to 32int
	h := uint(h64)
	w := uint(w64)

	// Runs re-size function
	m := resize.Resize(w, h, image, resize.Lanczos3)

	// Create new buffer of file
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, m, nil)
	return buf.Bytes()
}
