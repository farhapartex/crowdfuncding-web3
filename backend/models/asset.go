package models

import (
	"time"

	"gorm.io/gorm"
)

type Asset struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UploadedBy  string    `gorm:"not null" json:"uploadedBy"`
	Bucket      string    `gorm:"not null" json:"bucket"`
	ObjectKey   string    `gorm:"not null" json:"objectKey"`
	URL         string    `gorm:"not null" json:"url"`
	ContentType string    `gorm:"not null" json:"contentType"`
	SizeBytes   int64     `gorm:"not null" json:"sizeBytes"`
	CreatedAt   time.Time `json:"createdAt"`
}

func CreateAsset(db *gorm.DB, uploadedBy, bucket, objectKey, url, contentType string, sizeBytes int64) (*Asset, error) {
	asset := Asset{
		UploadedBy:  uploadedBy,
		Bucket:      bucket,
		ObjectKey:   objectKey,
		URL:         url,
		ContentType: contentType,
		SizeBytes:   sizeBytes,
	}

	if err := db.Create(&asset).Error; err != nil {
		return nil, err
	}

	return &asset, nil
}

func GetAssetsByIDs(db *gorm.DB, ids []uint64) ([]Asset, error) {
	if len(ids) == 0 {
		return []Asset{}, nil
	}

	var assets []Asset
	if err := db.Where("id IN ?", ids).Find(&assets).Error; err != nil {
		return nil, err
	}

	return assets, nil
}
