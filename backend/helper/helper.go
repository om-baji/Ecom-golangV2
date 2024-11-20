package helper

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, bool) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", false
	}

	return string(hash), true
}

func VerifyPassword(hashed string, password string) bool {

	isValid := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))

	return isValid == nil
}
