package main

import (
	"flag"
	"github.65twenty.com/65twenty/Ice/controllers"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"log"
	"os"
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

	var port = flag.String("port", "3000", "Enter a port number")
	flag.Parse()

	rc := controllers.NewResizeController(s3Config())

	r := gin.Default()
	r.Use(CorsMiddleware())

	// Routes
	r.POST("/resize", rc.PostResize)
	r.POST("/resize-base64", rc.PostBase64Resize)
	r.GET("/resize", rc.GetResize)

	r.Run(":" + *port)
}

// AWS s3 config
func s3Config() *s3.Bucket {
	auth, err := aws.EnvAuth()
	if err != nil {
		log.Fatal(err)
	}
	// @todo - Make region configurable
	client := s3.New(auth, aws.EUWest)
	bucketName := os.Getenv("AWS_BUCKET_NAME")
	if bucketName != "" {
		return client.Bucket(bucketName)
	}
	return client.Bucket("20.65twenty.com")
}
