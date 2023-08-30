package segment

import (
	"avito-backend/internal/database"

	"github.com/lib/pq"
)

func SelectSegmentsByName(segments []string, ex database.QueryExecutor) (map[int32]string, error) {
	rows, err := ex.Query("select id, name from segments where name = ANY($1)", pq.Array(segments))
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

func SelectSegmentsById(segmentsIds []int32, ex database.QueryExecutor) (map[int32]string, error) {
	rows, err := ex.Query("select id, name from segments where id = ANY($1)", pq.Array(segmentsIds))
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

func IsSegmentExists(name string, ex database.QueryExecutor) (bool, error) {
	var rowExists bool
	err := ex.QueryRow( /* sql */ `select exists(select true from segments where name=$1)`, name).Scan(&rowExists)
	if err != nil {
		return rowExists, err
	}
	return rowExists, nil
}

func InsertSegment(name string, ex database.QueryExecutor) (int32, error) {
	segmentId := 0
	err := ex.QueryRow( /* sql */ `insert into segments (name) values ($1) returning id`, name).Scan(&segmentId)
	if err != nil {
		return 0, err
	}
	return int32(segmentId), nil
}

func DelSegment(name string, ex database.QueryExecutor) error {
	_, err := ex.Exec( /* sql */ `delete from segments values where name = $1`, name)
	if err != nil {
		return err
	}
	return nil
}
