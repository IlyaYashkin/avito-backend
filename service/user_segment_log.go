package service

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func addInfoToLog(userId int32, segments map[int32]string, operation string, tx *sql.Tx) error {
	sqlString, values := buildUserSegmentLogInsertString(userId, segments, operation)

	_, err := tx.Exec(sqlString, values...)
	if err != nil {
		return err
	}

	return nil
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
	sqlString := strings.TrimSuffix(sbSql.String(), ",")

	return sqlString, values
}
