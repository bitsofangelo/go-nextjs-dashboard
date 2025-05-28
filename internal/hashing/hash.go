package hashing

import (
	"fmt"
)

type Hasher interface {
	Hash(password string) (string, error)
	Check(password, hash string) (bool, error)
}

type Hash struct {
	hasher Hasher
}

func New(hasher Hasher) *Hash {
	return &Hash{
		hasher: hasher,
	}
}

// Make hashes a plaintext password using bcrypt
func (a *Hash) Make(password string) (string, error) {
	s, err := a.hasher.Hash(password)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return s, nil
}

// Check checks if the given password matches the hashed password
func (a *Hash) Check(password, hash string) (bool, error) {
	match, err := a.hasher.Check(password, hash)
	if err != nil {
		return false, fmt.Errorf("check password hash: %w", err)
	}
	return match, nil
}
