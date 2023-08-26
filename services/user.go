package services

import (
	"avito-backend/database"
	"avito-backend/dtos"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

func UpdateUserSegments(requestData dtos.UpdateUserSegments) error {
	db := database.Open()
	tx, err := db.Begin()
	if err != nil {
		return errors.New("begin transaction error")
	}
	defer tx.Rollback()

	err = createUser(requestData.UserId, tx)
	if err != nil {
		return err
	}

	err = addUserSegments(
		requestData.UserId,
		requestData.AddSegments,
		tx,
	)
	if err != nil {
		return err
	}

	err = deleteUserSegments(
		requestData.UserId,
		requestData.DeleteSegments,
		tx,
	)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func createUser(userId int32, tx *sql.Tx) error {
	_, err := tx.Exec("insert into users (id) values ($1) on conflict (id) do nothing", userId)
	if err != nil {
		return errors.New("insert user error")
	}
	return nil
}

func addUserSegments(userId int32, add []string, tx *sql.Tx) error {
	segmentsIds, err := getSegmentsIds(add, tx)
	if err != nil {
		return err
	}

	var sbSql strings.Builder
	sbSql.WriteString("insert into user_segment (user_id, segment_id) values ")
	values := []interface{}{}
	var i int32
	i = 1
	for idx, segmentId := range segmentsIds {
		values = append(values, userId, segmentId)
		if idx == len(segmentsIds)-1 {
			sbSql.WriteString(fmt.Sprintf("($%d,$%d)", i, i+1))
			break
		}
		sbSql.WriteString(fmt.Sprintf("($%d,$%d),", i, i+1))
		i += 2
	}
	sbSql.WriteString(" on conflict do nothing")
	sqlString := sbSql.String()

	_, err = tx.Exec(sqlString, values...)
	if err != nil {
		return err
	}

	return nil
}

func deleteUserSegments(userId int32, delete []string, tx *sql.Tx) error {
	return nil
}

func getSegmentsIds(segments []string, tx *sql.Tx) ([]int32, error) {
	rows, err := tx.Query("select id from segments where name = ANY($1)", pq.Array(segments))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var segmentsIds []int32

	for rows.Next() {
		var id int32
		err := rows.Scan(&id)
		if err != nil {
			return segmentsIds, err
		}
		segmentsIds = append(segmentsIds, id)
	}

	if err = rows.Err(); err != nil {
		return segmentsIds, err
	}

	return segmentsIds, nil
}
