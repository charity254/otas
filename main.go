package main

import (
	"log"
	"os"

	"otas/config"
	"otas/migrations"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}

	config.ConnectDB()
	migrations.Run()

	r := gin.Default()
	routes.Register(r, config.DB)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
