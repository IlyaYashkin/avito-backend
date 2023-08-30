package usersegmentlog

import (
	"avito-backend/internal/database"
)

func InsertLog(rows []UserSegmentLog, operation string, ex database.QueryExecutor) error {
	sqlString, values := BuildUserSegmentLogInsertString123(rows, operation)
	_, err := ex.Exec(sqlString, values...)
	if err != nil {
		return err
	}
	return nil
}
