package common

import (
	"math/rand"
	"time"
)

func GenerateRandom(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.NewSource(time.Now().UnixNano())
	s := make([]byte, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
