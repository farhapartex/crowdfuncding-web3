package services

import (
	"log"
	"math/big"
	"time"

	"gorm.io/gorm"

	"crowdfunding-backend/contract"
	"crowdfunding-backend/models"
)

const deadlineCheckInterval = 60 * time.Second
const autoArchiveNote = "Automatically archived because the campaign's funding deadline was reached."

func StartDeadlineArchiver(db *gorm.DB, crowdFunding *contract.CrowdFunding) {
	go func() {
		for {
			if err := archiveExpiredCampaigns(db, crowdFunding); err != nil {
				log.Printf("deadline archiver: %v", err)
			}
			time.Sleep(deadlineCheckInterval)
		}
	}()
}

func archiveExpiredCampaigns(db *gorm.DB, crowdFunding *contract.CrowdFunding) error {
	campaigns, err := models.ListPublishedCampaignsForArchiveCheck(db)
	if err != nil {
		return err
	}

	now := time.Now().Unix()

	for _, campaign := range campaigns {
		chainCampaign, err := crowdFunding.GetCampaign(nil, new(big.Int).SetUint64(*campaign.OnChainCampaignID))
		if err != nil {
			log.Printf("deadline archiver: failed to read on-chain campaign %d: %v", *campaign.OnChainCampaignID, err)
			continue
		}

		if chainCampaign.Deadline.Int64() > now {
			continue
		}

		if _, err := models.ArchiveCampaign(db, campaign.ID, autoArchiveNote); err != nil {
			log.Printf("deadline archiver: failed to archive campaign %d: %v", campaign.ID, err)
		}
	}

	return nil
}
