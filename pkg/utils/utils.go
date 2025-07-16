package utils

import (
	"math/rand"
)

func GenerateAddress() string {
	const charset = "abcdef0123456789"
	b := make([]byte, 64)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
