package apikey

import (
	"crypto/sha256"
	"fmt"

	"github.com/google/uuid"
)

// Generates a unique API key and a matching hash.
func GenerateAPIKey() (raw string, hash [32]byte) {
	// Note: Not suitable for production! Demonstration purposes only.
	// Use a proper API key generation tool to provide secure, protected API keys.
	raw = uuid.NewString()
	hash = HashAPIKey(raw)
	return
}

// Hashes a raw API key for comparison or storage.
func HashAPIKey(raw string) [32]byte {
	return sha256.Sum256([]byte(raw))
}

// Converts hash into string for storage.
func HashByteToString(hash [32]byte) string {
	return fmt.Sprintf("%x", hash)
}
