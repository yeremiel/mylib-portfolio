package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
)

// GenerateSecureToken generates a cryptographically secure random token:
// 32 bytes (256 bits) of random data, base64url-encoded.
func GenerateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// HashToken hashes a token with SHA-256.
// Store the hashed value; at verification time, hash the presented token and compare.
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// VerifyToken compares a raw token against its stored hash.
func VerifyToken(token, hashedToken string) bool {
	return HashToken(token) == hashedToken
}

// GenerateOTP generates a cryptographically secure 6-digit numeric OTP.
// Returns a string in the range "000000"-"999999".
func GenerateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
