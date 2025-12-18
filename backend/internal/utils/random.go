package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateShareID generates a URL-safe random string of the specified length
func GenerateShareID(length int) (string, error) {
	// Calculate how many bytes we need (base64 encoding is ~4/3 the size)
	byteLength := (length * 3) / 4
	if byteLength < 1 {
		byteLength = 1
	}

	bytes := make([]byte, byteLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Use URL-safe base64 encoding and trim to desired length
	encoded := base64.URLEncoding.EncodeToString(bytes)
	if len(encoded) > length {
		encoded = encoded[:length]
	}

	return encoded, nil
}
