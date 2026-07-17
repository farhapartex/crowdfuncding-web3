package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt/v5"
)

const sessionTTL = time.Hour

type AuthService struct {
	secret []byte
	nonces *NonceStore
}

func NewAuthService(secret []byte) *AuthService {
	return &AuthService{secret: secret, nonces: NewNonceStore()}
}

func (s *AuthService) IssueNonce(address string) (nonce, message string, err error) {
	if address == "" {
		return "", "", NewValidationError("address is required")
	}

	nonce, err = s.nonces.Issue(address)
	if err != nil {
		return "", "", err
	}

	return nonce, buildSignInMessage(nonce), nil
}

func (s *AuthService) VerifyAndIssueSession(address, signature string) (token, recoveredAddress string, err error) {
	nonce, ok := s.nonces.PeekAndDelete(address)
	if !ok {
		return "", "", NewUnauthorizedError("nonce not found or expired, request a new one")
	}

	expectedMessage := buildSignInMessage(nonce)
	recovered, err := recoverAddress(expectedMessage, signature)
	if err != nil {
		return "", "", NewUnauthorizedError("invalid signature")
	}

	if !strings.EqualFold(recovered.Hex(), address) {
		return "", "", NewUnauthorizedError("signature does not match address")
	}

	token, err = s.generateJWT(recovered.Hex())
	if err != nil {
		return "", "", err
	}

	return token, recovered.Hex(), nil
}

func (s *AuthService) ParseSession(tokenString string) (address string, err error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return s.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	address, ok = claims["address"].(string)
	if !ok {
		return "", errors.New("token missing address claim")
	}

	return address, nil
}

func (s *AuthService) generateJWT(address string) (string, error) {
	claims := jwt.MapClaims{
		"address": address,
		"exp":     time.Now().Add(sessionTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func buildSignInMessage(nonce string) string {
	return fmt.Sprintf("Sign this message to log in to Crowd Funding.\n\nNonce: %s", nonce)
}

func recoverAddress(message string, signatureHex string) (common.Address, error) {
	sig, err := hexutil.Decode(signatureHex)
	if err != nil {
		return common.Address{}, err
	}
	if len(sig) != 65 {
		return common.Address{}, errors.New("invalid signature length")
	}

	// Wallets produce v = 27/28; go-ethereum's recovery expects v = 0/1.
	if sig[64] == 27 || sig[64] == 28 {
		sig[64] -= 27
	}

	hash := accounts.TextHash([]byte(message))
	pubKey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return common.Address{}, err
	}

	return crypto.PubkeyToAddress(*pubKey), nil
}
