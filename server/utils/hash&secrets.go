package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateApiKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return "nxt_live_" + hex.EncodeToString(bytes), nil
}

func GenerateWebhookSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "nxt_whs_" + hex.EncodeToString(bytes), nil
}

func HashSecrets(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func VerifySecrets(input string, hash string) bool {
	return hash == HashSecrets(input)
}
