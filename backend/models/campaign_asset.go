package models

import (
	"time"

	"gorm.io/gorm"
)

type CampaignAsset struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	CampaignID uint64    `gorm:"index;not null" json:"campaignId"`
	AssetID    uint64    `gorm:"index;not null" json:"assetId"`
	IsCover    bool      `gorm:"not null;default:false" json:"isCover"`
	Position   int       `gorm:"not null;default:0" json:"position"`
	CreatedAt  time.Time `json:"createdAt"`
}

type CampaignAssetView struct {
	AssetID     uint64 `json:"assetId"`
	URL         string `json:"url"`
	ContentType string `json:"contentType"`
	IsCover     bool   `json:"isCover"`
	Position    int    `json:"position"`
}

func AttachAssetsToCampaign(db *gorm.DB, campaignID uint64, assetIDs []uint64, coverAssetID uint64) error {
	links := make([]CampaignAsset, len(assetIDs))
	for i, assetID := range assetIDs {
		links[i] = CampaignAsset{
			CampaignID: campaignID,
			AssetID:    assetID,
			IsCover:    assetID == coverAssetID,
			Position:   i,
		}
	}

	return db.Create(&links).Error
}

func GetCoverAssetsForCampaigns(db *gorm.DB, campaignIDs []uint64) (map[uint64]string, error) {
	result := make(map[uint64]string, len(campaignIDs))
	if len(campaignIDs) == 0 {
		return result, nil
	}

	type row struct {
		CampaignID uint64
		URL        string
	}

	var rows []row
	err := db.Table("campaign_assets").
		Select("campaign_assets.campaign_id, assets.url").
		Joins("JOIN assets ON assets.id = campaign_assets.asset_id").
		Where("campaign_assets.campaign_id IN ? AND campaign_assets.is_cover = true", campaignIDs).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, r := range rows {
		result[r.CampaignID] = r.URL
	}

	return result, nil
}

func GetOrphanableAssetsForCampaign(db *gorm.DB, campaignID uint64) ([]Asset, error) {
	var assetIDs []uint64
	if err := db.Model(&CampaignAsset{}).Where("campaign_id = ?", campaignID).Pluck("asset_id", &assetIDs).Error; err != nil {
		return nil, err
	}
	if len(assetIDs) == 0 {
		return []Asset{}, nil
	}

	var sharedAssetIDs []uint64
	err := db.Model(&CampaignAsset{}).
		Where("asset_id IN ? AND campaign_id != ?", assetIDs, campaignID).
		Pluck("asset_id", &sharedAssetIDs).Error
	if err != nil {
		return nil, err
	}

	shared := make(map[uint64]bool, len(sharedAssetIDs))
	for _, id := range sharedAssetIDs {
		shared[id] = true
	}

	orphanIDs := make([]uint64, 0, len(assetIDs))
	for _, id := range assetIDs {
		if !shared[id] {
			orphanIDs = append(orphanIDs, id)
		}
	}
	if len(orphanIDs) == 0 {
		return []Asset{}, nil
	}

	var assets []Asset
	if err := db.Where("id IN ?", orphanIDs).Find(&assets).Error; err != nil {
		return nil, err
	}

	return assets, nil
}

func GetCampaignAssets(db *gorm.DB, campaignID uint64) ([]CampaignAssetView, error) {
	var views []CampaignAssetView
	err := db.Table("campaign_assets").
		Select("campaign_assets.asset_id, assets.url, assets.content_type, campaign_assets.is_cover, campaign_assets.position").
		Joins("JOIN assets ON assets.id = campaign_assets.asset_id").
		Where("campaign_assets.campaign_id = ?", campaignID).
		Order("campaign_assets.position asc").
		Scan(&views).Error
	if err != nil {
		return nil, err
	}

	return views, nil
}
