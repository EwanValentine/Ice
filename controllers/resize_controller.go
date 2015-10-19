package controllers

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/goamz/s3"
	"github.com/nfnt/resize"
	_ "image"
	"image/jpeg"
	"log"
	"mime/multipart"
	_ "strconv"
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
	if err != nil {
		panic(err)
	}

	// Get original filename
	filename := header.Filename

	fmt.Println(c.Request.MultipartForm.Value["set"])

	// Foreach set of dimensions given
	for _, element := range c.Request.MultipartForm.Value["set"] {

		// Get height and width
		height := element[0]
		width := element[1]

		fmt.Println(element[0])
		fmt.Println(element[1])

		// Crop file
		finalFile := rc.Crop(height, width, file)

		// Upload file
		go rc.Upload(filename, finalFile, "image/jpeg", s3.BucketOwnerFull)
	}

	c.JSON(200, gin.H{"filename": filename})
}

// Uploads file to S3
func (rc *ResizeController) Upload(filename string, file []byte, enctype string, acl s3.ACL) (error, string) {
	err := rc.bucket.Put(filename, file, enctype, acl)
	return err, filename
}

// Crops image and returns []byte of file
func (rc *ResizeController) Crop(height uint8, width uint8, file multipart.File) []byte {

	// Decode image
	image, err := jpeg.Decode(file)

	// If error, of course
	if err != nil {
		panic(err)
	}

	// Convert string height and width into 64int
	//w64, err := strconv.ParseUint(width, 10, 32)
	// h64, err := strconv.ParseUint(height, 10, 32)

	// Convert 64int to 32int
	h := uint(height)
	w := uint(width)

	// Runs re-size function
	m := resize.Resize(w, h, image, resize.Lanczos3)

	// Create new buffer of file
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, m, nil)

	if err != nil {
		panic(err)
	}

	// Return buffer
	return buf.Bytes()
}
