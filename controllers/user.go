package controllers

import (
	"avito-backend/dtos"
	"avito-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateUserSegments(c *gin.Context) {
	var requestData dtos.UpdateUserSegments
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": gin.H{"message": "bind json error"}})
		return
	}

	_, _, err := services.UpdateUserSegments(requestData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": gin.H{"message": err.Error()}})
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": "success",
			"data": gin.H{
				"message": "User segments updated",
			},
		},
	)
}
