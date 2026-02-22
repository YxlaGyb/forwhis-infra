package utils

import (
	"regexp"
	"unicode"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func IsValidEmail(email string) bool {
	if len(email) == 0 || len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}

func IsValidPassword(password string) bool {
	if len(password) < 8 || len(password) > 72 {
		return false
	}
	hasUpper := false
	hasLower := false
	hasDigit := false
	for _, r := range password {
		if unicode.IsUpper(r) {
			hasUpper = true
		} else if unicode.IsLower(r) {
			hasLower = true
		} else if unicode.IsDigit(r) {
			hasDigit = true
		}
	}
	return hasUpper && hasLower && hasDigit
}

func IsValidName(name string) bool {
	if len(name) < 2 || len(name) > 50 {
		return false
	}
	return true
}
