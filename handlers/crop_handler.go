package handlers

import (
	"github.com/EwanValentine/Ice/services"
	"github.com/labstack/echo"
	"github.com/mitchellh/goamz/s3"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type CropHandler struct {
	bucket   *s3.Bucket
	uploader *services.UploadService
	cropper  *services.CropService
	decoder  *services.DecoderService
}

func NewCropHandler(
	bucket *s3.Bucket,
	uploader *services.UploadService,
	cropper *services.CropService,
	decoder *services.DecoderService,
) *CropHandler {
	return &CropHandler{
		bucket,
		uploader,
		cropper,
		decoder,
	}
}

// GetCrop - Fetch and crop an image on the fly
func (handler *CropHandler) GetCrop(c echo.Context) error {

	// Get height and width from url parameters
	width, _ := strconv.Atoi(c.QueryParam("width"))
	height, _ := strconv.Atoi(c.QueryParam("height"))

	// Get file parameter
	filename := c.QueryParam("file")
	ext := GetExtension(filename)

	// If no file given, not much we can do
	if filename == "" {
		return c.JSON(http.StatusNotFound, &Error{
			Message: "No file given",
			Code:    http.StatusNotFound,
		})
	}

	// Fetch image from S3
	file, err := handler.bucket.Get("content/" + filename)

	// If file can't be fetched, throw a 404
	if err != nil {
		return c.JSON(http.StatusNotFound, &Error{
			Message: "File not found in bucket",
			Code:    http.StatusNotFound,
		})
	}

	// Decode image from bytes to `image.Image`
	img, err := handler.decoder.DecodeBytes(file)

	// Crop image
	cropped, err := handler.cropper.Crop(uint(width), uint(height), img, ext)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Error{
			Message: err,
			Code:    http.StatusInternalServerError,
		})
	}

	// Return cropped image
	return c.File(string(cropped))
}

// PostBase64Crop - Crops a base64 image to various dimensions
// then uploads them to S3
func (handler *CropHandler) PostBase64Crop(c echo.Context) error {
	var setData SingleUpload

	var files []string

	c.Bind(&setData)

	filename := setData.Filename

	// Get ext
	ext := GetExtension(filename)

	file := setData.File

	// Foreach set of dimensions given
	for i := 0; i < len(setData.Dimensions); i++ {

		// Get height and width
		height := setData.Dimensions[i].Height
		width := setData.Dimensions[i].Width

		// Include dimensions in filename to stop file being overriden
		fileNameDimensions := GenerateDimensionFilename(string(height), string(width), filename, "crop")

		go func(file string, height, width uint, fileNameDimensions, ext string) {

			// Decode Base64 image to `image.Image`
			decodedImage, _ := handler.decoder.DecodeBase64(file)

			// Crop image
			finalFile, _ := handler.cropper.Crop(height, width, decodedImage, ext)

			handler.uploader.Upload("content/"+fileNameDimensions, finalFile, "image/"+ext, s3.BucketOwnerFull)
		}(file, height, width, fileNameDimensions, ext)

		files = append(files, fileNameDimensions)
	}

	return c.JSON(http.StatusOK, &Response{
		Data: files,
		Meta: map[string]interface{}{},
	})
}

// PostCrop
// Resizes images stored in S3 based on filename, against given
// dimensions.
// For example...
// POST { "files": [ { "filename": "123.jpg", "dimensions: [ { "width": 50, "height": 50 } ] } ] }
func (handler *CropHandler) PostCrop(c echo.Context) error {
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

				decodedImage, _ := handler.decoder.DecodeBytes(original)

				// Crop file
				file, cropErr := handler.cropper.Crop(
					height,
					width,
					decodedImage,
					ext,
				)

				if cropErr != nil {
					log.Println(cropErr)
				}

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
