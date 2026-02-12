package verification

import (
	"fmt"
	"math/rand"
	"time"
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateVerificationCode generates a 6-digit numeric code
func GenerateVerificationCode() string {
	code := seededRand.Intn(900000) + 100000
	return fmt.Sprintf("%06d", code)
}

// GenerateTempPassword generates a random alphanumeric password of given length
func GenerateTempPassword(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// GenerateRandomString generates a random string with given length using the default charset
func GenerateRandomString(length int) string {
	return GenerateTempPassword(length)
}
