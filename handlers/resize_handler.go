package handlers

import (
	"bytes"
	"github.com/EwanValentine/Ice/services"
	"github.com/labstack/echo"
	"github.com/mitchellh/goamz/s3"
	"image"
	_ "image/png"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

// Controller type
type ResizeHandler struct {
	bucket   *s3.Bucket
	uploader *services.UploadService
	resizer  *services.ResizeService
	decoder  *services.DecoderService
}

// Controller instance
func NewResizeHandler(
	bucket *s3.Bucket,
	uploader *services.UploadService,
	resizer *services.ResizeService,
	decoder *services.DecoderService,
) *ResizeHandler {
	return &ResizeHandler{
		bucket,
		uploader,
		resizer,
		decoder,
	}
}

// GetRezize - On the fly cropping and resizing.
// Crops images based on url parameters. For example...
// /resize?file=123.jpg&width=34&height=56
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

	// Get file ext
	ext := strings.Replace(filepath.Ext(filename), ".", "", -1)

	// Decode image
	decoded, _, err := image.Decode(img)

	// Resize image
	cropped, err := handler.resizer.Resize(uint(width), uint(height), decoded, ext)

	// Return cropped image
	return c.File(string(cropped))
}

// PostResize
// Resizes images stored in S3 based on filename, against given
// dimensions.
// For example...
// POST { "files": [ { "filename": "123.jpg", "dimensions: [ { "width": 50, "height": 50 } ] } ] }
func (handler *ResizeHandler) PostResize(c echo.Context) error {
	var data Mass
	var files []string

	c.Bind(&data)

	// Foreach file
	for i := 0; i < len(data.Files); i++ {

		// Get dimensions
		dimensions := data.Files[i].Dimensions

		// File
		filename := data.Files[i].Filename

		// Foreach dimension
		for d := 0; d < len(dimensions); d++ {

			// Get file extension
			ext := strings.Replace(filepath.Ext(filename), ".", "", -1)

			height := dimensions[d].Height
			width := dimensions[d].Width

			finalFilename := GenerateFilename(height, width, filename)

			// Handle upload and crop process in the background as we don't need to wait for this
			go func(height, width uint, filename, ext, finalFilename string) {

				// Everything from here onwards is done after the response :D
				original, err := handler.uploader.Get("content/" + filename)

				if err != nil {
					log.Println(err)
				}

				// Decode image
				decodedImage, _ := handler.decoder.DecodeBytes(original)

				// Crop file
				file, _ := handler.resizer.Resize(
					height,
					width,
					decodedImage,
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

		// Convert Base64 string to `image.Image`
		decodedImage, _ := handler.decoder.DecodeBase64(file)

		// Crop file
		finalFile, _ := handler.resizer.Resize(height, width, decodedImage, ext)

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
