package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Profile struct {
	Address     string `gorm:"primaryKey" json:"address"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func GetProfile(db *gorm.DB, address string) (*Profile, error) {
	normalized := strings.ToLower(address)

	var profile Profile
	err := db.First(&profile, "address = ?", normalized).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Profile{Address: normalized}, nil
	}
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func UpsertProfile(db *gorm.DB, address, displayName, email string) (*Profile, error) {
	profile := Profile{
		Address:     strings.ToLower(address),
		DisplayName: displayName,
		Email:       email,
	}

	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "address"}},
		DoUpdates: clause.AssignmentColumns([]string{"display_name", "email", "updated_at"}),
	}).Create(&profile).Error
	if err != nil {
		return nil, err
	}

	return &profile, nil
}
