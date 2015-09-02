package controllers

import (
	"github.com/gin-gonic/gin"
)

type ResizeController struct {
}

func NewResizeController() *ResizeController {
	return &ResizeController{}
}

func (rc ResizeController) Resize(c *gin.Context) {

	height := c.Query("height")
	width := c.Query("width")

	var Result struct {
		Width  string `json:"width"`
		Height string `json:"height"`
	}

	Result.Height = height
	Result.Width = width

	c.JSON(200, Result)
}
