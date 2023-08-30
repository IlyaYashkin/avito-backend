package usersegment

import (
	"avito-backend/internal/database"
	"avito-backend/internal/entity/segment"
	"avito-backend/internal/entity/segmentpercentage"
	"avito-backend/internal/entity/user"
	"avito-backend/internal/entity/usersegmentlog"
	"database/sql"
	"time"
)

type UserSegment struct {
	Id        int32
	UserId    int32
	SegmentId int32
	Ttl       string
}

type UpdatedUserSegments struct {
	AddedSegments           []string
	AddedTtlSegments        []string
	AddedPercentageSegments []string
	DeletedSegments         []string
}

type requestSegmentWithTtl struct {
	segment string
	ttl     string
}

const TTL_TIME_FORMAT = time.RFC3339

func updateUserSegments(requestData RequestUpdateUserSegments) (UpdatedUserSegments, error) {
	db := database.Get()
	tx, err := db.Begin()
	if err != nil {
		return UpdatedUserSegments{}, err
	}
	defer tx.Rollback()

	isUsrInserted, err := user.InsertUser(requestData.UserId, tx)
	if err != nil {
		return UpdatedUserSegments{}, err
	}

	var addedPercentageSegments []string
	if isUsrInserted {
		addedPercentageSegments, err = addPercentageSegments(requestData.UserId, tx)
		if err != nil {
			return UpdatedUserSegments{}, err
		}
	}

	segments, segmentsTtl := splitSegments(requestData.AddSegments)
	addedSegments, err := addSegments(
		requestData.UserId,
		segments,
		tx,
	)
	if err != nil {
		return UpdatedUserSegments{}, err
	}
	addedTtlSegments, err := addTtlSegments(
		requestData.UserId,
		segmentsTtl,
		tx,
	)
	if err != nil {
		return UpdatedUserSegments{}, err
	}

	deletedSegments, err := deleteSegments(
		requestData.UserId,
		requestData.DeleteSegments,
		tx,
	)
	if err != nil {
		return UpdatedUserSegments{}, err
	}

	tx.Commit()

	updatedUserSegments := UpdatedUserSegments{
		AddedSegments:           addedSegments,
		AddedTtlSegments:        addedTtlSegments,
		AddedPercentageSegments: addedPercentageSegments,
		DeletedSegments:         deletedSegments,
	}

	return updatedUserSegments, nil
}

func getUserSegments(userId int32) ([]string, error) {
	db := database.Get()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = clearExpiredTtlUserSegments(tx)
	if err != nil {
		return nil, err
	}

	userSegmentsNames, err := SelectUserSegmentNamesByUserId(userId, tx)
	if err != nil {
		return nil, err
	}

	tx.Commit()

	var segmentsRs []string
	for _, userSegmentName := range userSegmentsNames {
		segmentsRs = append(segmentsRs, userSegmentName)
	}

	return segmentsRs, nil
}

func clearExpiredTtlUserSegments(ex database.QueryExecutor) error {
	deletedSegments, err := deleteExpiredTtlUserSegments(ex)
	if err != nil {
		return err
	}

	if len(deletedSegments) == 0 {
		return nil
	}

	var segmentsIds []int32
	for _, deletedSegment := range deletedSegments {
		segmentsIds = append(segmentsIds, deletedSegment.SegmentId)
	}

	segments, err := segment.SelectSegmentsById(segmentsIds, ex)
	if err != nil {
		return err
	}

	var logRows []usersegmentlog.UserSegmentLog
	for _, deletedSegment := range deletedSegments {
		logRows = append(
			logRows,
			usersegmentlog.UserSegmentLog{
				UserId:      deletedSegment.UserId,
				SegmentName: segments[deletedSegment.SegmentId],
			},
		)
	}

	err = usersegmentlog.InsertLog(logRows, usersegmentlog.LOG_OPERATION_DELETE_TTL, ex)
	if err != nil {
		return err
	}

	return nil
}

func splitSegments(segments []any) ([]string, []requestSegmentWithTtl) {
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

func addPercentageSegments(userId int32, tx *sql.Tx) ([]string, error) {
	err := segmentpercentage.IncrementCounters(tx)
	if err != nil {
		return nil, err
	}
	segmentsIds, err := segmentpercentage.PickSegments(tx)
	if err != nil {
		return nil, err
	}

	if len(segmentsIds) == 0 {
		return nil, nil
	}

	segments, err := segment.SelectSegmentsById(segmentsIds, tx)

	err = InsertUserSegment(userId, segmentsIds, tx)
	if err != nil {
		return nil, err
	}

	var logRows []usersegmentlog.UserSegmentLog
	for _, segment := range segments {
		logRows = append(logRows, usersegmentlog.UserSegmentLog{UserId: userId, SegmentName: segment})
	}

	usersegmentlog.InsertLog(logRows, usersegmentlog.LOG_OPERATION_ADD_PERCENTAGE, tx)

	var rsSegments []string
	for _, segment := range segments {
		rsSegments = append(rsSegments, segment)
	}
	return rsSegments, nil
}

func addSegments(userId int32, segmentsToAdd []string, tx *sql.Tx) ([]string, error) {
	if len(segmentsToAdd) == 0 {
		return nil, nil
	}

	segments, err := segment.SelectSegmentsByName(segmentsToAdd, tx)
	if err != nil {
		return nil, err
	}
	userSegments, err := SelectUserSegmentByUserId(userId, tx)
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

	var segmentsIds []int32
	for id := range segments {
		segmentsIds = append(segmentsIds, id)
	}

	err = InsertUserSegment(userId, segmentsIds, tx)
	if err != nil {
		return nil, err
	}

	var logRows []usersegmentlog.UserSegmentLog
	for _, segment := range segments {
		logRows = append(logRows, usersegmentlog.UserSegmentLog{UserId: userId, SegmentName: segment})
	}

	err = usersegmentlog.InsertLog(logRows, usersegmentlog.LOG_OPERATION_ADD, tx)
	if err != nil {
		return nil, err
	}

	var rsSegments []string

	for _, segment := range segments {
		rsSegments = append(rsSegments, segment)
	}

	return rsSegments, nil
}

func addTtlSegments(userId int32, segmentsToAdd []requestSegmentWithTtl, tx *sql.Tx) ([]string, error) {
	if len(segmentsToAdd) == 0 {
		return nil, nil
	}

	var segmentsWithTtlArr []string
	for _, segmentTtl := range segmentsToAdd {
		segmentsWithTtlArr = append(segmentsWithTtlArr, segmentTtl.segment)
	}
	segments, err := segment.SelectSegmentsByName(segmentsWithTtlArr, tx)
	if err != nil {
		return nil, err
	}
	userSegments, err := SelectUserSegmentByUserId(userId, tx)
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

	segments, timeTtls := sanitizeTtls(segments, ttls, userSegments)

	if len(segments) == 0 {
		return nil, nil
	}

	err = InsertUserTtlSegment(userId, segments, timeTtls, tx)
	if err != nil {
		return nil, err
	}

	var logRows []usersegmentlog.UserSegmentLog
	for _, segment := range segments {
		logRows = append(logRows, usersegmentlog.UserSegmentLog{UserId: userId, SegmentName: segment})
	}

	err = usersegmentlog.InsertLog(logRows, usersegmentlog.LOG_OPERATION_ADD, tx)
	if err != nil {
		return nil, err
	}

	var rsSegments []string
	for _, segment := range segments {
		rsSegments = append(rsSegments, segment)
	}
	return rsSegments, nil
}

func sanitizeTtls(segments map[int32]string, ttls map[int32]string, userSegments []UserSegment) (map[int32]string, map[int32]time.Time) {
	timeTtls := make(map[int32]time.Time)
	userTimeTtls := getUserSegmentsTtlsMap(userSegments)

	for id, ttl := range ttls {
		time, err := time.Parse(TTL_TIME_FORMAT, ttl)
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

func deleteSegments(userId int32, segmentsToDelete []string, tx *sql.Tx) ([]string, error) {
	if len(segmentsToDelete) == 0 {
		return nil, nil
	}

	segments, err := segment.SelectSegmentsByName(segmentsToDelete, tx)
	if err != nil {
		return nil, err
	}
	userSegments, err := SelectUserSegmentByUserId(userId, tx)
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

	var segmentsIdsToDelete []int32
	for i := range segments {
		segmentsIdsToDelete = append(segmentsIdsToDelete, i)
	}
	err = DeleteUserSegment(userId, segmentsIdsToDelete, tx)
	if err != nil {
		return nil, err
	}

	var logRows []usersegmentlog.UserSegmentLog
	for _, segment := range segments {
		logRows = append(logRows, usersegmentlog.UserSegmentLog{UserId: userId, SegmentName: segment})
	}

	err = usersegmentlog.InsertLog(logRows, usersegmentlog.LOG_OPERATION_DELETE, tx)
	if err != nil {
		return nil, err
	}

	var rsSegments []string

	for _, segment := range segments {
		rsSegments = append(rsSegments, segment)
	}

	return rsSegments, nil
}

func getUserSegmentsIds(userSegments []UserSegment) []int32 {
	var userSegmentsIds []int32
	for _, userSegment := range userSegments {
		userSegmentsIds = append(userSegmentsIds, userSegment.SegmentId)
	}
	return userSegmentsIds
}

func getUserSegmentsTtlsMap(userSegments []UserSegment) map[int32]time.Time {
	userSegmentsTtlsMap := make(map[int32]time.Time, len(userSegments))
	for _, userSegment := range userSegments {
		ttlTime, _ := time.Parse(TTL_TIME_FORMAT, userSegment.Ttl)
		userSegmentsTtlsMap[userSegment.SegmentId] = ttlTime
	}
	return userSegmentsTtlsMap
}

func flipSegmentsMap(segments map[int32]string) map[string]int32 {
	reversedSegments := make(map[string]int32, len(segments))
	for idx, segment := range segments {
		reversedSegments[segment] = idx
	}
	return reversedSegments
}
