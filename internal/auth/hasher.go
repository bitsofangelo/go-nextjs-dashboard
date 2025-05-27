package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

type ArgonHasher struct{}

func NewArgonHasher() *ArgonHasher {
	return &ArgonHasher{}
}

func (a *ArgonHasher) HashPassword(password string) (string, error) {
	// bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	s, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("argon hash password: %w", err)
	}
	return s, nil
}

func (a *ArgonHasher) CheckPasswordHash(password, hash string) (bool, error) {
	// err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	ok, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("argon compare password hash: %w", err)
	}
	return ok, nil
}
