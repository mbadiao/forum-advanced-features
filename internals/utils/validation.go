package utils

import (
	"regexp"
	"unicode/utf8"
)

func IsValidEmail(email string) bool {
	// Regular expression for validating email addresses
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(regex, email)
	return match
}

func IsValidPassword(password string) bool {
	runeCount := utf8.RuneCountInString(password)
	return runeCount >= 5
}
