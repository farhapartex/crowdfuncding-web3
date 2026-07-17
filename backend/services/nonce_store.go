package services

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"sync"
	"time"
)

const nonceTTL = 5 * time.Minute

type nonceEntry struct {
	value     string
	expiresAt time.Time
}

type NonceStore struct {
	mu     sync.Mutex
	nonces map[string]nonceEntry
}

func NewNonceStore() *NonceStore {
	return &NonceStore{nonces: make(map[string]nonceEntry)}
}

func (s *NonceStore) Issue(address string) (string, error) {
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

func (s *NonceStore) PeekAndDelete(address string) (string, bool) {
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
