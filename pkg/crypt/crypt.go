package crypt

import "golang.org/x/crypto/bcrypt"

// GeneratePassword set Password in crypt
func GeneratePassword(password []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	return string(hash), err
}

func ComparePassword(password []byte, dbPassword []byte) error {
	return bcrypt.CompareHashAndPassword(password, dbPassword)
}
