package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"testing"
	"time"
)

func TestRefreshTokensHaveIndependentRandomIDs(t *testing.T) {
	user := primitive.NewObjectID()
	secret := strings.Repeat("k", MinJWTSecretBytes)
	seen := map[string]bool{}
	for i := 0; i < 100; i++ {
		pair, err := GenerateTokenPair(user, "test@example.invalid", secret, 0, time.Minute, time.Hour)
		if err != nil {
			t.Fatal(err)
		}
		claims, err := ValidateToken(pair.RefreshToken, secret)
		if err != nil {
			t.Fatal(err)
		}
		if len(claims.ID) < 32 || seen[claims.ID] {
			t.Fatal("missing/duplicate random token ID")
		}
		seen[claims.ID] = true
	}
}
func TestLegacyRefreshWithoutIDStillValid(t *testing.T) {
	secret := strings.Repeat("k", MinJWTSecretBytes)
	claims := Claims{UserID: primitive.NewObjectID().Hex(), TokenType: RefreshToken, RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour))}}
	raw, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := ValidateToken(raw, secret)
	if err != nil || parsed.ID != "" {
		t.Fatal("legacy validation changed")
	}
}
