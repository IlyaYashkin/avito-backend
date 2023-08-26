package main

import (
	"avito-backend/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	r := gin.Default()
	routes.InitRoutes(r)
	r.Run("localhost:8080")
}
