package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID         string `gorm:"primaryKey"`
	CampaignID string `gorm:"index;not null"`
	AuthorSub  string `gorm:"not null"`
	AuthorName string `gorm:"not null"`
	Text       string `gorm:"type:text;not null"`
	ParentID   *string
	CreatedAt  time.Time
}

func ListCommentsByCampaign(db *gorm.DB, campaignID string, offset, limit uint64) ([]Comment, int64, error) {
	var total int64
	if err := db.Model(&Comment{}).Where("campaign_id = ?", campaignID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var comments []Comment
	err := db.Where("campaign_id = ?", campaignID).
		Order("created_at desc").
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&comments).Error
	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}
