package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestPasswordService() *passwordService {
	return NewPasswordService(8).(*passwordService)
}

func TestHashPassword_Success(t *testing.T) {
	svc := newTestPasswordService()

	hash, err := svc.HashPassword("ValidP@ss1")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, "ValidP@ss1", hash)
}

func TestVerifyPassword_Success(t *testing.T) {
	svc := newTestPasswordService()

	hash, err := svc.HashPassword("ValidP@ss1")
	require.NoError(t, err)

	err = svc.VerifyPassword(hash, "ValidP@ss1")
	assert.NoError(t, err)
}

func TestVerifyPassword_Wrong(t *testing.T) {
	svc := newTestPasswordService()

	hash, err := svc.HashPassword("ValidP@ss1")
	require.NoError(t, err)

	err = svc.VerifyPassword(hash, "WrongP@ss1")
	assert.Error(t, err)
}

func TestValidatePasswordStrength_TooShort(t *testing.T) {
	svc := newTestPasswordService()

	err := svc.ValidatePasswordStrength("Sh@rt1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least 8 characters")
}

func TestValidatePasswordStrength_TooLong(t *testing.T) {
	svc := newTestPasswordService()

	longPass := make([]byte, 129)
	for i := range longPass {
		longPass[i] = 'A'
	}
	longPass[0] = '@'

	err := svc.ValidatePasswordStrength(string(longPass))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "less than 128")
}

func TestValidatePasswordStrength_NoUppercase(t *testing.T) {
	svc := newTestPasswordService()

	err := svc.ValidatePasswordStrength("lowercase@1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "uppercase")
}

func TestValidatePasswordStrength_NoSpecialChar(t *testing.T) {
	svc := newTestPasswordService()

	err := svc.ValidatePasswordStrength("NoSpecial1A")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "special character")
}

func TestValidatePasswordStrength_SequentialChars(t *testing.T) {
	svc := newTestPasswordService()

	err := svc.ValidatePasswordStrength("Abc@pass1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sequential")
}

func TestValidatePasswordStrength_SequentialNumbers(t *testing.T) {
	svc := newTestPasswordService()

	err := svc.ValidatePasswordStrength("Pass@123word")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sequential")
}

func TestValidatePasswordStrength_RepeatedChars(t *testing.T) {
	svc := newTestPasswordService()

	err := svc.ValidatePasswordStrength("Paaa@sword1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "consecutive identical")
}

func TestValidatePasswordStrength_Valid(t *testing.T) {
	svc := newTestPasswordService()

	validPasswords := []string{
		"ValidP@ss1",
		"MyS3cure!Pass",
		"G00dP@ssw0rd",
		"Str0ng#Key!",
	}

	for _, pw := range validPasswords {
		err := svc.ValidatePasswordStrength(pw)
		assert.NoError(t, err, "password %q should be valid", pw)
	}
}

func TestGenerateRandomPassword(t *testing.T) {
	pw := GenerateRandomPassword(16)
	assert.Len(t, pw, 16)
}

func TestGenerateRandomPassword_MinLength(t *testing.T) {
	pw := GenerateRandomPassword(4)
	assert.Len(t, pw, 8) // minimum is 8
}
