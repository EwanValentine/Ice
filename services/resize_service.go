package services

import (
	"github.com/nfnt/resize"
	"image"
)

type ResizeService struct {
	decoder *DecoderService
}

func NewResizeService(decoder *DecoderService) *ResizeService {
	return &ResizeService{
		decoder,
	}
}

// Resize
// Generic Resize function, take height, width, image file and file extension
func (resizer *ResizeService) Resize(height, width uint, image image.Image, ext string) ([]byte, error) {

	// Runs re-size function
	file := resize.Resize(width, height, image, resize.Lanczos3)

	// Return buffer
	return resizer.decoder.EncodeImage(file, ext)
}
