package dbaccess

import (
	"avito-backend/database"
)

func InsertLog(userId int32, segments map[int32]string, operation string, ex QueryExecutor) error {
	sqlString, values := database.BuildUserSegmentLogInsertString(userId, segments, operation)
	_, err := ex.Exec(sqlString, values...)
	if err != nil {
		return err
	}
	return nil
}
