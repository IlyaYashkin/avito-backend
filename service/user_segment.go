package service

import (
	"avito-backend/database"
	"avito-backend/database/dbaccess"
	"avito-backend/dto"
	"avito-backend/model"
	"database/sql"
	"time"
)

type UpdatedUserSegments struct {
	AddedSegments           []string
	AddedSegmentsTtl        []string
	AddedSegmentsPercentage []string
	DeletedSegments         []string
}

type requestSegmentWithTtl struct {
	segment string
	ttl     string
}

const ttlTimeFormat = time.RFC3339

func UpdateUserSegments(requestData dto.UpdateUserSegments) (UpdatedUserSegments, error) {
	db := database.Open()
	tx, err := db.Begin()
	if err != nil {
		return UpdatedUserSegments{}, err
	}
	defer tx.Rollback()

	isUsrInserted, err := dbaccess.InsertUser(requestData.UserId, tx)
	if err != nil {
		return UpdatedUserSegments{}, err
	}

	if isUsrInserted {
		_, err = addPercentageSegments(requestData.UserId, tx)
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
		AddedSegments:    addedSegments,
		AddedSegmentsTtl: addedTtlSegments,
		DeletedSegments:  deletedSegments,
	}

	return updatedUserSegments, nil
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
	err := dbaccess.IncrementCounters(tx)
	if err != nil {
		return nil, err
	}
	segments, err := dbaccess.PickSegments(tx)
	if err != nil {
		return nil, err
	}

	if len(segments) == 0 {
		return nil, nil
	}

	err = dbaccess.InsertUsrSegments(userId, segments, tx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func addSegments(userId int32, segmentsToAdd []string, tx *sql.Tx) ([]string, error) {
	if len(segmentsToAdd) == 0 {
		return nil, nil
	}

	segments, err := dbaccess.GetMatchedSegments(segmentsToAdd, tx)
	if err != nil {
		return nil, err
	}
	userSegments, err := dbaccess.GetUsrSegments(userId, tx)
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

	err = dbaccess.InsertUsrSegments(userId, segmentsIds, tx)
	if err != nil {
		return nil, err
	}

	err = dbaccess.InsertLog(userId, segments, "addition", tx)
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
	segments, err := dbaccess.GetMatchedSegments(segmentsWithTtlArr, tx)
	if err != nil {
		return nil, err
	}
	userSegments, err := dbaccess.GetUsrSegments(userId, tx)
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

	err = dbaccess.InsertUsrSegmentsTtl(userId, segments, timeTtls, tx)
	if err != nil {
		return nil, err
	}

	err = dbaccess.InsertLog(userId, segments, "addition", tx)
	if err != nil {
		return nil, err
	}

	var rsSegments []string
	for _, segment := range segments {
		rsSegments = append(rsSegments, segment)
	}
	return rsSegments, nil
}

func sanitizeTtls(segments map[int32]string, ttls map[int32]string, userSegments []model.UserSegment) (map[int32]string, map[int32]time.Time) {
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

func deleteSegments(userId int32, segmentsToDelete []string, tx *sql.Tx) ([]string, error) {
	if len(segmentsToDelete) == 0 {
		return nil, nil
	}

	segments, err := dbaccess.GetMatchedSegments(segmentsToDelete, tx)
	if err != nil {
		return nil, err
	}
	userSegments, err := dbaccess.GetUsrSegments(userId, tx)
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
	err = dbaccess.DeleteUsrSegments(userId, segmentsIdsToDelete, tx)
	if err != nil {
		return nil, err
	}

	err = dbaccess.InsertLog(userId, segments, "deletion", tx)
	if err != nil {
		return nil, err
	}

	var rsSegments []string

	for _, segment := range segments {
		rsSegments = append(rsSegments, segment)
	}

	return rsSegments, nil
}

func getUserSegmentsIds(userSegments []model.UserSegment) []int32 {
	var userSegmentsIds []int32
	for _, userSegment := range userSegments {
		userSegmentsIds = append(userSegmentsIds, userSegment.SegmentId)
	}
	return userSegmentsIds
}

func getUserSegmentsTtlsMap(userSegments []model.UserSegment) map[int32]time.Time {
	userSegmentsTtlsMap := make(map[int32]time.Time, len(userSegments))
	for _, userSegment := range userSegments {
		ttlTime, _ := time.Parse(ttlTimeFormat, userSegment.Ttl)
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
