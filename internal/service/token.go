package service

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateToken(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return hex.EncodeToString(b)
}
