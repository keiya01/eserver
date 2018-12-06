package auth

import (
	"log"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(plainPassword string) (string, error) {
	password := []byte(plainPassword)
	encryptedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)

	if err != nil {
		return "", errors.Wrap(err, "EncryptPassword()")
	}

	return string(encryptedPassword), nil
}

func ComparePassword(encryptedPassword string, plainPassword string) bool {
	byteEncryptedPassword := []byte(encryptedPassword)
	bytePlainPassword := []byte(plainPassword)

	err := bcrypt.CompareHashAndPassword(byteEncryptedPassword, bytePlainPassword)
	if err != nil {
		log.Printf("ComparePassword(): %v", err)
		return false
	}

	return true
}
