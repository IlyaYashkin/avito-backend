package service

import (
	"avito-backend/database"
	"avito-backend/dto"
	"database/sql"
	"errors"
)

func CreateSegment(requestData dto.UpdateSegment) error {
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

func DeleteSegment(requestData dto.UpdateSegment) error {
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
		return rowExists, err
	}
	return rowExists, nil
}

func insertSegment(name string, db *sql.DB) error {
	_, err := db.Exec("insert into segments (name) values ($1)", name)
	if err != nil {
		return err
	}
	return nil
}

func deleteSegment(name string, db *sql.DB) error {
	_, err := db.Exec("delete from segments values where name = $1", name)
	if err != nil {
		return err
	}
	return nil
}
