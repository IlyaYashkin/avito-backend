package app

import (
	"avito-backend/internal/database"
	"avito-backend/internal/entity/segment"
	"avito-backend/internal/entity/usersegment"
	"avito-backend/internal/entity/usersegmentlog"

	"github.com/gin-gonic/gin"
)

func Start() {
	db := database.Open()
	defer db.Close()

	engine := gin.Default()
	engine.SetTrustedProxies(nil)

	engine.POST("/create-segment", segment.CreateSegment)
	engine.DELETE("/delete-segment", segment.DeleteSegment)
	engine.POST("/update-user-segments", usersegment.UpdateUserSegments)
	engine.GET("/get-user-segments/:user_id", usersegment.GetUserSegments)
	engine.GET("/get-user-segment-log", usersegmentlog.GetUserSegmentLog)

	engine.Run("0.0.0.0:8080")
}
