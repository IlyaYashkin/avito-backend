package usersegment

import (
	"avito-backend/internal/database"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

func GetMatchedSegments(segments []string, ex database.QueryExecutor) (map[int32]string, error) {
	rows, err := ex.Query("select id, name from segments where name = ANY($1)", pq.Array(segments))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	segmentsMap := make(map[int32]string)

	for rows.Next() {
		var id int32
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			return segmentsMap, err
		}
		segmentsMap[id] = name
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return segmentsMap, nil
}

func GetUsrSegments(userId int32, ex database.QueryExecutor) ([]UserSegment, error) {
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

func InsUsrSeg(userId int32, segments []int32, ex database.QueryExecutor) error {
	sqlString, values := BuildUserSegmentInsertString(userId, segments)
	_, err := ex.Exec(sqlString, values...)
	if err != nil {
		return err
	}
	return nil
}

func InsUsrTtlSeg(userId int32, segments map[int32]string, ttls map[int32]time.Time, ex database.QueryExecutor) error {
	sqlString, values := BuildUserSegmentTtlInsertString(userId, segments, ttls)
	_, err := ex.Exec(sqlString, values...)
	if err != nil {
		return err
	}
	return nil
}

func DelUsrSeg(userId int32, segmentsIds []int32, ex database.QueryExecutor) error {
	sqlString := "delete from user_segment where user_id = $1 and segment_id = ANY($2)"
	_, err := ex.Exec(sqlString, userId, pq.Array(segmentsIds))
	if err != nil {
		return err
	}
	return nil
}