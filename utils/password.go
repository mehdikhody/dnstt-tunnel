package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GeneratePassword(length int) (string, error) {
	b := make([]byte, length)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	encoder := base64.StdEncoding.WithPadding(base64.NoPadding)
	password := encoder.EncodeToString(b)
	return password, nil
}
