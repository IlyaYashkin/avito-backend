package app

import (
	"avito-backend/internal/database"
	"avito-backend/internal/entity/segment"
	"avito-backend/internal/entity/usersegment"

	"github.com/gin-gonic/gin"
)

func Start() {
	db := database.Open()
	defer db.Close()

	r := gin.Default()
	initRoutes(r)
	r.Run("0.0.0.0:8080")
}

func initRoutes(r *gin.Engine) {
	r.POST("/create-segment", segment.CreateSegment)
	r.DELETE("/delete-segment", segment.DeleteSegment)
	r.POST("/add-segments-to-user", usersegment.UpdateUserSegments)
}
