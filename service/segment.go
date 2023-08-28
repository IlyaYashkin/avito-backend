package service

import (
	"avito-backend/database"
	"avito-backend/database/dbaccess"
	"avito-backend/dto"
	"errors"
)

func CreateSegment(requestData dto.UpdateSegment) error {
	db := database.Get()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	rowExists, err := dbaccess.IsSegmentExists(requestData.Name, tx)
	if err != nil {
		return err
	}
	if rowExists {
		return errors.New("segment with this name already exists")
	}

	segmentId, err := dbaccess.InsertSegment(requestData.Name, tx)
	if err != nil {
		return err
	}
	err = dbaccess.InsertSegmentPercentage(segmentId, requestData.UserPercentage, tx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}

func DeleteSegment(requestData dto.UpdateSegment) error {
	db := database.Get()

	rowExists, err := dbaccess.IsSegmentExists(requestData.Name, db)
	if err != nil {
		return err
	}
	if !rowExists {
		return errors.New("segment with this name not exists")
	}

	err = dbaccess.DeleteSegment(requestData.Name, db)
	if err != nil {
		return err
	}

	return err
}
