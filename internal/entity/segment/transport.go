package segment

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestUpdateSegment struct {
	Name           string  `json:"name" binding:"required"`
	UserPercentage float32 `json:"user_percentage"`
}

func CreateSegment(c *gin.Context) {
	var requestData RequestUpdateSegment
	if err := c.ShouldBindJSON(&requestData); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": gin.H{"message": "bind json error"}})
		return
	}

	err := CreateSegmentService(requestData)
	if err != nil {
		log.Println(err)
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
	var requestData RequestUpdateSegment
	if err := c.ShouldBindJSON(&requestData); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": gin.H{"message": "bind json error"}})
		return
	}

	err := DeleteSegmentService(requestData)
	if err != nil {
		log.Println(err)
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
