package services

import (
	"log"
	"math/big"

	"gorm.io/gorm"

	"crowdfunding-backend/contract"
	"crowdfunding-backend/models"
)

type PublicCampaignService struct {
	db           *gorm.DB
	crowdFunding *contract.CrowdFunding
	idMask       *IDMaskService
}

func NewPublicCampaignService(db *gorm.DB, crowdFunding *contract.CrowdFunding, idMask *IDMaskService) *PublicCampaignService {
	return &PublicCampaignService{db: db, crowdFunding: crowdFunding, idMask: idMask}
}

type CampaignResponse struct {
	ID                string `json:"id"`
	OnChainCampaignID uint64 `json:"onChainCampaignId"`
	Owner             string `json:"owner"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	CurrencyMode      string `json:"currencyMode"`
	TokenAddress      string `json:"tokenAddress,omitempty"`
	TokenSymbol       string `json:"tokenSymbol,omitempty"`
	TokenDecimals     uint8  `json:"tokenDecimals,omitempty"`
	GoalEth           string `json:"goalEth"`
	GoalToken         string `json:"goalToken"`
	Deadline          string `json:"deadline"`
	AmountRaisedEth   string `json:"amountRaisedEth"`
	AmountRaisedToken string `json:"amountRaisedToken"`
	Withdrawn         bool   `json:"withdrawn"`
	Status            string `json:"status"`
	Country           string `json:"country"`
	Category          string `json:"category"`
	CoverURL          string `json:"coverUrl"`
	IsArchived        bool   `json:"isArchived"`
	ArchiveNote       string `json:"archiveNote,omitempty"`
}

func toCampaignResponse(maskedID string, onChainID uint64, campaign contract.Campaign, dbCampaign *models.Campaign, coverURL string) CampaignResponse {
	description := campaign.Description
	country := ""
	category := ""
	currencyMode := models.CurrencyModeEth
	tokenSymbol := ""
	var tokenDecimals uint8
	isArchived := false
	archiveNote := ""
	if dbCampaign != nil {
		description = dbCampaign.Description
		country = dbCampaign.Country
		category = dbCampaign.Category
		currencyMode = dbCampaign.CurrencyMode
		if dbCampaign.TokenSymbol != nil {
			tokenSymbol = *dbCampaign.TokenSymbol
		}
		if dbCampaign.TokenDecimals != nil {
			tokenDecimals = *dbCampaign.TokenDecimals
		}
		isArchived = dbCampaign.Status == models.CampaignStatusArchived
		if dbCampaign.ArchiveNote != nil {
			archiveNote = *dbCampaign.ArchiveNote
		}
	}

	return CampaignResponse{
		ID:                maskedID,
		OnChainCampaignID: onChainID,
		Owner:             campaign.Owner.Hex(),
		Title:             campaign.Title,
		Description:       description,
		CurrencyMode:      currencyMode,
		TokenAddress:      campaign.Token.Hex(),
		TokenSymbol:       tokenSymbol,
		TokenDecimals:     tokenDecimals,
		GoalEth:           campaign.GoalEth.String(),
		GoalToken:         campaign.GoalToken.String(),
		Deadline:          campaign.Deadline.String(),
		AmountRaisedEth:   campaign.AmountRaisedEth.String(),
		AmountRaisedToken: campaign.AmountRaisedToken.String(),
		Withdrawn:         campaign.Withdrawn,
		Status:            campaignStatus(campaign),
		Country:           country,
		Category:          category,
		CoverURL:          coverURL,
		IsArchived:        isArchived,
		ArchiveNote:       archiveNote,
	}
}

func (s *PublicCampaignService) CountPublished(category string) (int64, error) {
	if category != "" && !models.IsValidCampaignCategory(category) {
		return 0, NewValidationError("invalid category")
	}

	return models.CountPublishedCampaigns(s.db, category)
}

func (s *PublicCampaignService) ListPublished(category string, offset, limit uint64) (items []CampaignResponse, total int64, err error) {
	if category != "" && !models.IsValidCampaignCategory(category) {
		return nil, 0, NewValidationError("invalid category")
	}

	dbCampaigns, total, err := models.ListPublishedCampaigns(s.db, category, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	campaignIDs := make([]uint64, len(dbCampaigns))
	for i, dbCampaign := range dbCampaigns {
		campaignIDs[i] = dbCampaign.ID
	}

	coverURLs, err := models.GetCoverAssetsForCampaigns(s.db, campaignIDs)
	if err != nil {
		return nil, 0, err
	}

	items = make([]CampaignResponse, 0, len(dbCampaigns))
	for _, dbCampaign := range dbCampaigns {
		if dbCampaign.OnChainCampaignID == nil {
			continue
		}

		onChainID := *dbCampaign.OnChainCampaignID
		chainCampaign, err := s.crowdFunding.GetCampaign(nil, new(big.Int).SetUint64(onChainID))
		if err != nil {
			continue
		}

		maskedID, err := s.idMask.Mask(dbCampaign.ID)
		if err != nil {
			return nil, 0, err
		}

		items = append(items, toCampaignResponse(maskedID, onChainID, chainCampaign, &dbCampaign, coverURLs[dbCampaign.ID]))
	}

	return items, total, nil
}

func (s *PublicCampaignService) GetPublished(maskedID string) (*CampaignResponse, error) {
	dbID, err := s.idMask.Unmask(maskedID)
	if err != nil {
		return nil, NewValidationError("invalid campaign id")
	}

	dbCampaign, err := models.GetCampaignByID(s.db, dbID)
	if err != nil {
		return nil, err
	}
	isVisible := dbCampaign != nil &&
		(dbCampaign.Status == models.CampaignStatusPublished || dbCampaign.Status == models.CampaignStatusArchived) &&
		dbCampaign.OnChainCampaignID != nil
	if !isVisible {
		return nil, NewNotFoundError("campaign not found")
	}

	onChainID := *dbCampaign.OnChainCampaignID
	chainCampaign, err := s.crowdFunding.GetCampaign(nil, new(big.Int).SetUint64(onChainID))
	if err != nil {
		log.Printf("campaigns/%s: failed to read on-chain campaign %d: %v", maskedID, onChainID, err)
		return nil, NewUnavailableError("campaign data is temporarily unavailable, please try again shortly")
	}

	coverURLs, err := models.GetCoverAssetsForCampaigns(s.db, []uint64{dbCampaign.ID})
	if err != nil {
		return nil, err
	}

	response := toCampaignResponse(maskedID, onChainID, chainCampaign, dbCampaign, coverURLs[dbCampaign.ID])
	return &response, nil
}

type ContributorDTO struct {
	Address      string  `json:"address"`
	DisplayName  string  `json:"displayName"`
	Amount       string  `json:"amount"`
	TokenAddress *string `json:"tokenAddress,omitempty"`
}

func (s *PublicCampaignService) GetContributors(maskedID string) ([]ContributorDTO, error) {
	dbID, err := s.idMask.Unmask(maskedID)
	if err != nil {
		return nil, NewValidationError("invalid campaign id")
	}

	dbCampaign, err := models.GetCampaignByID(s.db, dbID)
	if err != nil {
		return nil, err
	}
	if dbCampaign == nil || dbCampaign.OnChainCampaignID == nil {
		return nil, NewNotFoundError("campaign not found")
	}

	summaries, err := models.GetContributorsForCampaign(s.db, *dbCampaign.OnChainCampaignID)
	if err != nil {
		return nil, err
	}

	response := make([]ContributorDTO, len(summaries))
	for i, summary := range summaries {
		displayName := ""
		if profile, err := models.GetProfile(s.db, summary.Contributor); err == nil {
			displayName = profile.DisplayName
		}

		response[i] = ContributorDTO{
			Address:      summary.Contributor,
			DisplayName:  displayName,
			Amount:       summary.TotalAmount,
			TokenAddress: summary.TokenAddress,
		}
	}

	return response, nil
}

func (s *PublicCampaignService) EnsureCommentable(maskedID string) error {
	dbID, err := s.idMask.Unmask(maskedID)
	if err != nil {
		return NewValidationError("invalid campaign id")
	}

	campaign, err := models.GetCampaignByID(s.db, dbID)
	if err != nil {
		return err
	}
	if campaign == nil {
		return NewNotFoundError("campaign not found")
	}

	switch campaign.Status {
	case models.CampaignStatusPublished:
		return nil
	case models.CampaignStatusArchived:
		return NewValidationError("this campaign is archived and no longer accepts comments")
	default:
		return NewValidationError("this campaign is not accepting comments")
	}
}

func (s *PublicCampaignService) GetCampaignTransactions(maskedID string, offset, limit uint64) (items []models.Transaction, total int64, err error) {
	dbID, err := s.idMask.Unmask(maskedID)
	if err != nil {
		return nil, 0, NewValidationError("invalid campaign id")
	}

	dbCampaign, err := models.GetCampaignByID(s.db, dbID)
	if err != nil {
		return nil, 0, err
	}
	if dbCampaign == nil || dbCampaign.OnChainCampaignID == nil {
		return nil, 0, NewNotFoundError("campaign not found")
	}

	return models.GetTransactionsForCampaign(s.db, *dbCampaign.OnChainCampaignID, offset, limit)
}
