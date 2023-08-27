package main

import (
	"avito-backend/database"
	"avito-backend/router"
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
	db := database.Open()
	defer db.Close()

	r := gin.Default()
	router.InitRoutes(r)
	r.Run("localhost:8080")
}
