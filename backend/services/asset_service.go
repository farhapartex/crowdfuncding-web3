package services

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"crowdfunding-backend/models"
)

const maxAssetUploadSize = 10 * 1024 * 1024

type AssetService struct {
	db        *gorm.DB
	r2Client  *s3.Client
	bucket    string
	publicURL string
}

func NewAssetService(db *gorm.DB, r2Client *s3.Client, bucket, publicURL string) *AssetService {
	return &AssetService{db: db, r2Client: r2Client, bucket: bucket, publicURL: publicURL}
}

func (s *AssetService) UploadAsset(ctx context.Context, sub string, fileHeader *multipart.FileHeader) (*models.Asset, error) {
	if s.bucket == "" || s.publicURL == "" {
		return nil, NewUnavailableError("asset storage is not configured yet")
	}
	if fileHeader.Size > maxAssetUploadSize {
		return nil, NewValidationError("file exceeds 10MB limit")
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return nil, NewValidationError("only image uploads are allowed")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)
	objectKey := fmt.Sprintf("uploads/campaign/covers/%s%s", uuid.NewString(), ext)

	if err := uploadObjectToR2(ctx, s.r2Client, s.bucket, objectKey, contentType, file, fileHeader.Size); err != nil {
		return nil, err
	}

	assetURL := fmt.Sprintf("%s/%s", strings.TrimRight(s.publicURL, "/"), objectKey)

	return models.CreateAsset(s.db, sub, s.bucket, objectKey, assetURL, contentType, fileHeader.Size)
}

func (s *AssetService) DeleteAssets(ctx context.Context, assets []models.Asset) {
	for _, asset := range assets {
		if err := deleteObjectFromR2(ctx, s.r2Client, asset.Bucket, asset.ObjectKey); err != nil {
			log.Printf("failed to delete r2 object %s/%s: %v", asset.Bucket, asset.ObjectKey, err)
		}
	}
}
