package handlers

import (
	"strconv"
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

// GenerateFilename
// Take dimensions, original filename, converts dimensions to
// integer and forms the final filename including the dimensions.
func GenerateFilename(height, width uint, filename string) string {

	// Convert dimensions to string
	heightString := strconv.Itoa(int(height))
	widthString := strconv.Itoa(int(width))

	// Final file name, i.e `h50w50-original-filename.jpg`
	return "h" + heightString + "w" + widthString + "-" + filename
}
