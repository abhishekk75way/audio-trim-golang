package utils

import (
	"errors"
	"regexp"
	"strings"
)

func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

func ValidateEmail(email string) (string, error) {
	email = TrimSpace(email)

	if email == "" {
		return "", errors.New("email cannot be empty")
	}

	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !regex.MatchString(email) {
		return "", errors.New("invalid email format")
	}

	return email, nil
}

func ValidatePassword(password string) (string, error) {
	password = TrimSpace(password)

	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	if len(password) < 4 {
		return "", errors.New("password must be at least 4 characters")
	}

	return password, nil
}

func ValidatePasswordChange(oldPass, newPass string) (string, string, error) {
	oldPass = TrimSpace(oldPass)
	newPass = TrimSpace(newPass)

	if oldPass == "" {
		return "", "", errors.New("old password cannot be empty")
	}

	if len(newPass) < 4 {
		return "", "", errors.New("new password must be at least 4 characters")
	}

	return oldPass, newPass, nil
}
