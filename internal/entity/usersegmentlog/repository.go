package usersegmentlog

import (
	"avito-backend/internal/database"
	"time"
)

type UserSegmentLog struct {
	UserId             int32
	SegmentName        string
	Operation          string
	OperationTimestamp string
}

type SegmentRepository interface {
	Save(rows []UserSegmentLog, operation string) error
	Get(userId int32, date time.Time) ([]UserSegmentLog, error)
}

type UserSegmentLogRepositoryDB struct {
	ex database.QueryExecutor
}

func NewUserSegmentLogRepository(ex database.QueryExecutor) *UserSegmentLogRepositoryDB {
	return &UserSegmentLogRepositoryDB{ex: ex}
}

func (repo UserSegmentLogRepositoryDB) Save(rows []UserSegmentLog, operation string) error {
	sqlString, values := buildUserSegmentLogInsertString(rows, operation)
	_, err := repo.ex.Exec(sqlString, values...)
	if err != nil {
		return err
	}
	return nil
}

func (repo UserSegmentLogRepositoryDB) Get(userId int32, date time.Time) ([]UserSegmentLog, error) {
	sqlString, values := buildUserSegmentLogSelectString(userId, date)
	rows, err := repo.ex.Query(sqlString, values...)
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
