package secretcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
)

// keyFromSecret converts the app secret into a fixed 32-byte key,
// then hash it with SHA256.
func keyFromSecret(appSecret string) []byte {
	sum := sha256.Sum256([]byte(appSecret))
	return sum[:]
}

func Encrypt(plaintext string, appSecret string) (string, error) {
	// Get 32 byte AES key from app secret
	key := keyFromSecret(appSecret)

	// Create AES block cipher using derived key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Wrap AES block in GCM mode for authenticated encryption.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// A 'nonce' stands for a 'number once.' It is a number
	// that can only be used once in the key's value, and
	// helps with ensuring uniqueness when re-encyryping
	// the same text with the same key.
	nonce := make([]byte, gcm.NonceSize())

	// Fill nonce with 'random' bytes.
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt input plaintext
	cipherText := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	// Prepend nonce to ciphertext so it can be used during decryption
	out := append(nonce, cipherText...)

	// Encode result as base64 sso it can be stored as text in the database
	return base64.StdEncoding.EncodeToString(out), nil
}
