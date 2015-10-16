package main

import (
	"github.65twenty.com/65twenty/Ice/controllers"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.Abort()
			return
		}

		c.Next()
	}
}

func main() {
	rc := controllers.NewResizeController()

	r := gin.Default()
	r.Use(CorsMiddleware())

	r.POST("/resize", rc.Resize)

	r.Run(":3000")
}
