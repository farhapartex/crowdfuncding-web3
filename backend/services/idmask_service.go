package services

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"math/rand"

	"github.com/sqids/sqids-go"
)

const baseIDMaskAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type IDMaskService struct {
	sqids *sqids.Sqids
}

func NewIDMaskService(secret string) (*IDMaskService, error) {
	s, err := sqids.New(sqids.Options{
		Alphabet:  shuffleAlphabet(baseIDMaskAlphabet, secret),
		MinLength: 8,
	})
	if err != nil {
		return nil, err
	}

	return &IDMaskService{sqids: s}, nil
}

func shuffleAlphabet(alphabet, secret string) string {
	hash := sha256.Sum256([]byte(secret))
	seed := int64(binary.BigEndian.Uint64(hash[:8]))
	r := rand.New(rand.NewSource(seed))

	chars := []rune(alphabet)
	r.Shuffle(len(chars), func(i, j int) {
		chars[i], chars[j] = chars[j], chars[i]
	})

	return string(chars)
}

func (s *IDMaskService) Mask(id uint64) (string, error) {
	return s.sqids.Encode([]uint64{id})
}

func (s *IDMaskService) Unmask(masked string) (uint64, error) {
	numbers := s.sqids.Decode(masked)
	if len(numbers) != 1 {
		return 0, errors.New("invalid id")
	}
	return numbers[0], nil
}
