package handlers

import (
	"bytes"
	"github.com/EwanValentine/Ice/services"
	"github.com/labstack/echo"
	"github.com/mitchellh/goamz/s3"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

// Mass crop type
type Mass struct {
	Files []File `json:"files"`
}

// Single file upload and crop type - with an array of dimensions to
// crop that single image to.
type SingleUpload struct {
	Filename   string      `json:"filename"`
	Dimensions []Dimension `json:"dimensions"`
	File       string      `json:"file"`
}

type Dimension struct {
	Height uint `json:"height"`
	Width  uint `json:"width"`
}

type File struct {
	Filename   string      `json:"filename"`
	Dimensions []Dimension `json:"dimensions"`
}

type Response struct {
	Data interface{}            `json:"data"`
	Meta map[string]interface{} `json:"_meta"`
}

type Error struct {
	Message string `json:"_message"`
	Code    int    `json:"code"`
}

// Controller type
type ResizeHandler struct {
	bucket   *s3.Bucket
	uploader *services.UploadService
	cropper  *services.CropService
}

// Controller instance
func NewResizeHandler(
	bucket *s3.Bucket,
	uploader *services.UploadService,
	cropper *services.CropService,
) *ResizeHandler {
	return &ResizeHandler{
		bucket,
		uploader,
		cropper,
	}
}

// GetRezize - On the fly cropping and resizing.
func (handler *ResizeHandler) GetResize(c echo.Context) error {

	// Set height and width
	width, _ := strconv.Atoi(c.QueryParam("width"))
	height, _ := strconv.Atoi(c.QueryParam("height"))

	// Get file
	filename := c.Param("file")

	// Get image
	file, err := handler.bucket.Get("content/" + filename)

	if err != nil {
		log.Println(err)
	}

	// Create new byte stream
	img := bytes.NewReader(file)

	// Decode image
	decoded, _, err := image.Decode(img)

	// Crop image
	cropped := resize.Resize(uint(width), uint(height), decoded, resize.Lanczos3)

	// Create new output buffer
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, cropped, nil)

	// Response
	response := buf.Bytes()

	// Return cropped image
	return c.File(string(response))
}

func (handler *ResizeHandler) PostResize(c echo.Context) error {
	var data Mass
	var files []string

	c.Bind(&data)

	log.Println(data)

	for i := 0; i < len(data.Files); i++ {

		// Get dimensions
		dimensions := data.Files[i].Dimensions

		// File
		filename := data.Files[i].Filename

		for d := 0; d < len(dimensions); d++ {

			// Get file extension
			ext := strings.Replace(filepath.Ext(filename), ".", "", -1)

			height := dimensions[d].Height
			width := dimensions[d].Width

			heightString := strconv.Itoa(int(height))
			widthString := strconv.Itoa(int(width))

			finalFilename := "h" + heightString + "w" + widthString + "-" + filename

			go func(height, width uint, filename, ext, finalFilename string) {

				// Everything from here onwards is done after the response :D
				original, err := handler.uploader.Get("content/" + filename)

				if err != nil {
					log.Println(err)
				}

				// Crop file
				file := handler.cropper.CropByte(
					height,
					width,
					original,
					ext,
				)

				// Upload file
				go handler.uploader.Upload("content/"+finalFilename, file, ext, s3.BucketOwnerFull)

			}(height, width, filename, ext, finalFilename)

			files = append(files, finalFilename)
		}
	}

	return c.JSON(200, &Response{
		Data: files,
		Meta: map[string]interface{}{
			"count": len(files),
		},
	})
}

// Upload a base64 image
// This allows you to post a single Base64 file with a set of
// dimensions, the image will be cropped to those dimensions,
// then uploaded to S3.
func (handler *ResizeHandler) PostBase64Resize(c echo.Context) error {
	var setData SingleUpload

	var files []string

	c.Bind(&setData)

	filename := setData.Filename

	// Get extension
	ext := strings.Replace(filepath.Ext(filename), ".", "", -1)

	file := setData.File

	// Foreach set of dimensions given
	for i := 0; i < len(setData.Dimensions); i++ {

		// Get height and width
		height := setData.Dimensions[i].Height
		width := setData.Dimensions[i].Width

		// Crop file
		finalFile := handler.cropper.CropBase64(height, width, file, ext)

		// Include dimensions in filename to stop file being overriden
		fileNameDimensions := "w" + string(width) + "h" + string(height) + "-" + filename

		// Upload file
		handler.uploader.Upload("content/"+fileNameDimensions, finalFile, "image/"+ext, s3.BucketOwnerFull)

		// Append file name to files list
		files = append(files, fileNameDimensions)
	}

	return c.JSON(200, &Response{
		Data: files,
		Meta: map[string]interface{}{},
	})
}
