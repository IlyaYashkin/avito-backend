package usersegment

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestUpdateUserSegments struct {
	UserId         int32    `json:"user_id" binding:"required"`
	AddSegments    []any    `json:"add_segments"`
	DeleteSegments []string `json:"delete_segments"`
}

func UpdateUserSegments(c *gin.Context) {
	var requestData RequestUpdateUserSegments
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": gin.H{"message": "bind json error"}})
		return
	}

	updatedUserSegments, err := updateUserSegments(requestData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": gin.H{"message": err.Error()}})
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": "success",
			"data": gin.H{
				"message":            "User segments updated",
				"added_segments":     updatedUserSegments.AddedSegments,
				"added_segments_ttl": updatedUserSegments.AddedSegmentsTtl,
				"deleted_segments":   updatedUserSegments.DeletedSegments,
			},
		},
	)
}
