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

	err := services.UpdateUserSegments(requestData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": gin.H{"message": err.Error()}})
		return
	}

}
