package service

import (
	"avito-backend/database"
	"avito-backend/database/dbaccess"
	"avito-backend/dto"
	"errors"
)

func CreateSegment(requestData dto.UpdateSegment) error {
	db := database.Get()

	rowExists, err := dbaccess.IsSegmentExists(requestData.Name, db)
	if err != nil {
		return err
	}
	if rowExists {
		return errors.New("segment with this name already exists")
	}

	err = dbaccess.InsertSegment(requestData.Name, requestData.UserPercentage, db)
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
