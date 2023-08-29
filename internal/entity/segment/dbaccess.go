package segment

import "avito-backend/internal/database"

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
