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

type Collection struct {
	Items []*Resize
}

type Resize struct {
	File   string `json:"file" binding:"required"`
	Width  string `json:"width"`
	Height string `json:"height"`
}

type ResizeController struct {
	bucket *s3.Bucket
}

func NewResizeController(bucket *s3.Bucket) *ResizeController {
	return &ResizeController{
		bucket,
	}
}

func (rc *ResizeController) Resize(c *gin.Context) {

	for index, element := range c.Request.PostForm("item") {
		fmt.Println(element.Height)
	}

	file, header, err := c.Request.FormFile("file")
	filename := header.Filename

	if err != nil {
		log.Fatal(err)
	}

	height := c.PostForm("height")
	width := c.PostForm("width")

	finalFile := rc.Crop(height, width, file)

	go rc.Upload(filename, finalFile, "image/jpeg", s3.BucketOwnerFull)

	c.JSON(200, gin.H{"filename": filename})
}

func (rc *ResizeController) Upload(filename string, file []byte, enctype string, acl s3.ACL) (error, string) {
	err := rc.bucket.Put(filename, file, enctype, acl)
	return err, filename
}

func (rc *ResizeController) Crop(height string, width string, file multipart.File) []byte {

	image, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	h64, err := strconv.ParseUint(height, 10, 32)
	w64, err := strconv.ParseUint(width, 10, 32)
	h := uint(h64)
	w := uint(w64)

	m := resize.Resize(w, h, image, resize.Lanczos3)

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, m, nil)
	return buf.Bytes()
}
