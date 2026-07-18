package services

import (
	"time"

	"crowdfunding-backend/contract"
)

func campaignGoalReached(campaign contract.Campaign) bool {
	return (campaign.GoalEth.Sign() > 0 && campaign.AmountRaisedEth.Cmp(campaign.GoalEth) >= 0) ||
		(campaign.GoalToken.Sign() > 0 && campaign.AmountRaisedToken.Cmp(campaign.GoalToken) >= 0)
}

func campaignStatus(campaign contract.Campaign) string {
	if campaignGoalReached(campaign) {
		return "Successful"
	}
	if time.Now().Unix() >= campaign.Deadline.Int64() {
		return "Failed"
	}
	return "Active"
}
