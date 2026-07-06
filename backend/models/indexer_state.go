package models

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IndexerState struct {
	ID                 uint64 `gorm:"primaryKey"`
	LastProcessedBlock uint64
}

func GetLastProcessedBlock(db *gorm.DB) (uint64, error) {
	var state IndexerState
	err := db.First(&state, "id = ?", 1).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return state.LastProcessedBlock, nil
}

func SetLastProcessedBlock(db *gorm.DB, block uint64) error {
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"last_processed_block"}),
	}).Create(&IndexerState{ID: 1, LastProcessedBlock: block}).Error
}
