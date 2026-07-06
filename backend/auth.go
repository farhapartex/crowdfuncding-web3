package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const nonceTTL = 5 * time.Minute
const sessionTTL = time.Hour

type nonceEntry struct {
	value     string
	expiresAt time.Time
}

type nonceStore struct {
	mu     sync.Mutex
	nonces map[string]nonceEntry
}

func newNonceStore() *nonceStore {
	return &nonceStore{nonces: make(map[string]nonceEntry)}
}

func (s *nonceStore) Issue(address string) (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	nonce := hex.EncodeToString(buf)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.nonces[strings.ToLower(address)] = nonceEntry{value: nonce, expiresAt: time.Now().Add(nonceTTL)}

	return nonce, nil
}

func (s *nonceStore) PeekAndDelete(address string) (string, bool) {
	key := strings.ToLower(address)

	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.nonces[key]
	delete(s.nonces, key)
	if !ok || time.Now().After(entry.expiresAt) {
		return "", false
	}

	return entry.value, true
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

func generateJWT(secret []byte, address string) (string, error) {
	claims := jwt.MapClaims{
		"address": address,
		"exp":     time.Now().Add(sessionTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func parseJWT(secret []byte, tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	address, ok := claims["address"].(string)
	if !ok {
		return "", errors.New("token missing address claim")
	}

	return address, nil
}

func authMiddleware(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		tokenString := strings.TrimPrefix(header, "Bearer ")
		if tokenString == "" || tokenString == header {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing bearer token"})
			return
		}

		address, err := parseJWT(secret, tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set("address", address)
		c.Next()
	}
}
