package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"go-nextjs-dashboard/internal/logger"
)

var (
	ErrPasswordIncorrect = errors.New("password incorrect")
	ErrJWTExpired        = errors.New("JWT is expired")
	ErrJWTInvalid        = errors.New("JWT is invalid")
)

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) (bool, error)
}

type Claims struct {
	Issuer    string
	Subject   string
	Audience  []string
	ExpiresAt time.Time
	NotBefore time.Time
	IssuedAt  time.Time
	ID        string

	UserID uuid.UUID
}

type JWT interface {
	NewAccess(uid uuid.UUID) (string, time.Time, error)
	ParseAccess(token string) (Claims, error)
}

type Service struct {
	hasher       PasswordHasher
	jwt          JWT
	refreshStore RefreshStore
	logger       logger.Logger
}

func New(hasher PasswordHasher, jwt JWT, refreshStore RefreshStore, logger logger.Logger) *Service {
	return &Service{
		hasher:       hasher,
		jwt:          jwt,
		refreshStore: refreshStore,
		logger:       logger.With("component", "auth"),
	}
}

// HashPassword hashes a plaintext password using bcrypt
func (a *Service) HashPassword(password string) (string, error) {
	s, err := a.hasher.HashPassword(password)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return s, nil
}

// CheckPasswordHash checks if the given password matches the hashed password
func (a *Service) CheckPasswordHash(password, hash string) (bool, error) {
	match, err := a.hasher.CheckPasswordHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("check password hash: %w", err)
	}
	return match, nil
}

func (a *Service) NewJWT(uid uuid.UUID) (string, time.Time, error) {
	s, exp, err := a.jwt.NewAccess(uid)
	if err != nil {
		return "", time.Now(), fmt.Errorf("jwt new access: %w", err)
	}
	return s, exp, nil
}

func (a *Service) ParseJWT(token string) (Claims, error) {
	claims, err := a.jwt.ParseAccess(token)
	if err != nil {
		return Claims{}, fmt.Errorf("jwt parse access: %w", err)
	}
	return claims, nil
}

func (a *Service) CreateRefreshToken(ctx context.Context, uid uuid.UUID) (string, error) {
	r, err := a.refreshStore.Insert(ctx, RefreshSession{
		ID:        uuid.New(),
		UserID:    uid,
		Used:      false,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
	})

	if err != nil {
		return "", fmt.Errorf("create refresh session: %w", err)
	}

	return r.ID.String(), nil
}

func (a *Service) ExchangeRefreshToken(ctx context.Context, token string) (string, error) {
	id, err := uuid.Parse(token)
	if err != nil {
		return "", fmt.Errorf("token uuid parse: %w", err)
	}

	sess, err := a.refreshStore.Get(ctx, id)
	if err != nil {
		return "", fmt.Errorf("get refresh session: %w", err)
	}

	sess.Used = true
	if err = a.refreshStore.Update(ctx, sess); err != nil {
		return "", fmt.Errorf("update refresh session: %w", err)
	}

	newToken, err := a.CreateRefreshToken(ctx, sess.UserID)
	if err != nil {
		return "", fmt.Errorf("create refresh token: %w", err)
	}

	return newToken, nil
}
