package usersegmentlog

import (
	"fmt"
	"strings"
	"time"
)

func BuildUserSegmentLogInsertString(rows []UserSegmentLog, operation string) (string, []interface{}) {
	var sbSql strings.Builder
	sbSql.WriteString("insert into user_segment_log (user_id, segment_name, operation, operation_timestamp) values ")
	values := []interface{}{}
	var i int32
	i = 1
	for _, row := range rows {
		values = append(values, row.UserId, row.SegmentName, operation, time.Now())
		sbSql.WriteString(fmt.Sprintf("($%d,$%d,$%d,$%d),", i, i+1, i+2, i+3))
		i += 4
	}
	sqlString := strings.TrimSuffix(sbSql.String(), ",")
	return sqlString, values
}

func BuildUserSegmentLogSelectString(userId int32, date time.Time) (string, []interface{}) {
	var sbSql strings.Builder
	sbSql.WriteString( /* sql */ `
		select user_id, segment_name, operation, operation_timestamp
		from user_segment_log
	`)
	var values []interface{}
	var conditions []string
	if userId != 0 || !date.IsZero() {
		sbSql.WriteString(" where ")
	}
	if userId != 0 {
		conditions = append(conditions, "user_id = $1")
		values = append(values, userId)
	}
	if !date.IsZero() {
		conditions = append(conditions, "operation_timestamp <= $2")
		values = append(values, date)
	}
	sbSql.WriteString(strings.Join(conditions, " and "))
	return sbSql.String(), values
}
