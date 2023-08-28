package usersegmentlog

import (
	"fmt"
	"strings"
	"time"
)

func BuildUserSegmentLogInsertString(userId int32, segments map[int32]string, operation string) (string, []interface{}) {
	var sbSql strings.Builder
	sbSql.WriteString("insert into user_segment_log (user_id, segment_name, operation, operation_timestamp) values ")
	values := []interface{}{}
	var i int32
	i = 1
	for _, segment := range segments {
		values = append(values, userId, segment, operation, time.Now())
		sbSql.WriteString(fmt.Sprintf("($%d,$%d,$%d,$%d),", i, i+1, i+2, i+3))
		i += 4
	}
	sqlString := strings.TrimSuffix(sbSql.String(), ",")

	return sqlString, values
}
