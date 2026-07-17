package models

import (
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

const (
	CampaignStatusDraft     = "draft"
	CampaignStatusPublished = "published"
	CampaignStatusArchived  = "archived"
)

var CampaignCategories = []string{
	"Medical & Health",
	"Education",
	"Community & Environment",
	"Animals & Pets",
	"Emergency Relief",
	"Other",
}

func IsValidCampaignCategory(category string) bool {
	for _, c := range CampaignCategories {
		if c == category {
			return true
		}
	}
	return false
}

type Campaign struct {
	ID                uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	OwnerSub          string     `gorm:"index;not null" json:"ownerSub"`
	Country           string     `gorm:"not null" json:"country"`
	Category          string     `gorm:"not null;default:Other" json:"category"`
	Title             string     `gorm:"not null" json:"title"`
	Description       string     `gorm:"type:text" json:"description"`
	TargetEth         string     `gorm:"not null" json:"targetEth"`
	DurationDays      uint32     `gorm:"not null;default:30" json:"durationDays"`
	FundraisingFor    string     `gorm:"not null" json:"fundraisingFor"`
	Status            string     `gorm:"not null;default:draft" json:"status"`
	WalletAddress     *string    `gorm:"index" json:"walletAddress"`
	OnChainCampaignID *uint64    `gorm:"uniqueIndex" json:"onChainCampaignId"`
	PublishedAt       *time.Time `json:"publishedAt"`
	ArchivedAt        *time.Time `json:"archivedAt"`
	ArchiveNote       *string    `json:"archiveNote"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}

func CreateCampaign(db *gorm.DB, ownerSub, country, category, title, description, targetEth string, durationDays uint32, fundraisingFor string) (*Campaign, error) {
	campaign := Campaign{
		OwnerSub:       ownerSub,
		Country:        country,
		Category:       category,
		Title:          title,
		Description:    description,
		TargetEth:      targetEth,
		DurationDays:   durationDays,
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

func GetCampaignByOnChainID(db *gorm.DB, onChainCampaignID uint64) (*Campaign, error) {
	var campaign Campaign
	err := db.First(&campaign, "on_chain_campaign_id = ?", onChainCampaignID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &campaign, nil
}

var (
	ErrCampaignNotDraft             = errors.New("campaign is not a draft")
	ErrOnChainCampaignAlreadyLinked = errors.New("this on-chain campaign is already linked to another draft")
	ErrCampaignNotPublished         = errors.New("campaign is not published")
)

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func PublishCampaign(db *gorm.DB, id uint64, walletAddress string, onChainCampaignID uint64) (*Campaign, error) {
	result := db.Model(&Campaign{}).
		Where("id = ? AND status = ?", id, CampaignStatusDraft).
		Updates(map[string]any{
			"status":               CampaignStatusPublished,
			"wallet_address":       strings.ToLower(walletAddress),
			"on_chain_campaign_id": onChainCampaignID,
			"published_at":         time.Now(),
		})
	if result.Error != nil {
		if isUniqueViolation(result.Error) {
			return nil, ErrOnChainCampaignAlreadyLinked
		}
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, ErrCampaignNotDraft
	}

	return GetCampaignByID(db, id)
}

func ArchiveCampaign(db *gorm.DB, id uint64, note string) (*Campaign, error) {
	result := db.Model(&Campaign{}).
		Where("id = ? AND status = ?", id, CampaignStatusPublished).
		Updates(map[string]any{
			"status":       CampaignStatusArchived,
			"archived_at":  time.Now(),
			"archive_note": note,
		})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, ErrCampaignNotPublished
	}

	return GetCampaignByID(db, id)
}

func ListPublishedCampaignsForArchiveCheck(db *gorm.DB) ([]Campaign, error) {
	var campaigns []Campaign
	err := db.Where("status = ? AND on_chain_campaign_id IS NOT NULL", CampaignStatusPublished).Find(&campaigns).Error
	if err != nil {
		return nil, err
	}

	return campaigns, nil
}

func DeleteCampaign(db *gorm.DB, id uint64, orphanAssetIDs []uint64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("campaign_id = ?", id).Delete(&CampaignAsset{}).Error; err != nil {
			return err
		}

		if len(orphanAssetIDs) > 0 {
			if err := tx.Where("id IN ?", orphanAssetIDs).Delete(&Asset{}).Error; err != nil {
				return err
			}
		}

		result := tx.Where("id = ? AND status = ?", id, CampaignStatusDraft).Delete(&Campaign{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrCampaignNotDraft
		}

		return nil
	})
}

func CountPublishedCampaigns(db *gorm.DB, category string) (int64, error) {
	query := db.Model(&Campaign{}).Where("status IN ?", []string{CampaignStatusPublished, CampaignStatusArchived})
	if category != "" {
		query = query.Where("category = ?", category)
	}

	var total int64
	err := query.Count(&total).Error
	return total, err
}

func ListPublishedCampaigns(db *gorm.DB, category string, offset, limit uint64) ([]Campaign, int64, error) {
	total, err := CountPublishedCampaigns(db, category)
	if err != nil {
		return nil, 0, err
	}

	query := db.Where("status IN ?", []string{CampaignStatusPublished, CampaignStatusArchived})
	if category != "" {
		query = query.Where("category = ?", category)
	}

	var campaigns []Campaign
	err = query.Order("published_at desc").Offset(int(offset)).Limit(int(limit)).Find(&campaigns).Error
	if err != nil {
		return nil, 0, err
	}

	return campaigns, total, nil
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
