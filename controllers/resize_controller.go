package controllers

import (
	"bitbucket.org/65twenty/ice/services"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"image/jpeg"
	"log"
	"os"
	"strconv"
)

type ResizeController struct {
}

func NewResizeController() *ResizeController {
	return &ResizeController{}
}

func (rc ResizeController) Resize(c *gin.Context) {

	height := c.Query("height")
	width := c.Query("width")
	filepath := c.Query("file")

	h64, err := strconv.ParseUint(height, 10, 32)
	w64, err := strconv.ParseUint(width, 10, 32)
	h := uint(h64)
	w := uint(w64)

	file, err := os.Open("./test_images/" + filepath)

	if err != nil {
		log.Fatal(err)
	}

	image, err := jpeg.Decode(file)

	if err != nil {
		log.Fatal(err)
	}

	m := resize.Resize(w, h, image, resize.Lanczos3)

	buf := new(bytes.Buffer)
	jpeg.Encode(buf, m, nil)
	response := buf.Bytes()

	c.Data(200, "image/jpeg", response)
}

func (rc ResizeController) FromS3(c *gin.Context) {

	height := c.Query("height")
	width := c.Query("width")
	filepath := c.Query("file")

	h64, err := strconv.ParseUint(height, 10, 32)
	w64, err := strconv.ParseUint(width, 10, 32)
	h := uint(h64)
	w := uint(w64)

	s3 := services.S3("submissions")

	file, err := s3.Get(filepath)

	image, err := jpeg.Decode(file)

	if err != nil {
		log.Fatal(err)
	}

	m := resize.Resize(w, h, image, resize.Lanczos3)

	buf := new(bytes.Buffer)
	jpeg.Encode(buf, m, nil)
	response := buf.Bytes()

	c.Data(200, "image/jpeg", response)
}
