package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateRandomString(n int) (string, error) {

	if n <= 0 {
		return "", fmt.Errorf("invalid length: %d", n)
	}

	bytes := make([]byte, n)

	_, err := rand.Read(bytes)

	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return hex.EncodeToString(bytes), nil

}
