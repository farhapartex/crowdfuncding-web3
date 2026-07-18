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
	TokenAddress   *string   `gorm:"index" json:"tokenAddress"`
	GasFeeWei      *string   `json:"gasFeeWei"`
	BlockNumber    uint64    `gorm:"not null" json:"blockNumber"`
	BlockTimestamp time.Time `json:"blockTimestamp"`
	TxHash         string    `gorm:"uniqueIndex:idx_transactions_tx_log;not null" json:"txHash"`
	LogIndex       uint      `gorm:"uniqueIndex:idx_transactions_tx_log;not null" json:"logIndex"`
	CreatedAt      time.Time `json:"createdAt"`
}

type ContributorSummary struct {
	Contributor  string  `json:"address"`
	TotalAmount  string  `json:"amount"`
	TokenAddress *string `json:"tokenAddress"`
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

type contributorKey struct {
	address string
	token   string
}

func GetContributorsForCampaign(db *gorm.DB, campaignID uint64) ([]ContributorSummary, error) {
	var transactions []Transaction
	if err := db.Where("campaign_id = ? AND type = ?", campaignID, TransactionTypeContribution).
		Order("created_at asc").Find(&transactions).Error; err != nil {
		return nil, err
	}

	// Contributions in different currencies (ETH vs. an ERC20 token) are never
	// summed together — a "Both" mode campaign tracks each currency's totals
	// independently, so contributor totals must too.
	totals := make(map[contributorKey]*big.Int)
	order := make([]contributorKey, 0)

	for _, t := range transactions {
		amount, ok := new(big.Int).SetString(t.Amount, 10)
		if !ok {
			continue
		}

		tokenKey := ""
		if t.TokenAddress != nil {
			tokenKey = strings.ToLower(*t.TokenAddress)
		}
		key := contributorKey{address: t.Address, token: tokenKey}

		if _, exists := totals[key]; !exists {
			order = append(order, key)
			totals[key] = big.NewInt(0)
		}
		totals[key].Add(totals[key], amount)
	}

	summaries := make([]ContributorSummary, 0, len(order))
	for _, key := range order {
		var tokenAddress *string
		if key.token != "" {
			token := key.token
			tokenAddress = &token
		}

		summaries = append(summaries, ContributorSummary{
			Contributor:  key.address,
			TotalAmount:  totals[key].String(),
			TokenAddress: tokenAddress,
		})
	}

	return summaries, nil
}
