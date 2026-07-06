package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Sub           string    `gorm:"primaryKey" json:"sub"`
	Email         string    `json:"email"`
	DisplayName   string    `json:"displayName"`
	WalletAddress *string   `json:"walletAddress"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func GetUser(db *gorm.DB, sub string) (*User, error) {
	var user User
	err := db.First(&user, "sub = ?", sub).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByWalletAddress(db *gorm.DB, address string) (*User, error) {
	normalized := strings.ToLower(address)

	var user User
	err := db.First(&user, "wallet_address = ?", normalized).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpsertUserFromAuth0(db *gorm.DB, sub, email, name string) (*User, error) {
	existing, err := GetUser(db, sub)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		user := User{Sub: sub, Email: email, DisplayName: name}
		if err := db.Create(&user).Error; err != nil {
			return nil, err
		}
		return &user, nil
	}

	if err := db.Model(existing).Update("email", email).Error; err != nil {
		return nil, err
	}
	existing.Email = email

	return existing, nil
}

func UpdateDisplayName(db *gorm.DB, sub, displayName string) (*User, error) {
	existing, err := GetUser(db, sub)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("user not found")
	}

	if err := db.Model(existing).Update("display_name", displayName).Error; err != nil {
		return nil, err
	}
	existing.DisplayName = displayName

	return existing, nil
}

func SetWalletAddress(db *gorm.DB, sub, address string) (*User, error) {
	existing, err := GetUser(db, sub)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("user not found")
	}

	normalized := strings.ToLower(address)
	if err := db.Model(existing).Update("wallet_address", normalized).Error; err != nil {
		return nil, err
	}
	existing.WalletAddress = &normalized

	return existing, nil
}
