package services

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"image/png"
	_ "image/png"
	"mime/multipart"
	"strings"
)

type DecoderService struct{}

func NewDecoderService() *DecoderService {
	return &DecoderService{}
}

// DecodeBase64
// Decodes a base64 string into an `image.Image` type
func (decoder *DecoderService) DecodeBase64(file string) (image.Image, error) {

	// Remove meta prefix data
	b64data := file[strings.IndexByte(file, ',')+1:]

	// Create byte stream
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64data))

	// Decode image
	image, _, err := image.Decode(reader)

	return image, err
}

// DecodeFile
// Decodes `multipart.File` to type `image.Image`
func (decoder *DecoderService) DecodeFile(file multipart.File) (image.Image, error) {

	// Decode file to type `image.Image`
	image, _, err := image.Decode(file)

	return image, err
}

// DecodeBytes
// Decodes `bytes` to `image.Image` - S3 returns image data as bytes
func (decoder *DecoderService) DecodeBytes(file []byte) (image.Image, error) {

	// Convert bytes to a `io.Reader` type, required by `image.Decode`
	reader := bytes.NewReader(file)

	// Decode image
	image, _, err := image.Decode(reader)

	return image, err
}

// EncodeImage - Encodes image back into original format
func (decoder *DecoderService) EncodeImage(file image.Image, ext string) ([]byte, error) {

	var err error

	buffer := new(bytes.Buffer)

	switch {
	case "jpg" == ext || "jpeg" == ext:
		err = jpeg.Encode(buffer, file, nil)
	case "png" == ext:
		err = png.Encode(buffer, file)
	}

	return buffer.Bytes(), err
}
