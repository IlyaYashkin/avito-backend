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

	engine := gin.Default()
	engine.SetTrustedProxies(nil)

	engine.POST("/create-segment", segment.CreateSegment)
	engine.DELETE("/delete-segment", segment.DeleteSegment)
	engine.POST("/add-segments-to-user", usersegment.UpdateUserSegments)
	engine.GET("/get-user-segments/:user_id", usersegment.GetUserSegments)
	// engine.GET("/get-user-segment-log/*user_id", usersegment.GetUserSegments)

	engine.Run("0.0.0.0:8080")
}
