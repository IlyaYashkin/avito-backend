package usersegmentlog

import (
	"avito-backend/internal/database"
	"encoding/csv"
	"fmt"

	"github.com/gin-gonic/gin"
)

type UserSegmentLog struct {
	UserId             int32
	SegmentName        string
	Operation          string
	OperationTimestamp string
}

const LOG_OPERATION_ADD = "added"
const LOG_OPERATION_ADD_PERCENTAGE = "added by percentage"
const LOG_OPERATION_DELETE = "deleted"
const LOG_OPERATION_DELETE_TTL = "deleted by ttl"

func getUserSegmentLog(requestData RequestGetUserSegmentLog, w gin.ResponseWriter) (*csv.Writer, error) {
	db := database.Get()
	writer := csv.NewWriter(w)

	userSegmentLogs, err := SelectLog(requestData.UserId, requestData.Date, db)
	if err != nil {
		return writer, err
	}

	for _, userSegmentLog := range userSegmentLogs {
		writer.Write(
			[]string{
				fmt.Sprint(userSegmentLog.UserId),
				userSegmentLog.SegmentName,
				userSegmentLog.Operation,
				userSegmentLog.OperationTimestamp,
			},
		)
	}

	return writer, nil
}
