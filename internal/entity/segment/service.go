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

	rowExists, err := IsSegExists(requestData.Name, tx)
	if err != nil {
		return err
	}
	if rowExists {
		return errors.New("segment with this name already exists")
	}

	segmentId, err := InsSeg(requestData.Name, tx)
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

	rowExists, err := IsSegExists(requestData.Name, db)
	if err != nil {
		return err
	}
	if !rowExists {
		return errors.New("segment with this name not exists")
	}

	err = DelSeg(requestData.Name, db)
	if err != nil {
		return err
	}

	return err
}
