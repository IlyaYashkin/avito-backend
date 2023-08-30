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

	segmentRepo := NewSegmentRepository(tx)
	segmentPercentageRepo := segmentpercentage.NewSegmentPercentageRepository(tx)

	rowExists, err := segmentRepo.IsExistsByName(requestData.Name)
	if err != nil {
		return err
	}
	if rowExists {
		return errors.New("segment with this name already exists")
	}
	segmentId, err := segmentRepo.Save(requestData.Name)
	if err != nil {
		return err
	}
	err = segmentPercentageRepo.Save(segmentId, requestData.UserPercentage)
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

	repo := NewSegmentRepository(db)

	rowExists, err := repo.IsExistsByName(requestData.Name)
	if err != nil {
		return err
	}
	if !rowExists {
		return errors.New("segment with this name not exists")
	}

	err = repo.DeleteByName(requestData.Name)
	if err != nil {
		return err
	}

	return err
}
