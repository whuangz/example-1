package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/scrypt"
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	shash, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		return "", err
	}
	hashed := fmt.Sprintf("%s.%s", hex.EncodeToString(shash), hex.EncodeToString(salt))
	return hashed, nil
}

func ValidateHash(hashPassword string, suppliedPassword string) (bool, error) {
	pwSalt := strings.Split(hashPassword, ".")
	salt, err := hex.DecodeString(pwSalt[1])
	if err != nil {
		return false, fmt.Errorf("Unable to verify user password")
	}

	shash, err := scrypt.Key([]byte(suppliedPassword), salt, 32768, 8, 1, 32)
	if err != nil {
		return false, fmt.Errorf("Unable to verify user password")
	}
	return hex.EncodeToString(shash) == pwSalt[0], nil

}
