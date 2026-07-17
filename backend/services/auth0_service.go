package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"crowdfunding-backend/models"
)

type Auth0Service struct {
	db       *gorm.DB
	domain   string
	audience string
	keyFunc  jwt.Keyfunc
	issuer   string
}

func NewAuth0Service(db *gorm.DB, domain, audience string) (*Auth0Service, error) {
	jwksURL := fmt.Sprintf("https://%s/.well-known/jwks.json", domain)

	k, err := keyfunc.NewDefaultCtx(context.Background(), []string{jwksURL})
	if err != nil {
		return nil, err
	}

	return &Auth0Service{
		db:       db,
		domain:   domain,
		audience: audience,
		keyFunc:  k.Keyfunc,
		issuer:   fmt.Sprintf("https://%s/", domain),
	}, nil
}

func (s *Auth0Service) ValidateToken(tokenString string) (sub string, err error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, s.keyFunc,
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithIssuer(s.issuer),
		jwt.WithAudience(s.audience),
	)
	if err != nil || !token.Valid {
		return "", NewUnauthorizedError("invalid token")
	}

	return claims.Subject, nil
}

func (s *Auth0Service) FetchUserInfo(accessToken string) (email, name string, err error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s/userinfo", s.domain), nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("userinfo request failed with status %d", resp.StatusCode)
	}

	var payload struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", "", err
	}

	return payload.Email, payload.Name, nil
}

func (s *Auth0Service) SyncUser(sub, accessToken string) (*models.User, error) {
	email, name, err := s.FetchUserInfo(accessToken)
	if err != nil {
		return nil, err
	}

	return models.UpsertUserFromAuth0(s.db, sub, email, name)
}

func (s *Auth0Service) GetUser(sub string) (*models.User, error) {
	user, err := models.GetUser(s.db, sub)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, NewNotFoundError("user not found")
	}

	return user, nil
}
