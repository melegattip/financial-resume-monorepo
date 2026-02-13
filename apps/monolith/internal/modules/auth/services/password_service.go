package services

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"unicode"

	"golang.org/x/crypto/bcrypt"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/ports"
)

const bcryptCost = 12

type passwordService struct {
	minLength int
}

// NewPasswordService creates a new password service.
func NewPasswordService(minLength int) ports.PasswordService {
	return &passwordService{minLength: minLength}
}

func (p *passwordService) HashPassword(password string) (string, error) {
	if err := p.ValidatePasswordStrength(password); err != nil {
		return "", fmt.Errorf("password validation failed: %w", err)
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedBytes), nil
}

func (p *passwordService) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (p *passwordService) ValidatePasswordStrength(password string) error {
	if len(password) < p.minLength {
		return fmt.Errorf("password must be at least %d characters long", p.minLength)
	}

	if len(password) > 128 {
		return fmt.Errorf("password must be less than 128 characters long")
	}

	var (
		hasUpper   bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	if err := checkCommonPatterns(password); err != nil {
		return err
	}

	return nil
}

var sequentialPattern = regexp.MustCompile(`(?i)(abc|bcd|cde|def|efg|fgh|ghi|hij|ijk|jkl|klm|lmn|mno|nop|opq|pqr|qrs|rst|stu|tuv|uvw|vwx|wxy|xyz|123|234|345|456|567|678|789|890)`)

func checkCommonPatterns(password string) error {
	if sequentialPattern.MatchString(password) {
		return fmt.Errorf("password cannot contain sequential characters")
	}

	if hasRepeatedChars(password, 3) {
		return fmt.Errorf("password cannot contain more than 2 consecutive identical characters")
	}

	return nil
}

// hasRepeatedChars checks if the string has n or more consecutive identical characters.
func hasRepeatedChars(s string, n int) bool {
	if len(s) < n {
		return false
	}
	count := 1
	for i := 1; i < len(s); i++ {
		if s[i] == s[i-1] {
			count++
			if count >= n {
				return true
			}
		} else {
			count = 1
		}
	}
	return false
}

// secureRandomInt returns a cryptographically secure random int in [0, max).
func secureRandomInt(max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic(fmt.Sprintf("crypto/rand failed: %v", err))
	}
	return int(n.Int64())
}

// GenerateRandomPassword creates a random password meeting all strength requirements.
func GenerateRandomPassword(length int) string {
	if length < 8 {
		length = 8
	}

	const (
		lower   = "abcdefghijkmnopqrstuvwxyz"
		upper   = "ABCDEFGHJKLMNPQRSTUVWXYZ"
		digits  = "23456789"
		special = "!@#$%^&*"
	)
	charset := lower + upper + digits + special

	password := make([]byte, length)
	password[0] = lower[secureRandomInt(len(lower))]
	password[1] = upper[secureRandomInt(len(upper))]
	password[2] = digits[secureRandomInt(len(digits))]
	password[3] = special[secureRandomInt(len(special))]

	for i := 4; i < length; i++ {
		password[i] = charset[secureRandomInt(len(charset))]
	}

	for i := len(password) - 1; i > 0; i-- {
		j := secureRandomInt(i + 1)
		password[i], password[j] = password[j], password[i]
	}

	return string(password)
}
