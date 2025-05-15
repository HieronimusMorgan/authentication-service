package utils

import (
	"authentication/internal/dto/out"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ValidationTrimSpace(s string) string {
	trim := strings.TrimSpace(s)
	trim = strings.Join(strings.Fields(trim), " ") // Remove extra spaces
	return trim
}

// Custom error messages
var (
	ErrUsernameLength  = errors.New("username must be between 3 and 20 characters")
	ErrUsernameInvalid = errors.New("username can only contain alphanumeric characters and underscores")
)

// ValidateUsername checks if the username meets the criteria
func ValidateUsername(username string) error {
	username = strings.TrimSpace(username) // Trim spaces

	if len(username) < 3 || len(username) > 20 {
		return ErrUsernameLength
	}

	validUsername := regexp.MustCompile(`^[a-zA-Z0-9@#$%&_\-.]+$`)
	if !validUsername.MatchString(username) {
		return ErrUsernameInvalid
	}

	return nil
}

func DecryptOptionalString(value *string, encryption Encryption) *string {
	if value == nil {
		return nil
	}
	decrypted, err := encryption.Decrypt(*value)
	if err != nil {
		return value
	}
	return &decrypted
}

func ConvertToUint(input string) (uint, error) {
	parsed, err := strconv.ParseUint(input, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid uint value: %w", err)
	}
	return uint(parsed), nil
}

func ValidateEmail(email string) error {
	validEmail := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !validEmail.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func ContainsRole(roles []out.RoleResponse, id uint) bool {
	for _, r := range roles {
		if r.RoleID == id {
			return true
		}
	}
	return false
}

func ContainsResource(resources []out.ResourceResponse, id uint) bool {
	for _, r := range resources {
		if r.ResourceID == id {
			return true
		}
	}
	return false
}

func DerefStr(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func DerefInt(i *int) int {
	if i != nil {
		return *i
	}
	return 0
}
