package usersegment

import (
	"avito-backend/internal/database"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type UserSegment struct {
	Id        int32
	UserId    int32
	SegmentId int32
	Ttl       string
}

type UserSegmentRepository interface {
	GetByUserId(userId int32) ([]UserSegment, error)
	GetSegmentsNamesByUserId(userId int32) ([]string, error)
	BulkSaveForUser(userId int32, segments []int32) error
	BulkSaveForUserWithTtl(userId int32, segments []int32) error
	BulkDeleteForUser(userId int32, segmentsIds []int32) error
	BulkDeleteByExpiredTtl() ([]UserSegment, error)
}

type UserSegmentRepositoryDB struct {
	ex database.QueryExecutor
}

func NewUserSegmentRepository(ex database.QueryExecutor) *UserSegmentRepositoryDB {
	return &UserSegmentRepositoryDB{ex: ex}
}

func (repo UserSegmentRepositoryDB) GetByUserId(userId int32) ([]UserSegment, error) {
	rows, err := repo.ex.Query( /* sql */ `select segment_id, ttl from user_segment where user_id = $1`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userSegments []UserSegment

	for rows.Next() {
		var segment_id int32
		var ttl sql.NullString
		err := rows.Scan(&segment_id, &ttl)
		if err != nil {
			return nil, err
		}
		userSegments = append(userSegments, UserSegment{SegmentId: segment_id, Ttl: ttl.String})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userSegments, nil
}

func (repo UserSegmentRepositoryDB) GetSegmentsNamesByUserId(userId int32) ([]string, error) {
	rows, err := repo.ex.Query( /* sql */ `
		select segments.name from user_segment
		left join segments
		on user_segment.segment_id = segments.id
		where user_id = $1
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var segmentsNames []string

	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		segmentsNames = append(segmentsNames, name)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return segmentsNames, nil
}

func (repo UserSegmentRepositoryDB) BulkSaveForUser(userId int32, segments []int32) error {
	sqlString, values := BuildUserSegmentInsertString(userId, segments)
	_, err := repo.ex.Exec(sqlString, values...)
	if err != nil {
		return err
	}
	return nil
}

func (repo UserSegmentRepositoryDB) BulkSaveForUserWithTtl(userId int32, segments map[int32]string, ttls map[int32]time.Time) error {
	sqlString, values := BuildUserSegmentTtlInsertString(userId, segments, ttls)
	_, err := repo.ex.Exec(sqlString, values...)
	if err != nil {
		return err
	}
	return nil
}

func (repo UserSegmentRepositoryDB) BulkDeleteForUser(userId int32, segmentsIds []int32) error {
	sqlString := /* sql */ `delete from user_segment where user_id = $1 and segment_id = ANY($2)`
	_, err := repo.ex.Exec(sqlString, userId, pq.Array(segmentsIds))
	if err != nil {
		return err
	}
	return nil
}

func (repo UserSegmentRepositoryDB) BulkDeleteByExpiredTtl() ([]UserSegment, error) {
	sqlString := /* sql */ `
		delete from user_segment where ttl < now()
		returning user_id, segment_id
	`
	rows, err := repo.ex.Query(sqlString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deletedSegments := []UserSegment{}
	for rows.Next() {
		var user_id int32
		var segment_id int32
		err := rows.Scan(&user_id, &segment_id)
		if err != nil {
			return nil, err
		}
		deletedSegments = append(deletedSegments, UserSegment{UserId: user_id, SegmentId: segment_id})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return deletedSegments, nil
}
