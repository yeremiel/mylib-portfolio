package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MinJWTSecretBytes is the minimum JWT secret length in bytes (RFC 7518 recommendation).
const MinJWTSecretBytes = 32

// ValidateJWTSecretEntropy checks that the JWT secret is at least 32 bytes long.
func ValidateJWTSecretEntropy(secret string) error {
	if len([]byte(secret)) < MinJWTSecretBytes {
		return fmt.Errorf("JWT secret must be at least %d bytes (got %d bytes)", MinJWTSecretBytes, len([]byte(secret)))
	}
	return nil
}

// Default expiry times (used when not set in config).
const (
	DefaultAccessTokenExpiry  = 15 * time.Minute   // 15 minutes
	DefaultRefreshTokenExpiry = 7 * 24 * time.Hour // 7 days
)

// TokenType identifies the kind of token.
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// Claims is the JWT payload.
type Claims struct {
	UserID       string    `json:"user_id"`
	Email        string    `json:"email"`
	TokenType    TokenType `json:"token_type"`
	TokenVersion int       `json:"token_version"`
	jwt.RegisteredClaims
}

// TokenPair holds an access token and a refresh token.
type TokenPair struct {
	AccessToken           string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	AccessTokenExpiresAt  int64  `json:"access_token_expires_at"`
	RefreshTokenExpiresAt int64  `json:"refresh_token_expires_at"`
}

// JWT-related errors.
var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("expired token")
	ErrInvalidTokenType = errors.New("invalid token type")
)

// GenerateTokenPair issues an access token and a refresh token.
// The refresh token carries a random token ID (jti).
func GenerateTokenPair(userID primitive.ObjectID, email, secretKey string, tokenVersion int, accessExpiry, refreshExpiry time.Duration) (*TokenPair, error) {
	// Apply defaults
	if accessExpiry == 0 {
		accessExpiry = DefaultAccessTokenExpiry
	}
	if refreshExpiry == 0 {
		refreshExpiry = DefaultRefreshTokenExpiry
	}

	tokenID, err := GenerateSecureToken()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	userIDStr := userID.Hex()

	// Build the access token
	accessTokenExpiresAt := now.Add(accessExpiry)
	accessClaims := Claims{
		UserID:       userIDStr,
		Email:        email,
		TokenType:    AccessToken,
		TokenVersion: tokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   userIDStr,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	// Build the refresh token
	refreshTokenExpiresAt := now.Add(refreshExpiry)
	refreshClaims := Claims{
		UserID:       userIDStr,
		Email:        email,
		TokenType:    RefreshToken,
		TokenVersion: tokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			ExpiresAt: jwt.NewNumericDate(refreshTokenExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   userIDStr,
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:           accessTokenString,
		RefreshToken:          refreshTokenString,
		AccessTokenExpiresAt:  accessTokenExpiresAt.Unix(),
		RefreshTokenExpiresAt: refreshTokenExpiresAt.Unix(),
	}, nil
}

// GenerateAccessToken issues only an access token (used when refreshing).
func GenerateAccessToken(userID primitive.ObjectID, email, secretKey string, tokenVersion int, expiry time.Duration) (string, int64, error) {
	if expiry == 0 {
		expiry = DefaultAccessTokenExpiry
	}

	now := time.Now()
	userIDStr := userID.Hex()
	expiresAt := now.Add(expiry)

	claims := Claims{
		UserID:       userIDStr,
		Email:        email,
		TokenType:    AccessToken,
		TokenVersion: tokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   userIDStr,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt.Unix(), nil
}

// ValidateToken verifies a token and returns its claims.
func ValidateToken(tokenString, secretKey string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing algorithm (HS256 only)
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		// Detect expiry
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GetRefreshTokenKey builds the Redis key under which a refresh token is stored.
// Format: "refresh_token:{user_id}"
func GetRefreshTokenKey(userID string) string {
	return "refresh_token:" + userID
}
