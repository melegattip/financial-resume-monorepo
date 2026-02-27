package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/domain"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/ports"
)

type jwtService struct {
	secretKey     string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	issuer        string
}

// NewJWTService creates a new JWT service.
func NewJWTService(secretKey string, accessExpiry, refreshExpiry time.Duration, issuer string) ports.JWTService {
	return &jwtService{
		secretKey:     secretKey,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
		issuer:        issuer,
	}
}

func (j *jwtService) GenerateTokens(userID string, email string, tenantID string, role string) (*domain.TokenPair, error) {
	accessToken, accessExpiry, err := j.generateToken(userID, email, tenantID, role, "access", j.accessExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, _, err := j.generateToken(userID, email, tenantID, role, "refresh", j.refreshExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiry,
		TokenType:    "Bearer",
	}, nil
}

func (j *jwtService) ValidateAccessToken(tokenString string) (*domain.Claims, error) {
	return j.validateToken(tokenString, "access")
}

func (j *jwtService) ValidateRefreshToken(tokenString string) (*domain.Claims, error) {
	return j.validateToken(tokenString, "refresh")
}

func (j *jwtService) GenerateEmailVerificationToken(userID string, email string) (string, error) {
	token, _, err := j.generateToken(userID, email, "", "", "email_verification", 24*time.Hour)
	return token, err
}

func (j *jwtService) GeneratePasswordResetToken(userID string, email string) (string, error) {
	token, _, err := j.generateToken(userID, email, "", "", "password_reset", 1*time.Hour)
	return token, err
}

func (j *jwtService) ValidateEmailVerificationToken(tokenString string) (*domain.Claims, error) {
	return j.validateToken(tokenString, "email_verification")
}

func (j *jwtService) ValidatePasswordResetToken(tokenString string) (*domain.Claims, error) {
	return j.validateToken(tokenString, "password_reset")
}

func (j *jwtService) generateToken(userID string, email, tenantID, role, tokenType string, expiry time.Duration) (string, time.Time, error) {
	expiresAt := time.Now().Add(expiry)

	claims := domain.Claims{
		UserID:    userID,
		Email:     email,
		TenantID:  tenantID,
		Role:      role,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    j.issuer,
			Subject:   userID,
			Audience:  []string{"financial-resume"},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresAt, nil
}

func (j *jwtService) validateToken(tokenString, expectedType string) (*domain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*domain.Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	if claims.TokenType != expectedType {
		return nil, fmt.Errorf("invalid token type: expected %s, got %s", expectedType, claims.TokenType)
	}

	return claims, nil
}
