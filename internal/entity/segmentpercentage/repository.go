package segmentpercentage

import "avito-backend/internal/database"

type SegmentPercentageRepository interface {
	Save(segmentId int32, userPercentage float32) error
	IncrementCounters() error
	PickSegments() (int32, error)
}

type SegmentPercentageRepositoryDB struct {
	ex database.QueryExecutor
}

func NewSegmentPercentageRepository(ex database.QueryExecutor) *SegmentPercentageRepositoryDB {
	return &SegmentPercentageRepositoryDB{ex: ex}
}

func (repo *SegmentPercentageRepositoryDB) Save(segmentId int32, userPercentage float32) error {
	_, err := repo.ex.Exec("insert into segment_percentage (segment_id, user_percentage) values ($1, $2)", segmentId, userPercentage)
	if err != nil {
		return err
	}
	return nil
}

func (repo *SegmentPercentageRepositoryDB) IncrementCounters() error {
	_, err := repo.ex.Exec( /* sql */ `
		update segment_percentage
		set user_counter = user_counter + 1
		where user_counter - 100 / user_percentage < 0
	`)
	if err != nil {
		return err
	}
	return nil
}

func (repo *SegmentPercentageRepositoryDB) PickSegments() ([]int32, error) {
	rows, err := repo.ex.Query( /* sql */ `
		update segment_percentage
		set user_counter = user_counter - 100 / user_percentage
		where user_counter - 100 / user_percentage >= 0
		returning segment_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var segments []int32
	for rows.Next() {
		var id int32
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		segments = append(segments, id)
	}

	return segments, nil
}
