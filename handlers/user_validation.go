package handlers

import (
	"fmt"
	"regexp"
	"strings"
)

func ValidatePassword(password string) error {
	if len(password) <= 5 {
		return fmt.Errorf("password cannot be less than 6 characters")
	}
	if len(password) > 16 {
		return fmt.Errorf("password cannot be more than 16 characters")
	}

	if strings.Contains(password, " ") {
		return fmt.Errorf("invalid password")
	}

	return nil
}

var userNameRegExp = regexp.MustCompile("[^A-Za-z0-9]")

func ValidateUsername(username string) error {
	if len(username) <= 2 {
		return fmt.Errorf("username cannot be less than 3 characters")
	}
	if len(username) > 16 {
		return fmt.Errorf("username cannot be more than 16 characters")
	}

	if userNameRegExp.MatchString(username) {
		return fmt.Errorf("username could contain only A-Z, a-z or 0-9 characters")
	}

	return nil
}
