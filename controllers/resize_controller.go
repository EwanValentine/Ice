package controllers

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"github.com/nfnt/resize"
	"image/jpeg"
	"log"
	"mime/multipart"
	"strconv"
)

type ResizeController struct {
}

func NewResizeController() *ResizeController {
	return &ResizeController{}
}

func (rc *ResizeController) Resize(c *gin.Context) {

	auth, err := aws.EnvAuth()

	if err != nil {
		log.Fatal(err)
	}

	client := s3.New(auth, aws.EUWest)
	bucket := client.Bucket("20.65twenty.com")

	file, header, err := c.Request.FormFile("file")
	filename := header.Filename

	if err != nil {
		log.Fatal(err)
	}

	height := c.PostForm("height")
	width := c.PostForm("width")

	fmt.Println("Height: " + height)
	fmt.Println("Width: " + width)

	finalFile := rc.Crop(height, width, file)

	err = bucket.Put(filename, finalFile, "image/jpeg", s3.BucketOwnerFull)
	if err != nil {
		panic(err)
	}

	c.JSON(200, gin.H{"filename": header.Filename})
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
