package router

import (
	"avito-backend/controller"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	r.POST("/create-segment", controller.CreateSegment)
	r.DELETE("/delete-segment", controller.DeleteSegment)
	r.POST("/add-segments-to-user", controller.UpdateUserSegments)
}
