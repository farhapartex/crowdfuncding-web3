package services

import (
	"time"

	"crowdfunding-backend/contract"
)

func campaignStatus(campaign contract.Campaign) string {
	if campaign.AmountRaised.Cmp(campaign.Goal) >= 0 {
		return "Successful"
	}
	if time.Now().Unix() >= campaign.Deadline.Int64() {
		return "Failed"
	}
	return "Active"
}
