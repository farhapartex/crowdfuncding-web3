package models

import (
	"math/big"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	TransactionTypeContribution   = "contribution"
	TransactionTypeWithdraw       = "withdraw"
	TransactionTypeRefund         = "refund"
	TransactionTypeCreateCampaign = "create_campaign"
	TransactionTypeCloseCampaign  = "close_campaign"
)

type Transaction struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	CampaignID     uint64    `gorm:"index;not null" json:"campaignId"`
	Type           string    `gorm:"index;not null" json:"type"`
	Address        string    `gorm:"index;not null" json:"address"`
	Amount         string    `gorm:"not null;default:0" json:"amount"`
	GasFeeWei      *string   `json:"gasFeeWei"`
	BlockNumber    uint64    `gorm:"not null" json:"blockNumber"`
	BlockTimestamp time.Time `json:"blockTimestamp"`
	TxHash         string    `gorm:"uniqueIndex:idx_transactions_tx_log;not null" json:"txHash"`
	LogIndex       uint      `gorm:"uniqueIndex:idx_transactions_tx_log;not null" json:"logIndex"`
	CreatedAt      time.Time `json:"createdAt"`
}

type ContributorSummary struct {
	Contributor string `json:"address"`
	TotalAmount string `json:"amount"`
}

func SaveTransaction(db *gorm.DB, tx *Transaction) error {
	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(tx).Error
}

func GetTransactionsForCampaign(db *gorm.DB, campaignID uint64, offset, limit uint64) ([]Transaction, int64, error) {
	types := []string{TransactionTypeContribution, TransactionTypeWithdraw}

	var total int64
	if err := db.Model(&Transaction{}).
		Where("campaign_id = ? AND type IN ?", campaignID, types).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var transactions []Transaction
	err := db.Where("campaign_id = ? AND type IN ?", campaignID, types).
		Order("block_number desc").
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func GetTransactionsForWallet(db *gorm.DB, address string, offset, limit uint64) ([]Transaction, int64, error) {
	normalized := strings.ToLower(address)

	var total int64
	if err := db.Model(&Transaction{}).Where("address = ?", normalized).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var transactions []Transaction
	err := db.Where("address = ?", normalized).
		Order("block_number desc").
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func GetContributorsForCampaign(db *gorm.DB, campaignID uint64) ([]ContributorSummary, error) {
	var transactions []Transaction
	if err := db.Where("campaign_id = ? AND type = ?", campaignID, TransactionTypeContribution).
		Order("created_at asc").Find(&transactions).Error; err != nil {
		return nil, err
	}

	totals := make(map[string]*big.Int)
	order := make([]string, 0)

	for _, t := range transactions {
		amount, ok := new(big.Int).SetString(t.Amount, 10)
		if !ok {
			continue
		}

		if _, exists := totals[t.Address]; !exists {
			order = append(order, t.Address)
			totals[t.Address] = big.NewInt(0)
		}
		totals[t.Address].Add(totals[t.Address], amount)
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
