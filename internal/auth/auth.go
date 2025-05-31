package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrPasswordIncorrect = errors.New("password incorrect")
	ErrJWTExpired        = errors.New("JWT is expired")
	ErrJWTInvalid        = errors.New("JWT is invalid")
)

type PasswordCredentials struct {
	Email    string
	Password string
}

type Authenticator[U any] interface {
	Authenticate(context.Context, PasswordCredentials) (*U, error)
}

type JWT interface {
	Sign(uid uuid.UUID) (string, time.Time, error)
	Parse(token string) (AccessClaims, error)
}

type AccessClaims struct {
	Issuer    string
	Subject   string
	Audience  []string
	ExpiresAt time.Time
	NotBefore time.Time
	IssuedAt  time.Time
	ID        string

	UserID uuid.UUID
}

type Token struct {
	jwt          JWT
	refreshStore RefreshStore
	// logger       logger.Logger
}

func NewToken(jwt JWT, refreshStore RefreshStore) *Token {
	return &Token{
		jwt:          jwt,
		refreshStore: refreshStore,
		// logger:       logger.With("component", "auth"),
	}
}

func (a *Token) SignJWT(uid uuid.UUID) (string, time.Time, error) {
	s, exp, err := a.jwt.Sign(uid)
	if err != nil {
		return "", time.Now(), fmt.Errorf("jwt new access: %w", err)
	}
	return s, exp, nil
}

func (a *Token) ParseJWT(token string) (AccessClaims, error) {
	claims, err := a.jwt.Parse(token)
	if err != nil {
		return AccessClaims{}, fmt.Errorf("jwt parse access: %w", err)
	}
	return claims, nil
}

func (a *Token) CreateRefresh(ctx context.Context, uid uuid.UUID) (string, error) {
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

func (a *Token) ExchangeRefresh(ctx context.Context, token string) (string, error) {
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

	newToken, err := a.CreateRefresh(ctx, sess.UserID)
	if err != nil {
		return "", fmt.Errorf("create refresh token: %w", err)
	}

	return newToken, nil
}
