package main

import (
	"bitbucket.org/65twenty/ice/controllers"
	"github.com/gin-gonic/contrib/jwt"
	"github.com/gin-gonic/gin"
	"time"
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

	r.GET("/resize", rc.Resize)

	r.Run(":8000")
}
