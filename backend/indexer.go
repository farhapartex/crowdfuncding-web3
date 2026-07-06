package main

import (
	"context"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"

	"crowdfunding-backend/contract"
	"crowdfunding-backend/models"
)

const contributionPollInterval = 5 * time.Second

func startContributionIndexer(db *gorm.DB, crowdFunding *contract.CrowdFunding, client *ethclient.Client) {
	go func() {
		for {
			if err := pollContributions(db, crowdFunding, client); err != nil {
				log.Printf("contribution indexer: %v", err)
			}
			time.Sleep(contributionPollInterval)
		}
	}()
}

func pollContributions(db *gorm.DB, crowdFunding *contract.CrowdFunding, client *ethclient.Client) error {
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

	iterator, err := crowdFunding.FilterContributionMade(opts, nil, nil)
	if err != nil {
		return err
	}
	defer iterator.Close()

	for iterator.Next() {
		event := iterator.Event
		contribution := &models.Contribution{
			CampaignID:  event.CampaignId.Uint64(),
			Contributor: event.Contributor.Hex(),
			Amount:      event.Amount.String(),
			BlockNumber: event.Raw.BlockNumber,
			TxHash:      event.Raw.TxHash.Hex(),
			LogIndex:    event.Raw.Index,
		}
		if err := models.SaveContribution(db, contribution); err != nil {
			return err
		}
	}
	if err := iterator.Error(); err != nil {
		return err
	}

	return models.SetLastProcessedBlock(db, latestBlock)
}
