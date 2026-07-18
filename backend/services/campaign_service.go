package services

import (
	"context"
	"errors"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"

	"crowdfunding-backend/contract"
	"crowdfunding-backend/models"
)

type CampaignService struct {
	db           *gorm.DB
	crowdFunding *contract.CrowdFunding
	idMask       *IDMaskService
	assets       *AssetService
	tokens       *TokenService
}

func NewCampaignService(db *gorm.DB, crowdFunding *contract.CrowdFunding, idMask *IDMaskService, assets *AssetService, tokens *TokenService) *CampaignService {
	return &CampaignService{db: db, crowdFunding: crowdFunding, idMask: idMask, assets: assets, tokens: tokens}
}

type CreateCampaignInput struct {
	Country        string
	Category       string
	Title          string
	Description    string
	DurationDays   uint32
	FundraisingFor string
	AssetIDs       []uint64
	CurrencyMode   string
	TargetEth      string
	TokenAddress   string
	GoalToken      string
}

type PublishCampaignInput struct {
	WalletAddress     string
	OnChainCampaignID uint64
	TxHash            string
}

func (s *CampaignService) ListMyCampaigns(sub string, offset, limit uint64) (items []map[string]any, total int64, err error) {
	campaigns, total, err := models.ListCampaignsByOwner(s.db, sub, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	campaignIDs := make([]uint64, len(campaigns))
	for i, campaign := range campaigns {
		campaignIDs[i] = campaign.ID
	}

	coverURLs, err := models.GetCoverAssetsForCampaigns(s.db, campaignIDs)
	if err != nil {
		return nil, 0, err
	}

	items = make([]map[string]any, len(campaigns))
	for i, campaign := range campaigns {
		maskedID, err := s.idMask.Mask(campaign.ID)
		if err != nil {
			return nil, 0, err
		}

		items[i] = map[string]any{
			"id":             maskedID,
			"title":          campaign.Title,
			"description":    campaign.Description,
			"currencyMode":   campaign.CurrencyMode,
			"targetEth":      campaign.TargetEth,
			"tokenSymbol":    campaign.TokenSymbol,
			"goalToken":      campaign.GoalToken,
			"country":        campaign.Country,
			"category":       campaign.Category,
			"durationDays":   campaign.DurationDays,
			"fundraisingFor": campaign.FundraisingFor,
			"status":         campaign.Status,
			"coverUrl":       coverURLs[campaign.ID],
			"createdAt":      campaign.CreatedAt,
		}
	}

	return items, total, nil
}

func (s *CampaignService) resolveToken(tokenAddress string) (*SupportedToken, error) {
	if tokenAddress == "" {
		return nil, NewValidationError("tokenAddress is required")
	}

	token := s.tokens.Find(tokenAddress)
	if token == nil {
		return nil, NewValidationError("unsupported token address")
	}

	return token, nil
}

func (s *CampaignService) CreateCampaign(sub string, input CreateCampaignInput) (map[string]any, error) {
	if !models.IsValidCampaignCategory(input.Category) {
		return nil, NewValidationError("invalid category")
	}
	if !models.IsValidCurrencyMode(input.CurrencyMode) {
		return nil, NewValidationError("invalid currency mode")
	}

	var tokenAddress, tokenSymbol *string
	var tokenDecimals *uint8
	var goalToken *string

	switch input.CurrencyMode {
	case models.CurrencyModeEth:
		if input.TargetEth == "" {
			return nil, NewValidationError("targetEth is required")
		}
	case models.CurrencyModeToken:
		token, err := s.resolveToken(input.TokenAddress)
		if err != nil {
			return nil, err
		}
		if input.GoalToken == "" {
			return nil, NewValidationError("goalToken is required")
		}
		tokenAddress, tokenSymbol, tokenDecimals = &token.Address, &token.Symbol, &token.Decimals
		goalToken = &input.GoalToken
	case models.CurrencyModeBoth:
		if input.TargetEth == "" || input.GoalToken == "" {
			return nil, NewValidationError("targetEth and goalToken are both required")
		}
		token, err := s.resolveToken(input.TokenAddress)
		if err != nil {
			return nil, err
		}
		tokenAddress, tokenSymbol, tokenDecimals = &token.Address, &token.Symbol, &token.Decimals
		goalToken = &input.GoalToken
	}

	assets, err := models.GetAssetsByIDs(s.db, input.AssetIDs)
	if err != nil {
		return nil, err
	}
	if len(assets) != len(input.AssetIDs) {
		return nil, NewValidationError("one or more assets do not exist")
	}
	for _, asset := range assets {
		if asset.UploadedBy != sub {
			return nil, NewForbiddenError("one or more assets do not belong to you")
		}
	}

	campaign, err := models.CreateCampaign(s.db, models.CreateCampaignParams{
		OwnerSub:       sub,
		Country:        input.Country,
		Category:       input.Category,
		Title:          input.Title,
		Description:    input.Description,
		CurrencyMode:   input.CurrencyMode,
		TargetEth:      input.TargetEth,
		TokenAddress:   tokenAddress,
		TokenSymbol:    tokenSymbol,
		TokenDecimals:  tokenDecimals,
		GoalToken:      goalToken,
		DurationDays:   input.DurationDays,
		FundraisingFor: input.FundraisingFor,
	})
	if err != nil {
		return nil, err
	}

	if err := models.AttachAssetsToCampaign(s.db, campaign.ID, input.AssetIDs, input.AssetIDs[0]); err != nil {
		return nil, err
	}

	maskedID, err := s.idMask.Mask(campaign.ID)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"id":             maskedID,
		"ownerSub":       campaign.OwnerSub,
		"country":        campaign.Country,
		"category":       campaign.Category,
		"title":          campaign.Title,
		"description":    campaign.Description,
		"currencyMode":   campaign.CurrencyMode,
		"targetEth":      campaign.TargetEth,
		"tokenAddress":   campaign.TokenAddress,
		"tokenSymbol":    campaign.TokenSymbol,
		"tokenDecimals":  campaign.TokenDecimals,
		"goalToken":      campaign.GoalToken,
		"durationDays":   campaign.DurationDays,
		"fundraisingFor": campaign.FundraisingFor,
		"status":         campaign.Status,
		"createdAt":      campaign.CreatedAt,
	}, nil
}

func (s *CampaignService) getOwnedCampaign(sub, maskedID string) (*models.Campaign, error) {
	campaignID, err := s.idMask.Unmask(maskedID)
	if err != nil {
		return nil, NewValidationError("invalid campaign id")
	}

	campaign, err := models.GetCampaignByID(s.db, campaignID)
	if err != nil {
		return nil, err
	}
	if campaign == nil {
		return nil, NewNotFoundError("campaign not found")
	}
	if campaign.OwnerSub != sub {
		return nil, NewForbiddenError("not the owner of this campaign")
	}

	return campaign, nil
}

func (s *CampaignService) GetMyCampaign(sub, maskedID string) (map[string]any, error) {
	campaign, err := s.getOwnedCampaign(sub, maskedID)
	if err != nil {
		return nil, err
	}

	assets, err := models.GetCampaignAssets(s.db, campaign.ID)
	if err != nil {
		return nil, err
	}

	response := map[string]any{
		"id":                maskedID,
		"ownerSub":          campaign.OwnerSub,
		"country":           campaign.Country,
		"category":          campaign.Category,
		"title":             campaign.Title,
		"description":       campaign.Description,
		"currencyMode":      campaign.CurrencyMode,
		"targetEth":         campaign.TargetEth,
		"tokenAddress":      campaign.TokenAddress,
		"tokenSymbol":       campaign.TokenSymbol,
		"tokenDecimals":     campaign.TokenDecimals,
		"goalToken":         campaign.GoalToken,
		"durationDays":      campaign.DurationDays,
		"fundraisingFor":    campaign.FundraisingFor,
		"status":            campaign.Status,
		"walletAddress":     campaign.WalletAddress,
		"onChainCampaignId": campaign.OnChainCampaignID,
		"publishedAt":       campaign.PublishedAt,
		"archivedAt":        campaign.ArchivedAt,
		"archiveNote":       campaign.ArchiveNote,
		"createdAt":         campaign.CreatedAt,
		"assets":            assets,
	}

	if campaign.OnChainCampaignID != nil {
		chainCampaign, err := s.crowdFunding.GetCampaign(nil, new(big.Int).SetUint64(*campaign.OnChainCampaignID))
		if err != nil {
			log.Printf("my-campaigns/%s: failed to read on-chain campaign %d: %v", maskedID, *campaign.OnChainCampaignID, err)
			response["onChainAvailable"] = false
		} else {
			response["onChainAvailable"] = true
			response["goalEth"] = chainCampaign.GoalEth.String()
			response["goalTokenOnChain"] = chainCampaign.GoalToken.String()
			response["amountRaisedEth"] = chainCampaign.AmountRaisedEth.String()
			response["amountRaisedToken"] = chainCampaign.AmountRaisedToken.String()
			response["deadline"] = chainCampaign.Deadline.String()
			response["withdrawn"] = chainCampaign.Withdrawn
			response["onChainStatus"] = campaignStatus(chainCampaign)
		}
	}

	return response, nil
}

func (s *CampaignService) DeleteCampaign(ctx context.Context, sub, maskedID string) error {
	campaign, err := s.getOwnedCampaign(sub, maskedID)
	if err != nil {
		return err
	}
	if campaign.Status != models.CampaignStatusDraft {
		return NewValidationError("only draft campaigns can be deleted")
	}

	orphanAssets, err := models.GetOrphanableAssetsForCampaign(s.db, campaign.ID)
	if err != nil {
		return err
	}

	orphanAssetIDs := make([]uint64, len(orphanAssets))
	for i, asset := range orphanAssets {
		orphanAssetIDs[i] = asset.ID
	}

	if err := models.DeleteCampaign(s.db, campaign.ID, orphanAssetIDs); err != nil {
		if errors.Is(err, models.ErrCampaignNotDraft) {
			return NewValidationError("only draft campaigns can be deleted")
		}
		return err
	}

	s.assets.DeleteAssets(ctx, orphanAssets)

	return nil
}

func (s *CampaignService) PublishCampaign(sub, maskedID string, input PublishCampaignInput) (map[string]any, error) {
	campaign, err := s.getOwnedCampaign(sub, maskedID)
	if err != nil {
		return nil, err
	}
	if campaign.Status != models.CampaignStatusDraft {
		return nil, NewValidationError("only draft campaigns can be published")
	}

	onChainCampaign, err := s.crowdFunding.GetCampaign(nil, new(big.Int).SetUint64(input.OnChainCampaignID))
	if err != nil {
		return nil, NewValidationError("could not find this campaign on-chain")
	}
	if !strings.EqualFold(onChainCampaign.Owner.Hex(), input.WalletAddress) {
		return nil, NewValidationError("on-chain campaign owner does not match the provided wallet address")
	}
	if err := s.verifyOnChainCurrencyMatches(campaign, onChainCampaign); err != nil {
		return nil, err
	}

	updated, err := models.PublishCampaign(s.db, campaign.ID, input.WalletAddress, input.OnChainCampaignID)
	if err != nil {
		if errors.Is(err, models.ErrCampaignNotDraft) {
			return nil, NewValidationError("only draft campaigns can be published")
		}
		if errors.Is(err, models.ErrOnChainCampaignAlreadyLinked) {
			return nil, NewConflictError(err.Error())
		}
		return nil, err
	}

	updatedMaskedID, err := s.idMask.Mask(updated.ID)
	if err != nil {
		return nil, err
	}

	log.Printf("campaign %d published on-chain as %d (tx %s)", updated.ID, input.OnChainCampaignID, input.TxHash)

	return map[string]any{
		"id":                updatedMaskedID,
		"status":            updated.Status,
		"walletAddress":     updated.WalletAddress,
		"onChainCampaignId": updated.OnChainCampaignID,
		"publishedAt":       updated.PublishedAt,
	}, nil
}

func currencyModeToOnChain(mode string) (uint8, bool) {
	switch mode {
	case models.CurrencyModeEth:
		return 0, true
	case models.CurrencyModeToken:
		return 1, true
	case models.CurrencyModeBoth:
		return 2, true
	default:
		return 0, false
	}
}

func (s *CampaignService) verifyOnChainCurrencyMatches(dbCampaign *models.Campaign, onChainCampaign contract.Campaign) error {
	expectedMode, ok := currencyModeToOnChain(dbCampaign.CurrencyMode)
	if !ok {
		return NewValidationError("invalid currency mode")
	}
	if onChainCampaign.CurrencyMode != expectedMode {
		return NewValidationError("on-chain currency mode does not match the draft campaign")
	}

	if dbCampaign.CurrencyMode == models.CurrencyModeEth {
		if onChainCampaign.Token != (common.Address{}) {
			return NewValidationError("on-chain token address does not match the draft campaign")
		}
		return nil
	}

	expectedToken := ""
	if dbCampaign.TokenAddress != nil {
		expectedToken = *dbCampaign.TokenAddress
	}
	if !strings.EqualFold(onChainCampaign.Token.Hex(), expectedToken) {
		return NewValidationError("on-chain token address does not match the draft campaign")
	}

	return nil
}

func (s *CampaignService) ArchiveCampaign(sub, maskedID, note string) (map[string]any, error) {
	note = strings.TrimSpace(note)
	if note == "" {
		return nil, NewValidationError("a note is required when archiving a campaign")
	}

	campaign, err := s.getOwnedCampaign(sub, maskedID)
	if err != nil {
		return nil, err
	}
	if campaign.Status != models.CampaignStatusPublished {
		return nil, NewValidationError("only published campaigns can be archived")
	}

	updated, err := models.ArchiveCampaign(s.db, campaign.ID, note)
	if err != nil {
		if errors.Is(err, models.ErrCampaignNotPublished) {
			return nil, NewValidationError("only published campaigns can be archived")
		}
		return nil, err
	}

	return map[string]any{
		"id":          maskedID,
		"status":      updated.Status,
		"archivedAt":  updated.ArchivedAt,
		"archiveNote": updated.ArchiveNote,
	}, nil
}
