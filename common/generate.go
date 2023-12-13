package common

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

func GenerateSHA256(salt, word string) string {
	payload := fmt.Sprint(salt, word)
	h := sha256.New()
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

func GenerateRandomString(chars string, strlen int) string {
	rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func GenerateOTP() string {
	return GenerateRandomString("1234567890", 6)
}
