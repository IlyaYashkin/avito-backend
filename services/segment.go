package services

import (
	"avito-backend/database"
	"avito-backend/models"
	"errors"
	"time"
)

func CreateSegment(requestData models.UpdateSegmentData) error {
	segmentName := requestData.Name

	db := database.Open()
	defer db.Close()

	var checkExists int
	err := db.QueryRow("select count(*) from segments where name = $1", segmentName).Scan(&checkExists)
	if err != nil {
		return errors.New("query error")
	}

	if checkExists != 0 {
		return errors.New("segment with this name already exists")
	}

	_, err = db.Exec("insert into segments (name, created_at) values ($1, $2)", segmentName, time.Now())
	if err != nil {
		return errors.New("insert error")
	}

	return nil
}
