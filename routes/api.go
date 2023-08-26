package routes

import (
	"avito-backend/controllers"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	r.POST("/create-segment", controllers.CreateSegment)
	r.DELETE("/delete-segment", controllers.DeleteSegment)
	r.POST("/add-segments-to-user", controllers.UpdateUserSegments)
}
