package utils

import (
	"crypto/rand"
	"math/big"
)

const (
	charsetAlphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charsetLowerNumeric = "abcdefghijklmnopqrstuvwxyz0123456789"
)

func SecureRandomString(length int) string {
	return secureRandomFromCharset(length, charsetAlphanumeric)
}

func SecureRandomLowerNumeric(length int) string {
	return secureRandomFromCharset(length, charsetLowerNumeric)
}

func secureRandomFromCharset(length int, charset string) string {
	b := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))
	for i := range b {
		n, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			b[i] = charset[0]
			continue
		}
		b[i] = charset[n.Int64()]
	}
	return string(b)
}
