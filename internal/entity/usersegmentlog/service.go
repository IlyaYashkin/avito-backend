package usersegmentlog

import (
	"avito-backend/internal/database"
	"encoding/csv"
	"fmt"

	"github.com/gin-gonic/gin"
)

const LOG_OPERATION_ADD = "added"
const LOG_OPERATION_ADD_PERCENTAGE = "added by percentage"
const LOG_OPERATION_DELETE = "deleted"
const LOG_OPERATION_DELETE_TTL = "deleted by ttl"

func GetUserSegmentLogService(requestData RequestGetUserSegmentLog, w gin.ResponseWriter) (*csv.Writer, error) {
	db := database.Get()
	writer := csv.NewWriter(w)

	userSegmentLogRepo := NewUserSegmentLogRepository(db)

	userSegmentLogs, err := userSegmentLogRepo.Get(requestData.UserId, requestData.Date)
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
