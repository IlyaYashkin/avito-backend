package services

import (
	"avito-backend/database"
	"avito-backend/dtos"
	"database/sql"
	"errors"
)

func CreateSegment(requestData dtos.UpdateSegment) error {
	db := database.Get()

	rowExists, err := checkSegmentExists(requestData.Name, db)
	if err != nil {
		return err
	}
	if rowExists {
		return errors.New("segment with this name already exists")
	}

	err = insertSegment(requestData.Name, db)
	if err != nil {
		return err
	}

	return err
}

func DeleteSegment(requestData dtos.UpdateSegment) error {
	db := database.Get()

	rowExists, err := checkSegmentExists(requestData.Name, db)
	if err != nil {
		return err
	}
	if !rowExists {
		return errors.New("segment with this name not exists")
	}

	err = deleteSegment(requestData.Name, db)
	if err != nil {
		return err
	}

	return err
}

func checkSegmentExists(name string, db *sql.DB) (bool, error) {
	var rowExists bool
	err := db.QueryRow("select exists(select true from segments where name=$1)", name).Scan(&rowExists)
	if err != nil {
		return rowExists, errors.New("query error")
	}
	return rowExists, nil
}

func insertSegment(name string, db *sql.DB) error {
	_, err := db.Exec("insert into segments (name) values ($1)", name)
	if err != nil {
		return errors.New("insert error")
	}
	return nil
}

func deleteSegment(name string, db *sql.DB) error {
	_, err := db.Exec("delete from segments values where name = $1", name)
	if err != nil {
		return errors.New("delete error")
	}
	return nil
}
