package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/argon2"
)

const (
	timeCost   = 3
	memoryCost = 64 * 1024
	threads    = 4
	keyLen     = 32
	saltLen    = 16
)

var ErrRandomReadFailed = errors.New("failed to generate random salt")

func HashPassword(password string) (string, error) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", ErrRandomReadFailed
	}
	hash := argon2.IDKey([]byte(password), salt, timeCost, memoryCost, threads, keyLen)
	return base64.StdEncoding.EncodeToString(append(salt, hash...)), nil
}

func VerifyPassword(hashed, password string) bool {
	data, err := base64.StdEncoding.DecodeString(hashed)
	if err != nil || len(data) != saltLen+keyLen {
		return false
	}
	salt := data[:saltLen]
	stored := data[saltLen:]
	hash := argon2.IDKey([]byte(password), salt, timeCost, memoryCost, threads, keyLen)
	for i := range hash {
		if hash[i] != stored[i] {
			return false
		}
	}
	return true
}
