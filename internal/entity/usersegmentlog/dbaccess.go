package usersegmentlog

import (
	"avito-backend/internal/database"
	"time"
)

func InsertLog(rows []UserSegmentLog, operation string, ex database.QueryExecutor) error {
	sqlString, values := BuildUserSegmentLogInsertString(rows, operation)
	_, err := ex.Exec(sqlString, values...)
	if err != nil {
		return err
	}
	return nil
}

func SelectLog(userId int32, date time.Time, ex database.QueryExecutor) ([]UserSegmentLog, error) {
	sqlString, values := BuildUserSegmentLogSelectString(userId, date)
	rows, err := ex.Query(sqlString, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userSegmentLogs []UserSegmentLog
	for rows.Next() {
		var user_id int32
		var segment_name string
		var operation string
		var operation_timestamp string
		err := rows.Scan(&user_id, &segment_name, &operation, &operation_timestamp)
		if err != nil {
			return nil, err
		}
		userSegmentLogs = append(userSegmentLogs, UserSegmentLog{
			UserId:             user_id,
			SegmentName:        segment_name,
			Operation:          operation,
			OperationTimestamp: operation_timestamp,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userSegmentLogs, nil
}
