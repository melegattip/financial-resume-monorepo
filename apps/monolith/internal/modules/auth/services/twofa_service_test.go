package services

import (
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestTwoFAService() *twoFAService {
	return NewTwoFAService("TestApp").(*twoFAService)
}

func TestGenerateSecret(t *testing.T) {
	svc := newTestTwoFAService()

	secret, qrCode, backupCodes, err := svc.GenerateSecret("user@example.com")
	require.NoError(t, err)
	assert.NotEmpty(t, secret)
	assert.NotEmpty(t, qrCode) // base64 QR code
	assert.Len(t, backupCodes, 8)
}

func TestValidateCode_ValidCode(t *testing.T) {
	svc := newTestTwoFAService()

	secret, _, _, err := svc.GenerateSecret("user@example.com")
	require.NoError(t, err)

	// Generate a valid TOTP code using the library directly
	code, err := totp.GenerateCode(secret, time.Now())
	require.NoError(t, err)

	valid := svc.ValidateCode(secret, code)
	assert.True(t, valid)
}

func TestValidateCode_InvalidCode(t *testing.T) {
	svc := newTestTwoFAService()

	secret, _, _, err := svc.GenerateSecret("user@example.com")
	require.NoError(t, err)

	valid := svc.ValidateCode(secret, "invalid")
	assert.False(t, valid)
}

func TestValidateCode_CleanInput(t *testing.T) {
	svc := newTestTwoFAService()

	secret, _, _, err := svc.GenerateSecret("user@example.com")
	require.NoError(t, err)

	// Should handle spaces and dashes gracefully (even if code is wrong)
	valid := svc.ValidateCode(secret, "123 456")
	assert.False(t, valid) // wrong code, but no panic

	valid = svc.ValidateCode(secret, "123-456")
	assert.False(t, valid)
}

func TestGenerateBackupCodes(t *testing.T) {
	svc := newTestTwoFAService()

	codes, err := svc.GenerateBackupCodes(8)
	require.NoError(t, err)
	assert.Len(t, codes, 8)

	for _, code := range codes {
		// Backup codes should be in XXXX-XXXX format
		assert.Len(t, code, 9) // 4 + 1 (dash) + 4
		assert.Contains(t, code, "-")
	}

	// All codes should be unique
	unique := make(map[string]bool)
	for _, code := range codes {
		unique[code] = true
	}
	assert.Len(t, unique, 8)
}

func TestValidateBackupCode_ValidCode(t *testing.T) {
	svc := newTestTwoFAService()

	codes, err := svc.GenerateBackupCodes(8)
	require.NoError(t, err)

	remaining, valid := svc.ValidateBackupCode(codes, codes[0])
	assert.True(t, valid)
	assert.Len(t, remaining, 7)
	assert.NotContains(t, remaining, codes[0])
}

func TestValidateBackupCode_InvalidCode(t *testing.T) {
	svc := newTestTwoFAService()

	codes, err := svc.GenerateBackupCodes(8)
	require.NoError(t, err)

	remaining, valid := svc.ValidateBackupCode(codes, "INVALID-CODE")
	assert.False(t, valid)
	assert.Len(t, remaining, 8)
}

func TestValidateBackupCode_CaseInsensitive(t *testing.T) {
	svc := newTestTwoFAService()

	codes := []string{"ABCD-EFGH", "IJKL-MNOP"}

	remaining, valid := svc.ValidateBackupCode(codes, "abcd-efgh")
	assert.True(t, valid)
	assert.Len(t, remaining, 1)
}

func TestGenerateQRCode(t *testing.T) {
	svc := newTestTwoFAService()

	secret, _, _, err := svc.GenerateSecret("user@example.com")
	require.NoError(t, err)

	pngData, err := svc.GenerateQRCode(secret, "user@example.com")
	require.NoError(t, err)
	assert.NotEmpty(t, pngData)

	// PNG files start with the PNG magic bytes
	assert.Equal(t, byte(0x89), pngData[0])
	assert.Equal(t, byte('P'), pngData[1])
	assert.Equal(t, byte('N'), pngData[2])
	assert.Equal(t, byte('G'), pngData[3])
}
