package usersegmentlog

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type RequestGetUserSegmentLog struct {
	UserId int32
	Date   time.Time
}

func GetUserSegmentLog(c *gin.Context) {
	c.Header("Content-Disposition", "inline; filename=\"user-segment-log.csv\"") // Заменили "attachment" на "inline"
	c.Header("Content-Type", "text/csv")

	userIdParam := c.Query("user_id")
	dateParam := c.Query("date")

	requestData := RequestGetUserSegmentLog{}

	if userIdParam != "" {
		userId, err := strconv.ParseInt(userIdParam, 10, 32)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": gin.H{"message": err.Error()}})
			return
		}
		requestData.UserId = int32(userId)
	}
	if dateParam != "" {
		date, err := time.Parse("2006-01", dateParam)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": gin.H{"message": err.Error()}})
			return
		}
		requestData.Date = date.AddDate(0, 1, -1)
	}

	log.Println(RequestGetUserSegmentLog{}.Date.IsZero())

	writer, err := getUserSegmentLog(requestData, c.Writer)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": gin.H{"message": err.Error()}})
		return
	}

	writer.Flush()
}
