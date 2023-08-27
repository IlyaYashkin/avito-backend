package controller

import (
	"avito-backend/dto"
	"avito-backend/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateSegment(c *gin.Context) {
	var requestData dto.UpdateSegment
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": gin.H{"message": "bind json error"}})
		return
	}

	err := service.CreateSegment(requestData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": gin.H{"message": err.Error()}})
		return
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"status": "success",
			"data": gin.H{
				"message": "Segment created",
				"name":    requestData.Name,
			},
		},
	)
}

func DeleteSegment(c *gin.Context) {
	var requestData dto.UpdateSegment
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": gin.H{"message": "bind json error"}})
		return
	}

	err := service.DeleteSegment(requestData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": gin.H{"message": err.Error()}})
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": "success",
			"data": gin.H{
				"message": "Segment deleted",
				"name":    requestData.Name,
			},
		},
	)
}
