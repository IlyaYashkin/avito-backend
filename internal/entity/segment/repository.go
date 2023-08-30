package segment

import (
	"avito-backend/internal/database"

	"github.com/lib/pq"
)

type SegmentRepository interface {
	GetByName(segments []string) (map[int32]string, error)
	GetById(ids []int32) (map[int32]string, error)
	IsExistsByName(name string) (bool, error)
	Save(name string) (int32, error)
	DeleteByName(name string) error
}

type SegmentRepositoryDB struct {
	ex database.QueryExecutor
}

func NewSegmentRepository(ex database.QueryExecutor) *SegmentRepositoryDB {
	return &SegmentRepositoryDB{ex: ex}
}

func (repo *SegmentRepositoryDB) GetByName(segments []string) (map[int32]string, error) {
	rows, err := repo.ex.Query( /* sql */ `select id, name from segments where name = ANY($1)`, pq.Array(segments))
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

func (repo *SegmentRepositoryDB) GetById(ids []int32) (map[int32]string, error) {
	rows, err := repo.ex.Query( /* sql */ `select id, name from segments where id = ANY($1)`, pq.Array(ids))
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

func (repo *SegmentRepositoryDB) IsExistsByName(name string) (bool, error) {
	var rowExists bool
	err := repo.ex.QueryRow( /* sql */ `select exists(select true from segments where name=$1)`, name).Scan(&rowExists)
	if err != nil {
		return rowExists, err
	}
	return rowExists, nil
}

func (repo *SegmentRepositoryDB) Save(name string) (int32, error) {
	segmentId := 0
	err := repo.ex.QueryRow( /* sql */ `insert into segments (name) values ($1) returning id`, name).Scan(&segmentId)
	if err != nil {
		return 0, err
	}
	return int32(segmentId), nil
}

func (repo *SegmentRepositoryDB) DeleteByName(name string) error {
	_, err := repo.ex.Exec( /* sql */ `delete from segments values where name = $1`, name)
	if err != nil {
		return err
	}
	return nil
}
