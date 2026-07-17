package services

import (
	"gorm.io/gorm"

	"crowdfunding-backend/models"
)

type WalletService struct {
	db *gorm.DB
}

func NewWalletService(db *gorm.DB) *WalletService {
	return &WalletService{db: db}
}

func (s *WalletService) GetTransactions(address string, offset, limit uint64) ([]models.Transaction, int64, error) {
	return models.GetTransactionsForWallet(s.db, address, offset, limit)
}
