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

type UpdatedUserSegments struct {
	AddedSegments           []string
	AddedTtlSegments        []string
	AddedPercentageSegments []string
	DeletedSegments         []string
}

type RequestSegmentWithTtl struct {
	Segment string
	Ttl     string
}

const TTL_TIME_FORMAT = time.RFC3339

func UpdateUserSegmentsService(requestData RequestUpdateUserSegments) (UpdatedUserSegments, error) {
	db := database.Get()
	tx, err := db.Begin()
	if err != nil {
		return UpdatedUserSegments{}, err
	}
	defer tx.Rollback()

	userRepository := user.NewUserRepository(tx)

	isUsrInserted, err := userRepository.Save(requestData.UserId)
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

func GetUserSegmentsService(userId int32) ([]string, error) {
	db := database.Get()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	userSegmentRepo := NewUserSegmentRepository(tx)

	err = ClearExpiredTtlUserSegments(tx)
	if err != nil {
		return nil, err
	}

	userSegmentsNames, err := userSegmentRepo.GetSegmentsNamesByUserId(userId)
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

func ClearExpiredTtlUserSegments(ex database.QueryExecutor) error {
	segmentRepo := segment.NewSegmentRepository(ex)
	userSegmentRepo := NewUserSegmentRepository(ex)
	userSegmentLogRepo := usersegmentlog.NewUserSegmentLogRepository(ex)

	deletedSegments, err := userSegmentRepo.BulkDeleteByExpiredTtl()
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

	segments, err := segmentRepo.GetById(segmentsIds)
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
	err = userSegmentLogRepo.Save(logRows, usersegmentlog.LOG_OPERATION_DELETE_TTL)
	if err != nil {
		return err
	}

	return nil
}

func addPercentageSegments(userId int32, tx *sql.Tx) ([]string, error) {
	segmentPercentageRepo := segmentpercentage.NewSegmentPercentageRepository(tx)
	segmentRepo := segment.NewSegmentRepository(tx)
	userSegmentRepo := NewUserSegmentRepository(tx)
	userSegmentLogRepo := usersegmentlog.NewUserSegmentLogRepository(tx)

	err := segmentPercentageRepo.IncrementCounters()
	if err != nil {
		return nil, err
	}
	segmentsIds, err := segmentPercentageRepo.PickSegments()
	if err != nil {
		return nil, err
	}
	if len(segmentsIds) == 0 {
		return nil, nil
	}

	segments, err := segmentRepo.GetById(segmentsIds)
	err = userSegmentRepo.BulkSaveForUser(userId, segmentsIds)
	if err != nil {
		return nil, err
	}

	var logRows []usersegmentlog.UserSegmentLog
	for _, segment := range segments {
		logRows = append(logRows, usersegmentlog.UserSegmentLog{UserId: userId, SegmentName: segment})
	}
	userSegmentLogRepo.Save(logRows, usersegmentlog.LOG_OPERATION_ADD_PERCENTAGE)

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

	segmentRepo := segment.NewSegmentRepository(tx)
	userSegmentRepo := NewUserSegmentRepository(tx)
	userSegmentLogRepo := usersegmentlog.NewUserSegmentLogRepository(tx)

	segments, err := segmentRepo.GetByName(segmentsToAdd)
	if err != nil {
		return nil, err
	}
	userSegments, err := userSegmentRepo.GetByUserId(userId)
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
	err = userSegmentRepo.BulkSaveForUser(userId, segmentsIds)
	if err != nil {
		return nil, err
	}

	var logRows []usersegmentlog.UserSegmentLog
	for _, segment := range segments {
		logRows = append(logRows, usersegmentlog.UserSegmentLog{UserId: userId, SegmentName: segment})
	}
	err = userSegmentLogRepo.Save(logRows, usersegmentlog.LOG_OPERATION_ADD)
	if err != nil {
		return nil, err
	}

	var rsSegments []string
	for _, segment := range segments {
		rsSegments = append(rsSegments, segment)
	}
	return rsSegments, nil
}

func addTtlSegments(userId int32, segmentsToAdd []RequestSegmentWithTtl, tx *sql.Tx) ([]string, error) {
	if len(segmentsToAdd) == 0 {
		return nil, nil
	}

	segmentRepo := segment.NewSegmentRepository(tx)
	userSegmentRepo := NewUserSegmentRepository(tx)
	userSegmentLogRepo := usersegmentlog.NewUserSegmentLogRepository(tx)

	var segmentsWithTtlArr []string
	for _, segmentTtl := range segmentsToAdd {
		segmentsWithTtlArr = append(segmentsWithTtlArr, segmentTtl.Segment)
	}
	segments, err := segmentRepo.GetByName(segmentsWithTtlArr)
	if err != nil {
		return nil, err
	}
	userSegments, err := userSegmentRepo.GetByUserId(userId)
	if err != nil {
		return nil, err
	}
	reversedSegments := flipSegmentsMap(segments)
	ttls := make(map[int32]string)
	for _, segmentTtl := range segmentsToAdd {
		if id, exists := reversedSegments[segmentTtl.Segment]; exists {
			ttls[id] = segmentTtl.Ttl
		}
	}
	segments, timeTtls := sanitizeTtls(segments, ttls, userSegments)
	if len(segments) == 0 {
		return nil, nil
	}

	err = userSegmentRepo.BulkSaveForUserWithTtl(userId, segments, timeTtls)
	if err != nil {
		return nil, err
	}

	var logRows []usersegmentlog.UserSegmentLog
	for _, segment := range segments {
		logRows = append(logRows, usersegmentlog.UserSegmentLog{UserId: userId, SegmentName: segment})
	}
	err = userSegmentLogRepo.Save(logRows, usersegmentlog.LOG_OPERATION_ADD)
	if err != nil {
		return nil, err
	}

	var rsSegments []string
	for _, segment := range segments {
		rsSegments = append(rsSegments, segment)
	}
	return rsSegments, nil
}

func deleteSegments(userId int32, segmentsToDelete []string, tx *sql.Tx) ([]string, error) {
	if len(segmentsToDelete) == 0 {
		return nil, nil
	}

	segmentRepo := segment.NewSegmentRepository(tx)
	userSegmentRepo := NewUserSegmentRepository(tx)
	userSegmentLogRepo := usersegmentlog.NewUserSegmentLogRepository(tx)

	segments, err := segmentRepo.GetByName(segmentsToDelete)
	if err != nil {
		return nil, err
	}
	userSegments, err := userSegmentRepo.GetByUserId(userId)
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
	err = userSegmentRepo.BulkDeleteForUser(userId, segmentsIdsToDelete)
	if err != nil {
		return nil, err
	}

	var logRows []usersegmentlog.UserSegmentLog
	for _, segment := range segments {
		logRows = append(logRows, usersegmentlog.UserSegmentLog{UserId: userId, SegmentName: segment})
	}
	err = userSegmentLogRepo.Save(logRows, usersegmentlog.LOG_OPERATION_DELETE)
	if err != nil {
		return nil, err
	}

	var rsSegments []string
	for _, segment := range segments {
		rsSegments = append(rsSegments, segment)
	}
	return rsSegments, nil
}

func splitSegments(segments []any) ([]string, []RequestSegmentWithTtl) {
	var addSegments []string
	var addSegmentsWithTtl []RequestSegmentWithTtl

	for _, value := range segments {
		segmentTtl, ok := value.(map[string]interface{})
		if ok {
			segment, segmentOk := segmentTtl["segment"].(string)
			ttl, ttlOk := segmentTtl["ttl"].(string)

			if segmentOk && ttlOk {
				segmentTtlStruct := RequestSegmentWithTtl{
					Segment: segment,
					Ttl:     ttl,
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
