package usersegment

import (
	"avito-backend/internal/database"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

func SelectUserSegmentByUserId(userId int32, ex database.QueryExecutor) ([]UserSegment, error) {
	rows, err := ex.Query("select segment_id, ttl from user_segment where user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userSegments []UserSegment

	for rows.Next() {
		var segment_id int32
		var ttl sql.NullString
		err := rows.Scan(&segment_id, &ttl)
		if err != nil {
			return nil, err
		}
		userSegments = append(userSegments, UserSegment{SegmentId: segment_id, Ttl: ttl.String})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userSegments, nil
}

func SelectUserSegmentNamesByUserId(userId int32, ex database.QueryExecutor) ([]string, error) {
	rows, err := ex.Query( /* sql */ `
		select segments.name from user_segment
		left join segments
		on user_segment.segment_id = segments.id
		where user_id = $1
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var segmentsNames []string

	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		segmentsNames = append(segmentsNames, name)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return segmentsNames, nil
}

func InsertUserSegment(userId int32, segments []int32, ex database.QueryExecutor) error {
	sqlString, values := BuildUserSegmentInsertString(userId, segments)
	_, err := ex.Exec(sqlString, values...)
	if err != nil {
		return err
	}
	return nil
}

func InsertUserTtlSegment(userId int32, segments map[int32]string, ttls map[int32]time.Time, ex database.QueryExecutor) error {
	sqlString, values := BuildUserSegmentTtlInsertString(userId, segments, ttls)
	_, err := ex.Exec(sqlString, values...)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUserSegment(userId int32, segmentsIds []int32, ex database.QueryExecutor) error {
	sqlString := "delete from user_segment where user_id = $1 and segment_id = ANY($2)"
	_, err := ex.Exec(sqlString, userId, pq.Array(segmentsIds))
	if err != nil {
		return err
	}
	return nil
}

func deleteExpiredTtlUserSegments(ex database.QueryExecutor) ([]UserSegment, error) {
	sqlString := /* sql */ `
		delete from user_segment where ttl < now()
		returning user_id, segment_id
	`
	rows, err := ex.Query(sqlString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deletedSegments := []UserSegment{}
	for rows.Next() {
		var user_id int32
		var segment_id int32
		err := rows.Scan(&user_id, &segment_id)
		if err != nil {
			return nil, err
		}
		deletedSegments = append(deletedSegments, UserSegment{UserId: user_id, SegmentId: segment_id})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return deletedSegments, nil
}
