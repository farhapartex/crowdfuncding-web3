package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

const (
	CampaignStatusDraft     = "draft"
	CampaignStatusPublished = "published"
)

type Campaign struct {
	ID                uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	OwnerSub          string     `gorm:"index;not null" json:"ownerSub"`
	Country           string     `gorm:"not null" json:"country"`
	Title             string     `gorm:"not null" json:"title"`
	Description       string     `gorm:"type:text" json:"description"`
	TargetEth         string     `gorm:"not null" json:"targetEth"`
	FundraisingFor    string     `gorm:"not null" json:"fundraisingFor"`
	Status            string     `gorm:"not null;default:draft" json:"status"`
	WalletAddress     *string    `gorm:"index" json:"walletAddress"`
	OnChainCampaignID *uint64    `gorm:"index" json:"onChainCampaignId"`
	PublishedAt       *time.Time `json:"publishedAt"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}

func CreateCampaign(db *gorm.DB, ownerSub, country, title, description, targetEth, fundraisingFor string) (*Campaign, error) {
	campaign := Campaign{
		OwnerSub:       ownerSub,
		Country:        country,
		Title:          title,
		Description:    description,
		TargetEth:      targetEth,
		FundraisingFor: fundraisingFor,
		Status:         CampaignStatusDraft,
	}

	if err := db.Create(&campaign).Error; err != nil {
		return nil, err
	}

	return &campaign, nil
}

func GetCampaignByID(db *gorm.DB, id uint64) (*Campaign, error) {
	var campaign Campaign
	err := db.First(&campaign, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &campaign, nil
}

func ListCampaignsByOwner(db *gorm.DB, ownerSub string, offset, limit uint64) ([]Campaign, int64, error) {
	var total int64
	if err := db.Model(&Campaign{}).Where("owner_sub = ?", ownerSub).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var campaigns []Campaign
	err := db.Where("owner_sub = ?", ownerSub).
		Order("created_at desc").
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&campaigns).Error
	if err != nil {
		return nil, 0, err
	}

	return campaigns, total, nil
}
