package app

import (
	"avito-backend/internal/database"
	"avito-backend/internal/entity/segment"
	"avito-backend/internal/entity/usersegment"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Start() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	db := database.Open()
	defer db.Close()

	r := gin.Default()
	initRoutes(r)
	r.Run("localhost:8080")
}

func initRoutes(r *gin.Engine) {
	r.POST("/create-segment", segment.CreateSegment)
	r.DELETE("/delete-segment", segment.DeleteSegment)
	r.POST("/add-segments-to-user", usersegment.UpdateUserSegments)
}
