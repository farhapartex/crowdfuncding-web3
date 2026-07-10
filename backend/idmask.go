package main

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"math/rand"

	"github.com/sqids/sqids-go"
)

const baseIDMaskAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func newIDMasker(secret string) (*sqids.Sqids, error) {
	return sqids.New(sqids.Options{
		Alphabet:  shuffleAlphabet(baseIDMaskAlphabet, secret),
		MinLength: 8,
	})
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

func maskID(masker *sqids.Sqids, id uint64) (string, error) {
	return masker.Encode([]uint64{id})
}

func unmaskID(masker *sqids.Sqids, masked string) (uint64, error) {
	numbers := masker.Decode(masked)
	if len(numbers) != 1 {
		return 0, errors.New("invalid id")
	}
	return numbers[0], nil
}
