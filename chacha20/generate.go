package chacha20

import (
	"math/rand"
)

func GenerateNonce() string {
	return generateRandomBytesString(12)
}

func GenerateKey() string {
	return generateRandomBytesString(32)
}

func generateRandomBytesString(length int) string {
	charset := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+<>/?")

	result := make([]rune, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
