package helper

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
)

func GenerateRandomToken(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("token length must be greater than zero")
	}

	buffer := make([]byte, length)
	if _, err := rand.Read(buffer); err != nil {
		return "", fmt.Errorf("generate random token: %w", err)
	}

	return hex.EncodeToString(buffer), nil
}

func SHA256Hex(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func BasicAuthPassword(password string) string {
	sum := sha512.Sum512([]byte(password))
	hash := base64.StdEncoding.EncodeToString(sum[:])
	return fmt.Sprint("{SHA512}", hash)
}

func ID() (uuid.UUID, error) {
	return uuid.NewV7()
}

func MustID() uuid.UUID {
	id, _ := ID()
	return id
}
