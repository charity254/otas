package main

import (
	"log"
	"os"

	"otas/config"
	"otas/internal/routes"
	"otas/migrations"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using defaults")
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
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Unable to start server:", err)
	}
}
