package services

import (
	"context"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"

	"crowdfunding-backend/contract"
	"crowdfunding-backend/models"
)

const transactionPollInterval = 5 * time.Second

func StartTransactionIndexer(db *gorm.DB, crowdFunding *contract.CrowdFunding, client *ethclient.Client) {
	go func() {
		for {
			if err := pollTransactions(db, crowdFunding, client); err != nil {
				log.Printf("transaction indexer: %v", err)
			}
			time.Sleep(transactionPollInterval)
		}
	}()
}

type blockTimestampCache struct {
	client *ethclient.Client
	cache  map[uint64]time.Time
}

func newBlockTimestampCache(client *ethclient.Client) *blockTimestampCache {
	return &blockTimestampCache{client: client, cache: make(map[uint64]time.Time)}
}

func (c *blockTimestampCache) get(ctx context.Context, blockNumber uint64) (time.Time, error) {
	if ts, ok := c.cache[blockNumber]; ok {
		return ts, nil
	}

	header, err := c.client.HeaderByNumber(ctx, new(big.Int).SetUint64(blockNumber))
	if err != nil {
		return time.Time{}, err
	}

	ts := time.Unix(int64(header.Time), 0)
	c.cache[blockNumber] = ts

	return ts, nil
}

func pollTransactions(db *gorm.DB, crowdFunding *contract.CrowdFunding, client *ethclient.Client) error {
	ctx := context.Background()

	latestBlock, err := client.BlockNumber(ctx)
	if err != nil {
		return err
	}

	lastProcessed, err := models.GetLastProcessedBlock(db)
	if err != nil {
		return err
	}

	startBlock := lastProcessed
	if lastProcessed > 0 {
		startBlock = lastProcessed + 1
	}

	if latestBlock < startBlock {
		return nil
	}

	opts := &bind.FilterOpts{Start: startBlock, End: &latestBlock, Context: ctx}
	timestamps := newBlockTimestampCache(client)

	if err := indexContributions(ctx, db, crowdFunding, opts, timestamps); err != nil {
		return err
	}
	if err := indexWithdrawals(ctx, db, crowdFunding, opts, timestamps); err != nil {
		return err
	}
	if err := indexRefunds(ctx, db, crowdFunding, opts, timestamps); err != nil {
		return err
	}

	return models.SetLastProcessedBlock(db, latestBlock)
}

func indexContributions(ctx context.Context, db *gorm.DB, crowdFunding *contract.CrowdFunding, opts *bind.FilterOpts, timestamps *blockTimestampCache) error {
	iterator, err := crowdFunding.FilterContributionMade(opts, nil, nil)
	if err != nil {
		return err
	}
	defer iterator.Close()

	for iterator.Next() {
		event := iterator.Event
		blockTime, err := timestamps.get(ctx, event.Raw.BlockNumber)
		if err != nil {
			return err
		}

		tx := &models.Transaction{
			CampaignID:     event.CampaignId.Uint64(),
			Type:           models.TransactionTypeContribution,
			Address:        strings.ToLower(event.Contributor.Hex()),
			Amount:         event.Amount.String(),
			BlockNumber:    event.Raw.BlockNumber,
			BlockTimestamp: blockTime,
			TxHash:         event.Raw.TxHash.Hex(),
			LogIndex:       event.Raw.Index,
		}
		if err := models.SaveTransaction(db, tx); err != nil {
			return err
		}
	}

	return iterator.Error()
}

func indexWithdrawals(ctx context.Context, db *gorm.DB, crowdFunding *contract.CrowdFunding, opts *bind.FilterOpts, timestamps *blockTimestampCache) error {
	iterator, err := crowdFunding.FilterFundsWithdrawn(opts, nil, nil)
	if err != nil {
		return err
	}
	defer iterator.Close()

	for iterator.Next() {
		event := iterator.Event
		blockTime, err := timestamps.get(ctx, event.Raw.BlockNumber)
		if err != nil {
			return err
		}

		tx := &models.Transaction{
			CampaignID:     event.CampaignId.Uint64(),
			Type:           models.TransactionTypeWithdraw,
			Address:        strings.ToLower(event.Owner.Hex()),
			Amount:         event.Amount.String(),
			BlockNumber:    event.Raw.BlockNumber,
			BlockTimestamp: blockTime,
			TxHash:         event.Raw.TxHash.Hex(),
			LogIndex:       event.Raw.Index,
		}
		if err := models.SaveTransaction(db, tx); err != nil {
			return err
		}
	}

	return iterator.Error()
}

func indexRefunds(ctx context.Context, db *gorm.DB, crowdFunding *contract.CrowdFunding, opts *bind.FilterOpts, timestamps *blockTimestampCache) error {
	iterator, err := crowdFunding.FilterContributionRefunded(opts, nil, nil)
	if err != nil {
		return err
	}
	defer iterator.Close()

	for iterator.Next() {
		event := iterator.Event
		blockTime, err := timestamps.get(ctx, event.Raw.BlockNumber)
		if err != nil {
			return err
		}

		tx := &models.Transaction{
			CampaignID:     event.CampaignId.Uint64(),
			Type:           models.TransactionTypeRefund,
			Address:        strings.ToLower(event.Contributor.Hex()),
			Amount:         event.Amount.String(),
			BlockNumber:    event.Raw.BlockNumber,
			BlockTimestamp: blockTime,
			TxHash:         event.Raw.TxHash.Hex(),
			LogIndex:       event.Raw.Index,
		}
		if err := models.SaveTransaction(db, tx); err != nil {
			return err
		}
	}

	return iterator.Error()
}
