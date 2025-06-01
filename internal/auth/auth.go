package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrPasswordIncorrect        = errors.New("password incorrect")
	ErrJWTExpired               = errors.New("JWT is expired")
	ErrJWTInvalid               = errors.New("JWT is invalid")
	ErrRefreshTokenExpired      = errors.New("refresh token is expired")
	ErrRefreshTokenUserMismatch = errors.New("token user does not match")
	ErrRefreshTokenUsed         = errors.New("refresh token is used")
)

type AccessToken struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}

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

func (a *Token) GetRefresh(ctx context.Context, id uuid.UUID) (RefreshSession, error) {
	r, err := a.refreshStore.Get(ctx, id)
	if err != nil {
		return RefreshSession{}, fmt.Errorf("get refresh: %w", err)
	}
	return r, nil
}

func (a *Token) CreateRefresh(ctx context.Context, uid uuid.UUID) (RefreshSession, error) {
	r, err := a.refreshStore.Insert(ctx, RefreshSession{
		ID:        uuid.New(),
		UserID:    uid,
		Used:      false,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
	})

	if err != nil {
		return RefreshSession{}, fmt.Errorf("create refresh session: %w", err)
	}

	return r, nil
}

func (a *Token) ExchangeRefresh(ctx context.Context, currRefresh RefreshSession) (RefreshSession, error) {
	if currRefresh.ExpiresAt.Before(time.Now()) {
		return RefreshSession{}, ErrRefreshTokenExpired
	}

	if currRefresh.Used {
		return RefreshSession{}, ErrRefreshTokenUsed
	}

	currRefresh.Used = true
	if err := a.refreshStore.Update(ctx, currRefresh); err != nil {
		return RefreshSession{}, fmt.Errorf("update refresh session: %w", err)
	}

	newRefresh, err := a.CreateRefresh(ctx, currRefresh.UserID)
	if err != nil {
		return RefreshSession{}, fmt.Errorf("create refresh token: %w", err)
	}

	return newRefresh, nil
}
