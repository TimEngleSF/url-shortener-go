package main

import (
	"link-short/controller"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(static.Serve("/", static.LocalFile("./client/dist", true)))

	// TODO:
	r.POST("/create", controller.Create)
	r.GET("/:paramURL", controller.Redirect)

	r.Run(":8000")
}
