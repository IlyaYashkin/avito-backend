package segment

import (
	"avito-backend/internal/database"
	"avito-backend/internal/entity/segmentpercentage"
	"errors"
)

func createSegment(requestData RequestUpdateSegment) error {
	db := database.Get()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	rowExists, err := IsSegmentExists(requestData.Name, tx)
	if err != nil {
		return err
	}
	if rowExists {
		return errors.New("segment with this name already exists")
	}

	segmentId, err := InsertSegment(requestData.Name, tx)
	if err != nil {
		return err
	}
	err = segmentpercentage.InsertSegmentPercentage(segmentId, requestData.UserPercentage, tx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}

func deleteSegment(requestData RequestUpdateSegment) error {
	db := database.Get()

	rowExists, err := IsSegmentExists(requestData.Name, db)
	if err != nil {
		return err
	}
	if !rowExists {
		return errors.New("segment with this name not exists")
	}

	err = DelSegment(requestData.Name, db)
	if err != nil {
		return err
	}

	return err
}
