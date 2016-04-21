package handlers

import (
	"github.com/EwanValentine/Ice/services"
	"github.com/labstack/echo"
	"github.com/mitchellh/goamz/s3"
	"log"
	"path/filepath"
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
