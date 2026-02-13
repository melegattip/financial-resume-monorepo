package services

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"image/png"
	"strings"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/ports"
)

type twoFAService struct {
	issuer string
}

// NewTwoFAService creates a new 2FA service.
func NewTwoFAService(issuer string) ports.TwoFAService {
	return &twoFAService{issuer: issuer}
}

func (t *twoFAService) GenerateSecret(userEmail string) (string, string, []string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      t.issuer,
		AccountName: userEmail,
		SecretSize:  32,
	})
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to generate TOTP key: %w", err)
	}

	backupCodes, err := t.GenerateBackupCodes(8)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to generate backup codes: %w", err)
	}

	// Generate QR code as base64 PNG
	qrBase64, err := t.generateQRBase64(key)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	return key.Secret(), qrBase64, backupCodes, nil
}

func (t *twoFAService) GenerateQRCode(secret, userEmail string) ([]byte, error) {
	key, err := otp.NewKeyFromURL(fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s",
		t.issuer, userEmail, secret, t.issuer))
	if err != nil {
		return nil, fmt.Errorf("failed to create key from URL: %w", err)
	}

	img, err := key.Image(256, 256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code image: %w", err)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("failed to encode QR code as PNG: %w", err)
	}

	return buf.Bytes(), nil
}

func (t *twoFAService) ValidateCode(secret, code string) bool {
	code = strings.ReplaceAll(code, " ", "")
	code = strings.ReplaceAll(code, "-", "")
	return totp.Validate(code, secret)
}

func (t *twoFAService) GenerateBackupCodes(count int) ([]string, error) {
	codes := make([]string, count)
	for i := 0; i < count; i++ {
		code, err := generateRandomCode(8)
		if err != nil {
			return nil, fmt.Errorf("failed to generate backup code %d: %w", i+1, err)
		}
		codes[i] = code
	}
	return codes, nil
}

func (t *twoFAService) ValidateBackupCode(codes []string, providedCode string) ([]string, bool) {
	providedCode = strings.ToUpper(strings.TrimSpace(providedCode))

	for i, code := range codes {
		if strings.ToUpper(code) == providedCode {
			remaining := make([]string, 0, len(codes)-1)
			remaining = append(remaining, codes[:i]...)
			remaining = append(remaining, codes[i+1:]...)
			return remaining, true
		}
	}

	return codes, false
}

func (t *twoFAService) generateQRBase64(key *otp.Key) (string, error) {
	img, err := key.Image(256, 256)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR image: %w", err)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", fmt.Errorf("failed to encode QR PNG: %w", err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func generateRandomCode(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)

	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to read random bytes: %w", err)
	}

	for i, v := range b {
		b[i] = charset[v%byte(len(charset))]
	}

	if length == 8 {
		return fmt.Sprintf("%s-%s", string(b[:4]), string(b[4:])), nil
	}

	return string(b), nil
}
