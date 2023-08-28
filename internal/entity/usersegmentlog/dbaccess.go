package usersegmentlog

import (
	"avito-backend/internal/database"
)

func InsertLog(userId int32, segments map[int32]string, operation string, ex database.QueryExecutor) error {
	sqlString, values := BuildUserSegmentLogInsertString(userId, segments, operation)
	_, err := ex.Exec(sqlString, values...)
	if err != nil {
		return err
	}
	return nil
}
