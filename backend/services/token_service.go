package services

import (
	"encoding/json"
	"strings"
)

type SupportedToken struct {
	Symbol   string `json:"symbol"`
	Address  string `json:"address"`
	Decimals uint8  `json:"decimals"`
}

type TokenService struct {
	tokens []SupportedToken
}

func NewTokenService(rawJSON string) (*TokenService, error) {
	var tokens []SupportedToken
	if rawJSON != "" {
		if err := json.Unmarshal([]byte(rawJSON), &tokens); err != nil {
			return nil, err
		}
	}

	return &TokenService{tokens: tokens}, nil
}

func (s *TokenService) List() []SupportedToken {
	return s.tokens
}

func (s *TokenService) Find(address string) *SupportedToken {
	for _, t := range s.tokens {
		if strings.EqualFold(t.Address, address) {
			return &t
		}
	}
	return nil
}
