package server

import "crypto/rand"

func GenerateRandomKey() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	return string(key), err
}
