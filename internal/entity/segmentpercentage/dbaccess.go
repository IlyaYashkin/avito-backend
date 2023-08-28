package segmentpercentage

import "avito-backend/internal/database"

func InsertSegmentPercentage(segmentId int32, userPercentage float32, ex database.QueryExecutor) error {
	_, err := ex.Exec("insert into segment_percentage (segment_id, user_percentage) values ($1, $2)", segmentId, userPercentage)
	if err != nil {
		return err
	}
	return nil
}

func IncrementCounters(ex database.QueryExecutor) error {
	_, err := ex.Exec( /* sql */ `
		update segment_percentage
		set user_counter = user_counter + 1
		where user_counter - 100 / user_percentage < 0
	`)
	if err != nil {
		return err
	}
	return nil
}

func PickSegments(ex database.QueryExecutor) ([]int32, error) {
	rows, err := ex.Query( /* sql */ `
		update segment_percentage
		set user_counter = user_counter - 100 / user_percentage
		where user_counter - 100 / user_percentage >= 0
		returning segment_id
	`)
	if err != nil {
		return nil, err
	}

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
