package controllers

import (
	"avito-backend/models"
	"avito-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateSegment(c *gin.Context) {
	var requestData models.UpdateSegmentData

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.CreateSegment(requestData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": gin.H{"message": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"message": "Segment created", "name": requestData.Name}})
}
