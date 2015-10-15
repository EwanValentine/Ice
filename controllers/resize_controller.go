package controllers

import (
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

	auth, err := aws.EnvAuth()

	if err != nil {
		log.Fatal(err)
	}

	client := s3.New(auth, aws.EUWest)

	file, header, err := c.Request.FormFile("file")
	filename := header.Filename

	height := c.Query("height")
	width := c.Query("width")

	h64, err := strconv.ParseUint(height, 10, 32)
	w64, err := strconv.ParseUint(width, 10, 32)
	h := uint(h64)
	w := uint(w64)

	m := resize.Resize(w, h, file, resize.Lanczos3)

	client.Put(os.Getenv("AWS_BUCKET_NAME"), m)

	c.Data(200, "image/jpeg", response)
}
