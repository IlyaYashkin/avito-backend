package usersegment

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RequestUpdateUserSegments struct {
	UserId         int32    `json:"user_id" binding:"required"`
	AddSegments    []any    `json:"add_segments"`
	DeleteSegments []string `json:"delete_segments"`
}

type RequestGetUserSegments struct {
	UserId int32 `json:"user_id" binding:"required"`
}

func UpdateUserSegments(c *gin.Context) {
	var requestData RequestUpdateUserSegments
	if err := c.ShouldBindJSON(&requestData); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": gin.H{"message": "bind json error"}})
		return
	}

	updatedUserSegments, err := updateUserSegments(requestData)
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
				"message":                   "User segments updated",
				"added_segments":            updatedUserSegments.AddedSegments,
				"added_ttl_segments":        updatedUserSegments.AddedTtlSegments,
				"added_percentage_segments": updatedUserSegments.AddedPercentageSegments,
				"deleted_segments":          updatedUserSegments.DeletedSegments,
			},
		},
	)
}

func GetUserSegments(c *gin.Context) {
	userIdParam := c.Param("user_id")
	if userIdParam == "" {
		c.JSON(400, gin.H{
			"status": "error",
			"data": gin.H{
				"message": "Parameter user_id is required",
			},
		})
		return
	}
	userId, err := strconv.ParseInt(userIdParam, 10, 32)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": gin.H{"message": err.Error()}})
		return
	}

	userSegments, err := getUserSegments(int32(userId))
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
				"user":     int32(userId),
				"segments": userSegments,
			},
		},
	)
}
