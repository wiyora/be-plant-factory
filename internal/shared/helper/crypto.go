package helper

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
)

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
