package services

import (
	"avito-backend/database"
	"avito-backend/dtos"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
)

func UpdateUserSegments(requestData dtos.UpdateUserSegments) (map[int32]string, map[int32]string, error) {
	db := database.Open()
	tx, err := db.Begin()
	if err != nil {
		return nil, nil, errors.New("begin transaction error")
	}
	defer tx.Rollback()

	err = createUserUsingTransaction(requestData.UserId, tx)
	if err != nil {
		return nil, nil, err
	}

	segmentsAdded, err := addUserSegments(
		requestData.UserId,
		requestData.AddSegments,
		tx,
	)
	if err != nil {
		return nil, nil, err
	}

	segmentsDeleted, err := deleteUserSegments(
		requestData.UserId,
		requestData.DeleteSegments,
		tx,
	)
	if err != nil {
		return nil, nil, err
	}

	tx.Commit()
	return segmentsAdded, segmentsDeleted, nil
}

func addUserSegments(userId int32, segmentsToAdd []string, tx *sql.Tx) (map[int32]string, error) {
	segments, err := getSegments(segmentsToAdd, tx)
	if err != nil {
		return nil, err
	}
	userSegmentsIds, err := getUserSegmentsIds(userId, tx)
	if err != nil {
		return nil, err
	}

	for _, segmentId := range userSegmentsIds {
		if _, exists := segments[segmentId]; exists {
			delete(segments, segmentId)
		}
	}
	if len(segments) == 0 {
		return segments, nil
	}

	sqlString, values := buildUserSegmentInsertString(userId, segments)
	_, err = tx.Exec(sqlString, values...)
	if err != nil {
		return nil, err
	}

	err = addInfoToLog(userId, segments, "addition", tx)
	if err != nil {
		return nil, err
	}

	return segments, nil
}

func deleteUserSegments(userId int32, segmentsToDelete []string, tx *sql.Tx) (map[int32]string, error) {
	segments, err := getSegments(segmentsToDelete, tx)
	if err != nil {
		return nil, err
	}
	userSegmentsIds, err := getUserSegmentsIds(userId, tx)
	if err != nil {
		return nil, err
	}

	for _, segmentId := range userSegmentsIds {
		if _, exists := segments[segmentId]; !exists {
			delete(segments, segmentId)
		}
	}
	if len(segments) == 0 {
		return segments, nil
	}

	sqlString := "delete from user_segment where segment_id = ANY($1)"
	var segmentsIdsToDelete []int32
	for i := range segments {
		segmentsIdsToDelete = append(segmentsIdsToDelete, i)
	}
	_, err = tx.Exec(sqlString, pq.Array(segmentsIdsToDelete))
	if err != nil {
		return nil, err
	}

	err = addInfoToLog(userId, segments, "deletion", tx)
	if err != nil {
		return nil, err
	}

	return segments, nil
}

func getSegments(segments []string, tx *sql.Tx) (map[int32]string, error) {
	rows, err := tx.Query("select id, name from segments where name = ANY($1)", pq.Array(segments))
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

func getUserSegmentsIds(userId int32, tx *sql.Tx) ([]int32, error) {
	rows, err := tx.Query("select segment_id from user_segment where user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userSegmentsIds []int32

	for rows.Next() {
		var id int32
		err := rows.Scan(&id)
		if err != nil {
			return userSegmentsIds, err
		}
		userSegmentsIds = append(userSegmentsIds, id)
	}

	if err = rows.Err(); err != nil {
		return userSegmentsIds, err
	}

	return userSegmentsIds, nil
}

func addInfoToLog(userId int32, segments map[int32]string, operation string, tx *sql.Tx) error {
	sqlString, values := buildUserSegmentLogInsertString(userId, segments, operation)

	_, err := tx.Exec(sqlString, values...)
	if err != nil {
		return err
	}

	return nil
}

func buildUserSegmentInsertString(userId int32, segments map[int32]string) (string, []interface{}) {
	var sbSql strings.Builder
	sbSql.WriteString("insert into user_segment (user_id, segment_id) values ")
	values := []interface{}{}
	var i int32
	i = 1
	for segmentId := range segments {
		values = append(values, userId, segmentId)
		sbSql.WriteString(fmt.Sprintf("($%d,$%d),", i, i+1))
		i += 2
	}
	sqlString := strings.Trim(sbSql.String(), ",")

	return sqlString, values
}

func buildUserSegmentLogInsertString(userId int32, segments map[int32]string, operation string) (string, []interface{}) {
	var sbSql strings.Builder
	sbSql.WriteString("insert into user_segment_log (user_id, segment_name, operation, operation_timestamp) values ")
	values := []interface{}{}
	var i int32
	i = 1
	for segmentId := range segments {
		values = append(values, userId, segmentId, operation, time.Now())
		sbSql.WriteString(fmt.Sprintf("($%d,$%d,$%d,$%d),", i, i+1, i+2, i+3))
		i += 4
	}
	sqlString := strings.Trim(sbSql.String(), ",")

	return sqlString, values
}
