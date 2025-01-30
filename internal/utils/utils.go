package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	result := make([]byte, length)

	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}

	return string(result)
}
