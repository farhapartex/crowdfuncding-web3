package services

import (
	"context"
	"errors"
	"log"
	"math/big"
	"strings"

	"gorm.io/gorm"

	"crowdfunding-backend/contract"
	"crowdfunding-backend/models"
)

type CampaignService struct {
	db           *gorm.DB
	crowdFunding *contract.CrowdFunding
	idMask       *IDMaskService
	assets       *AssetService
}

func NewCampaignService(db *gorm.DB, crowdFunding *contract.CrowdFunding, idMask *IDMaskService, assets *AssetService) *CampaignService {
	return &CampaignService{db: db, crowdFunding: crowdFunding, idMask: idMask, assets: assets}
}

type CreateCampaignInput struct {
	Country        string
	Category       string
	Title          string
	Description    string
	TargetEth      string
	DurationDays   uint32
	FundraisingFor string
	AssetIDs       []uint64
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
			"targetEth":      campaign.TargetEth,
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

func (s *CampaignService) CreateCampaign(sub string, input CreateCampaignInput) (map[string]any, error) {
	if !models.IsValidCampaignCategory(input.Category) {
		return nil, NewValidationError("invalid category")
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

	campaign, err := models.CreateCampaign(s.db, sub, input.Country, input.Category, input.Title, input.Description, input.TargetEth, input.DurationDays, input.FundraisingFor)
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
		"targetEth":      campaign.TargetEth,
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
		"targetEth":         campaign.TargetEth,
		"durationDays":      campaign.DurationDays,
		"fundraisingFor":    campaign.FundraisingFor,
		"status":            campaign.Status,
		"walletAddress":     campaign.WalletAddress,
		"onChainCampaignId": campaign.OnChainCampaignID,
		"publishedAt":       campaign.PublishedAt,
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
			response["goal"] = chainCampaign.Goal.String()
			response["amountRaised"] = chainCampaign.AmountRaised.String()
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
