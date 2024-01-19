package main

import (
	"link-short/controller"
	"log"
	"os"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err !=nil{
		log.Println("no .env file found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable")
	}

	r := gin.Default()

	r.Use(static.Serve("/", static.LocalFile("./client/dist", true)))

	// TODO:
	r.POST("/create", controller.Create)
	r.GET("/:paramURL", controller.Redirect)

	r.Run(":" + port)
}
