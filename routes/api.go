package routes

import (
	"avito-backend/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/create-segment", controllers.CreateSegment)
}
