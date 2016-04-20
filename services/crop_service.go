package services

import (
	"bytes"
	"encoding/base64"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"mime/multipart"
	"strings"
)

type CropService struct {
}

func NewCropService() *CropService {
	return &CropService{}
}

// Encodes base64 image, then calls standard crop function with base64 image now as a byte stream
func (cropper *CropService) CropBase64(height, width uint, file, ext string) []byte {

	// Remove meta prefix data
	b64data := file[strings.IndexByte(file, ',')+1:]

	// Create byte stream
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64data))

	image, _, err := image.Decode(reader)

	if err != nil {
		log.Println(err)
	}

	// Call crop function and return bytes
	return cropper.Crop(height, width, image, ext)
}

// Crops image and returns []byte of file - Converts `multipart.File` to `image.Image`
func (cropper *CropService) CropFile(height, width uint, file multipart.File, ext string) []byte {

	// Decode image
	image, _, err := image.Decode(file)

	if err != nil {
		log.Println(err)
	}

	// Call crop function and return bytes
	return cropper.Crop(height, width, image, ext)
}

// CropByte - AWS  returns files in byte format, so we need to decode this before passing
// it to the generic crop function, which expects `image.Image`
func (cropper *CropService) CropByte(height, width uint, file []byte, ext string) []byte {

	// Convert bytes to a `io.Reader` type, required by `image.Decode`
	reader := bytes.NewReader(file)

	// Decode image
	image, _, err := image.Decode(reader)

	if err != nil {
		log.Println(err)
	}

	return cropper.Crop(height, width, image, ext)
}

// Crop
// Generic Crop functions, take height, width, image file and file extension
func (cropper *CropService) Crop(height, width uint, image image.Image, ext string) []byte {

	var err error

	// Runs re-size function
	m := resize.Resize(width, height, image, resize.Lanczos3)

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
		log.Println(err)
	}

	// Return buffer
	return buf.Bytes()
}
