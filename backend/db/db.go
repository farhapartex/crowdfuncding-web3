package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"crowdfunding-backend/models"
)

func Connect(databaseURL string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Profile{},
		&models.User{},
		&models.Campaign{},
		&models.Asset{},
		&models.CampaignAsset{},
		&models.Contribution{},
		&models.IndexerState{},
	)
}
