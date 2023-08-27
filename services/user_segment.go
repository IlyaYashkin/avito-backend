package services

import (
	"avito-backend/database"
	"avito-backend/dtos"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
)

type requestSegmentWithTtl struct {
	segment string
	ttl     string
}

type segment struct {
	id   int32
	name string
}

type userSegment struct {
	segment_id int32
	ttl        string
}

type UpdatedUserSegments struct {
	AddedSegments        []string
	AddedSegmentsWithTtl []string
	DeletedSegments      []string
}

const ttlTimeFormat = time.RFC3339

func UpdateUserSegments(requestData dtos.UpdateUserSegments) (UpdatedUserSegments, error) {
	db := database.Open()
	tx, err := db.Begin()
	if err != nil {
		return UpdatedUserSegments{}, err
	}
	defer tx.Rollback()

	err = createUserUsingTransaction(requestData.UserId, tx)
	if err != nil {
		return UpdatedUserSegments{}, err
	}

	addSegments, addSegmentsWithTtl := splitAddSegments(requestData.AddSegments)

	addedSegments, err := addUserSegments(
		requestData.UserId,
		addSegments,
		tx,
	)
	if err != nil {
		return UpdatedUserSegments{}, err
	}

	addedSegmentsWithTtl, err := addUserSegmentsWithTtl(
		requestData.UserId,
		addSegmentsWithTtl,
		tx,
	)
	if err != nil {
		return UpdatedUserSegments{}, err
	}

	deletedSegments, err := deleteUserSegments(
		requestData.UserId,
		requestData.DeleteSegments,
		tx,
	)
	if err != nil {
		return UpdatedUserSegments{}, err
	}

	tx.Commit()

	updatedUserSegments := UpdatedUserSegments{
		AddedSegments:        addedSegments,
		AddedSegmentsWithTtl: addedSegmentsWithTtl,
		DeletedSegments:      deletedSegments,
	}

	return updatedUserSegments, nil
}

func splitAddSegments(segments []any) ([]string, []requestSegmentWithTtl) {
	var addSegments []string
	var addSegmentsWithTtl []requestSegmentWithTtl

	for _, value := range segments {
		segmentTtl, ok := value.(map[string]interface{})
		if ok {
			segment, segmentOk := segmentTtl["segment"].(string)
			ttl, ttlOk := segmentTtl["ttl"].(string)

			if segmentOk && ttlOk {
				segmentTtlStruct := requestSegmentWithTtl{
					segment: segment,
					ttl:     ttl,
				}
				addSegmentsWithTtl = append(addSegmentsWithTtl, segmentTtlStruct)
				continue
			}
		}
		segment, ok := value.(string)
		if ok {
			addSegments = append(addSegments, segment)
		}
	}

	return addSegments, addSegmentsWithTtl
}

func addUserSegments(userId int32, segmentsToAdd []string, tx *sql.Tx) ([]string, error) {
	segments, err := getSegments(segmentsToAdd, tx)
	if err != nil {
		return nil, err
	}
	userSegments, err := getUserSegments(userId, tx)
	if err != nil {
		return nil, err
	}
	userSegmentsIds := getUserSegmentsIds(userSegments)

	for _, segmentId := range userSegmentsIds {
		if _, exists := segments[segmentId]; exists {
			delete(segments, segmentId)
		}
	}
	if len(segments) == 0 {
		return nil, nil
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

	var rsSegments []string

	for _, segment := range segments {
		rsSegments = append(rsSegments, segment)
	}

	return rsSegments, nil
}

func addUserSegmentsWithTtl(userId int32, segmentsToAdd []requestSegmentWithTtl, tx *sql.Tx) ([]string, error) {
	var segmentsWithTtlArr []string
	for _, segmentTtl := range segmentsToAdd {
		segmentsWithTtlArr = append(segmentsWithTtlArr, segmentTtl.segment)
	}
	segments, err := getSegments(segmentsWithTtlArr, tx)
	if err != nil {
		return nil, err
	}
	userSegments, err := getUserSegments(userId, tx)
	if err != nil {
		return nil, err
	}
	reversedSegments := flipSegmentsMap(segments)
	ttls := make(map[int32]string)
	for _, segmentTtl := range segmentsToAdd {
		if id, exists := reversedSegments[segmentTtl.segment]; exists {
			ttls[id] = segmentTtl.ttl
		}
	}

	segments, timeTtls := sanitizeSegmentsTtls(segments, ttls, userSegments)

	if len(segments) == 0 {
		return nil, nil
	}

	sqlString, values := buildUserSegmentWithTtlInsertString(userId, segments, timeTtls)
	_, err = tx.Exec(sqlString, values...)
	if err != nil {
		return nil, err
	}

	err = addInfoToLog(userId, segments, "addition", tx)
	if err != nil {
		return nil, err
	}

	var rsSegments []string
	for _, segment := range segments {
		rsSegments = append(rsSegments, segment)
	}
	return rsSegments, nil
}

func sanitizeSegmentsTtls(segments map[int32]string, ttls map[int32]string, userSegments []userSegment) (map[int32]string, map[int32]time.Time) {
	timeTtls := make(map[int32]time.Time)
	userTimeTtls := getUserSegmentsTtlsMap(userSegments)

	for id, ttl := range ttls {
		time, err := time.Parse(ttlTimeFormat, ttl)
		if err != nil {
			delete(segments, id)
			continue
		}
		if time.Equal(userTimeTtls[id]) {
			delete(segments, id)
			continue
		}
		timeTtls[id] = time
	}

	return segments, timeTtls
}

func deleteUserSegments(userId int32, segmentsToDelete []string, tx *sql.Tx) ([]string, error) {
	segments, err := getSegments(segmentsToDelete, tx)
	if err != nil {
		return nil, err
	}
	userSegments, err := getUserSegments(userId, tx)
	if err != nil {
		return nil, err
	}
	userSegmentsIds := getUserSegmentsIds(userSegments)

	for _, segmentId := range userSegmentsIds {
		if _, exists := segments[segmentId]; !exists {
			delete(segments, segmentId)
		}
	}
	if len(segments) == 0 {
		return nil, nil
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

	var rsSegments []string

	for _, segment := range segments {
		rsSegments = append(rsSegments, segment)
	}

	return rsSegments, nil
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

func getUserSegments(userId int32, tx *sql.Tx) ([]userSegment, error) {
	rows, err := tx.Query("select segment_id, ttl from user_segment where user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userSegments []userSegment

	for rows.Next() {
		var segment_id int32
		var ttl sql.NullString
		err := rows.Scan(&segment_id, &ttl)
		if err != nil {
			return nil, err
		}
		userSegments = append(userSegments, userSegment{segment_id, ttl.String})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userSegments, nil
}

func getUserSegmentsIds(userSegments []userSegment) []int32 {
	var userSegmentsIds []int32
	for _, userSegment := range userSegments {
		userSegmentsIds = append(userSegmentsIds, userSegment.segment_id)
	}
	return userSegmentsIds
}

func getUserSegmentsTtlsMap(userSegments []userSegment) map[int32]time.Time {
	userSegmentsTtlsMap := make(map[int32]time.Time, len(userSegments))
	for _, userSegment := range userSegments {
		ttlTime, _ := time.Parse(ttlTimeFormat, userSegment.ttl)
		userSegmentsTtlsMap[userSegment.segment_id] = ttlTime
	}
	return userSegmentsTtlsMap
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
	sqlString := strings.TrimSuffix(sbSql.String(), ",")

	return sqlString, values
}

func buildUserSegmentWithTtlInsertString(userId int32, segments map[int32]string, ttls map[int32]time.Time) (string, []interface{}) {
	var sbSql strings.Builder
	sbSql.WriteString("insert into user_segment (user_id, segment_id, ttl) values ")
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

func flipSegmentsMap(segments map[int32]string) map[string]int32 {
	reversedSegments := make(map[string]int32, len(segments))
	for idx, segment := range segments {
		reversedSegments[segment] = idx
	}
	return reversedSegments
}
