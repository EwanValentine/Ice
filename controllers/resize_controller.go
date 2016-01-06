package controllers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/goamz/s3"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	_ "image/png"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
)

// Form struct
type Form struct {
	Width    []string `form:"width[]" json:"width[]" binding:"required"`
	Height   []string `form:"height[]" json:"height[]"`
	File     string   `form:"file" json:"file"`
	Filename string   `form:"filename" json:"filename"`
}

type Files struct {
	Filename []string
}

// Controller type
type ResizeController struct {
	bucket *s3.Bucket
}

// Controller instance
func NewResizeController(bucket *s3.Bucket) *ResizeController {
	return &ResizeController{
		bucket,
	}
}

// GetRezize - On the fly cropping and resizing.
func (rc *ResizeController) GetResize(c *gin.Context) {

	// Set height and width
	width := c.Query("width")
	height := c.Query("height")

	// Get file
	filename := c.Query("file")

	// Get image
	file, err := rc.bucket.Get("content/" + filename)

	if err != nil {
		panic(err)
	}

	// Create new byte stream
	img := bytes.NewReader(file)

	// Decode image
	decoded, _, err := image.Decode(img)

	// Convert string height and width into 64int
	w64, err := strconv.ParseUint(width, 10, 32)
	h64, err := strconv.ParseUint(height, 10, 32)

	// Convert 64int to 32int
	h := uint(h64)
	w := uint(w64)

	// Crop image
	cropped := resize.Resize(w, h, decoded, resize.Lanczos3)

	// Create new output buffer
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, cropped, nil)

	// Response
	response := buf.Bytes()

	// Return cropped image
	c.Data(200, "image/jpeg", response)
}

func (rc *ResizeController) PostBase64Resize(c *gin.Context) {
	setData := &Form{}

	var files []string

	c.Bind(setData)

	if setData == nil {
		c.BindJSON(setData)
	}

	fmt.Println(setData.Filename)

	filename := setData.Filename

	// Get extension
	ext := strings.Replace(filepath.Ext(filename), ".", "", -1)

	file := setData.File

	// Foreach set of dimensions given
	for i := 0; i < len(setData.Width); i++ {

		// Get height and width
		height := setData.Height[i]
		width := setData.Width[i]

		// Crop file
		finalFile := rc.CropBase64(height, width, file, ext)

		// Include dimensions in filename to stop file being overriden
		fileNameDimensions := "w" + width + "h" + height + "-" + filename

		// Upload file
		rc.Upload("content/"+fileNameDimensions, finalFile, "image/"+ext, s3.BucketOwnerFull)

		// Append file name to files list
		files = append(files, fileNameDimensions)
	}

	c.JSON(200, gin.H{"files": files})

}

// PostResize - function for taking images and meta data and resizing
func (rc *ResizeController) PostResize(c *gin.Context) {

	// Empty formdata struct
	setData := &Form{}

	// Empty files struct
	var files []string

	// Bind empty struct to context
	c.Bind(setData)

	if setData == nil {
		c.BindJSON(&setData)
	}

	// Get file
	file, header, _ := c.Request.FormFile("file")

	var filename string

	if file == nil {

		fmt.Println("File missing")

		var formFiles []*multipart.FileHeader

		form := c.Request.MultipartForm
		formFiles = form.File["file[]"]

		if formFiles == nil {
			formFiles = form.File["file"]
		}

		fileP, _ := formFiles[0].Open()
		defer fileP.Close()
		file = fileP

		filename = formFiles[0].Filename

	} else {

		// Get original filename
		filename = header.Filename
	}

	// Get extension
	ext := strings.Replace(filepath.Ext(filename), ".", "", -1)

	// Foreach set of dimensions given
	for i := 0; i < len(setData.Width); i++ {

		// Get height and width
		height := setData.Height[i]
		width := setData.Width[i]

		// Crop file
		finalFile := rc.Crop(height, width, file, ext)

		// Include dimensions in filename to stop file being overriden
		fileNameDimensions := "w" + width + "h" + height + "-" + filename

		// Upload file
		rc.Upload("content/"+fileNameDimensions, finalFile, "image/"+ext, s3.BucketOwnerFull)

		// Append file name to files list
		files = append(files, fileNameDimensions)

		// Seek file back to first byte, so it can be re-cropped
		if _, err := file.Seek(0, 0); err != nil {
			panic(err)
		}
	}

	c.JSON(200, gin.H{"files": files})
}

// Uploads file to S3
func (rc *ResizeController) Upload(filename string, file []byte, enctype string, acl s3.ACL) {

	// Made this into a goroutine, knocked 200ms off!!
	go func(filename string, file []byte, enctype string, acl s3.ACL) {
		err := rc.bucket.Put(filename, file, enctype, acl)
		if err != nil {
			panic(err)
		}
	}(filename, file, enctype, acl)
}

func (rc *ResizeController) CropBase64(height string, width string, file string, ext string) []byte {

	b64data := file[strings.IndexByte(file, ',')+1:]
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64data))

	image, _, err := image.Decode(reader)

	if err != nil {
		panic(err)
	}

	w64, err := strconv.ParseUint(width, 10, 32)
	h64, err := strconv.ParseUint(height, 10, 32)

	h := uint(h64)
	w := uint(w64)

	m := resize.Resize(w, h, image, resize.Lanczos3)

	buf := new(bytes.Buffer)

	switch {
	case "jpg" == ext || "jpeg" == ext:
		err = jpeg.Encode(buf, m, nil)
	case "png" == ext:
		err = png.Encode(buf, m)
	}

	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

// Crops image and returns []byte of file
func (rc *ResizeController) Crop(height string, width string, file multipart.File, ext string) []byte {

	// Decode image
	image, _, err := image.Decode(file)

	// If error, panic
	if err != nil {
		panic(err)
	}

	// Convert string height and width into 64int
	w64, err := strconv.ParseUint(width, 10, 32)
	h64, err := strconv.ParseUint(height, 10, 32)

	// Convert 64int to 32int
	h := uint(h64)
	w := uint(w64)

	// Runs re-size function
	m := resize.Resize(w, h, image, resize.Lanczos3)

	// Create new buffer of file
	buf := new(bytes.Buffer)

	// Use correct encoder for file type
	switch {
	case "jpg" == ext || "jpeg" == ext:
		err = jpeg.Encode(buf, m, nil)
	case "png" == ext:
		err = png.Encode(buf, m)
	}

	if err != nil {
		panic(err)
	}

	// Return buffer
	return buf.Bytes()
}
