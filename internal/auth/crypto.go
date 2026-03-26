package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// NewSecret generates a strong secret string
func NewSecret(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// HashSecret creates a hash of the secret to avoid
// storing the secret in plaintext
func HashSecret(secret string) string {
	sum := sha256.Sum256([]byte(secret))

	return hex.EncodeToString(sum[:])
}
