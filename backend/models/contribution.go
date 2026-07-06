package models

import (
	"math/big"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Contribution struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement"`
	CampaignID  uint64 `gorm:"index;not null"`
	Contributor string `gorm:"not null"`
	Amount      string `gorm:"not null"`
	BlockNumber uint64 `gorm:"not null"`
	TxHash      string `gorm:"uniqueIndex:idx_tx_log;not null"`
	LogIndex    uint   `gorm:"uniqueIndex:idx_tx_log;not null"`
	CreatedAt   time.Time
}

type ContributorSummary struct {
	Contributor string `json:"address"`
	TotalAmount string `json:"amount"`
}

func SaveContribution(db *gorm.DB, contribution *Contribution) error {
	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(contribution).Error
}

func GetContributorsForCampaign(db *gorm.DB, campaignID uint64) ([]ContributorSummary, error) {
	var contributions []Contribution
	if err := db.Where("campaign_id = ?", campaignID).Order("created_at asc").Find(&contributions).Error; err != nil {
		return nil, err
	}

	totals := make(map[string]*big.Int)
	order := make([]string, 0)

	for _, c := range contributions {
		amount, ok := new(big.Int).SetString(c.Amount, 10)
		if !ok {
			continue
		}

		if _, exists := totals[c.Contributor]; !exists {
			order = append(order, c.Contributor)
			totals[c.Contributor] = big.NewInt(0)
		}
		totals[c.Contributor].Add(totals[c.Contributor], amount)
	}

	summaries := make([]ContributorSummary, 0, len(order))
	for _, address := range order {
		summaries = append(summaries, ContributorSummary{
			Contributor: address,
			TotalAmount: totals[address].String(),
		})
	}

	return summaries, nil
}
