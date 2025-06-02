package hashing

import (
	"fmt"

	"github.com/gelozr/go-dash/internal/config"
)

type Hasher interface {
	Hash(password string) (string, error)
	Check(password, hash string) (bool, error)
}

type Manager struct {
	bcrypt        BcryptHasher
	argon2id      Argon2IDHasher
	defaultHasher string
}

func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		bcrypt:        NewBcryptHasher(),
		argon2id:      NewArgon2IDHasher(),
		defaultHasher: getDefaultDriver(cfg),
	}
}

func (m *Manager) getHasher(hasher string) Hasher {
	switch hasher {
	case "argon2id":
		return m.argon2id
	default:
		return m.bcrypt
	}
}

// Hash hashes a plaintext using bcrypt
func (m *Manager) Hash(text string) (string, error) {
	s, err := m.getHasher(m.defaultHasher).Hash(text)
	if err != nil {
		return "", fmt.Errorf("hash: %w", err)
	}
	return s, nil
}

// Check checks if the given text matches the hashed text
func (m *Manager) Check(text, hash string) (bool, error) {
	match, err := m.getHasher(m.defaultHasher).Check(text, hash)
	if err != nil {
		return false, fmt.Errorf("check hash: %w", err)
	}
	return match, nil
}

func getDefaultDriver(cfg *config.Config) string {
	defaultHasher := "bcrypt"
	if cfg.HashingDriver != "" {
		defaultHasher = cfg.HashingDriver
	}
	return defaultHasher
}
