package main

import (
	"flag"
	"github.com/EwanValentine/Ice/drivers"
	"github.com/EwanValentine/Ice/handlers"
	"github.com/EwanValentine/Ice/services"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	"github.com/labstack/echo/middleware"
	"log"
	"runtime"
)

func Init() {

	// Verbose loggin
	log.SetFlags(log.Lshortfile)

	// Use all available CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	// Initialise application runtime settings
	Init()

	// Create new Echo instance
	e := echo.New()

	// Apply middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Get port number from flags
	var port = flag.String("port", "3000", "Enter a port number")
	var bucketName = flag.String("bucket", "main", "Enter a bucket name")

	flag.Parse()

	// Drivers
	bucket := drivers.GetBucket(*bucketName)

	// Services
	uploader := services.NewUploadService(bucket)
	decoder := services.NewDecoderService()
	resizer := services.NewResizeService(decoder)
	cropper := services.NewCropService(decoder)

	// Handlers
	resizeHandler := handlers.NewResizeHandler(
		bucket,
		uploader,
		resizer,
		decoder,
	)

	cropHandler := handlers.NewCropHandler(
		bucket,
		uploader,
		cropper,
		decoder,
	)

	// Routes
	e.Post("/resize", resizeHandler.PostResize)
	e.Post("/resize-base64", resizeHandler.PostBase64Resize)
	e.Get("/resize", resizeHandler.GetResize)

	e.Post("/crop", cropHandler.PostCrop)

	// Run new fasthttp server instance
	e.Run(fasthttp.New(":" + *port))
}
