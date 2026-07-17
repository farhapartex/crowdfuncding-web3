package services

import (
	"gorm.io/gorm"

	"crowdfunding-backend/models"
)

type ProfileService struct {
	db *gorm.DB
}

func NewProfileService(db *gorm.DB) *ProfileService {
	return &ProfileService{db: db}
}

func (s *ProfileService) GetProfile(address string) (*models.Profile, error) {
	return models.GetProfile(s.db, address)
}

func (s *ProfileService) UpsertProfile(address, displayName, email string) (*models.Profile, error) {
	return models.UpsertProfile(s.db, address, displayName, email)
}
