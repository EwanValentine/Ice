package services

import (
	"github.com/oliamb/cutter"
	"image"
)

type CropService struct {
	decoder *DecoderService
}

func NewCropService(decoder *DecoderService) *CropService {
	return &CropService{
		decoder,
	}
}

// Crop
// Crops an image based on the given dimensions
func (cropper *CropService) Crop(height, width uint, image image.Image, ext string) ([]byte, error) {

	// Crop image
	file, err := cutter.Crop(image, cutter.Config{
		Width:   int(width),
		Height:  int(height),
		Options: cutter.Copy,
	})

	if err != nil {
		return nil, err
	}

	return cropper.decoder.EncodeImage(file, ext)
}
