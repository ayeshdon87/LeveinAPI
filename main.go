package main

import (
	"log"
	"os"

	"github.com/ayeshdon87/LeveinAPI/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	var port = os.Getenv("PORT")

	if port == "" {
		port = "8081"
	}

	router := gin.New()
	router.Use(gin.Logger())

	log.Println("CALL SERVER")
	routes.AuthRoutes(router)

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"success": "Welcome to API Service v1.0"})
	})

	router.Run(":" + port)
}
