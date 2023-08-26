package controllers

import (
	"avito-backend/dtos"
	"avito-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateSegment(c *gin.Context) {
	var requestData dtos.UpdateSegment
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": gin.H{"message": "bind json error"}})
		return
	}

	err := services.CreateSegment(requestData)
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
	var requestData dtos.UpdateSegment
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": gin.H{"message": "bind json error"}})
		return
	}

	err := services.DeleteSegment(requestData)
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
