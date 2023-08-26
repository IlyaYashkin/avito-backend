package services

import (
	"avito-backend/database"
	"avito-backend/dtos"
	"errors"
)

func CreateSegment(requestData dtos.UpdateSegment) error {
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

	_, err = db.Exec("insert into segments (name) values ($1)", segmentName)
	if err != nil {
		return errors.New("insert error")
	}

	return nil
}

func DeleteSegment(requestData dtos.UpdateSegment) error {
	segmentName := requestData.Name

	db := database.Open()
	defer db.Close()

	var checkExists int
	err := db.QueryRow("select count(*) from segments where name = $1", segmentName).Scan(&checkExists)
	if err != nil {
		return errors.New("query error")
	}
	if checkExists == 0 {
		return errors.New("segment with this name not exists")
	}

	_, err = db.Exec("delete from segments values where name = $1", segmentName)
	if err != nil {
		return errors.New("delete error")
	}

	return nil
}
