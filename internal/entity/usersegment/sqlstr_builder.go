package usersegment

import (
	"fmt"
	"strings"
	"time"
)

func buildUserSegmentInsertString(userId int32, segments []int32) (string, []interface{}) {
	var sbSql strings.Builder
	sbSql.WriteString( /* sql */ `insert into user_segment (user_id, segment_id) values `)
	values := []interface{}{}
	var i int32
	i = 1
	for _, segmentId := range segments {
		values = append(values, userId, segmentId)
		sbSql.WriteString(fmt.Sprintf("($%d,$%d),", i, i+1))
		i += 2
	}
	sqlString := strings.TrimSuffix(sbSql.String(), ",")

	return sqlString, values
}

func buildUserSegmentTtlInsertString(userId int32, segments map[int32]string, ttls map[int32]time.Time) (string, []interface{}) {
	var sbSql strings.Builder
	sbSql.WriteString( /* sql */ `insert into user_segment (user_id, segment_id, ttl) values `)
	values := []interface{}{}
	var i int32
	i = 1
	for segmentId := range segments {
		values = append(values, userId, segmentId, ttls[segmentId])
		sbSql.WriteString(fmt.Sprintf("($%d,$%d,$%d),", i, i+1, i+2))
		i += 3
	}
	buf := strings.TrimSuffix(sbSql.String(), ",")
	sbSql.Reset()
	sbSql.WriteString(buf)
	sbSql.WriteString(
		`on conflict (user_id, segment_id) do update
			set
			user_id=EXCLUDED.user_id,
			segment_id=EXCLUDED.segment_id,
			ttl=EXCLUDED.ttl`,
	)

	return sbSql.String(), values
}
