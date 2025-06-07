package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gelozr/forge/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/gelozr/go-dash/internal/config"
	"github.com/gelozr/go-dash/internal/user"
)

type AccessToken struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
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

type JWTClaims struct {
	UserID uuid.UUID
	jwt.RegisteredClaims
}

type JWTDriver struct {
	hmacKey           []byte
	refreshSessionSvc *Token
}

func NewJWTDriver(cfg *config.Config, refreshSessionSvc *Token) *JWTDriver {
	return &JWTDriver{
		hmacKey:           []byte(cfg.JWTHmacKey),
		refreshSessionSvc: refreshSessionSvc,
	}
}

func (d *JWTDriver) Sign(uid uuid.UUID) (string, time.Time, error) {
	claims := JWTClaims{
		UserID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Issuer:    "myapp",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := tok.SignedString(d.hmacKey)
	if err != nil {
		return "", time.Now(), fmt.Errorf("signing token: %w", err)
	}

	exp, err := tok.Claims.GetExpirationTime()
	if err != nil {
		return "", time.Now(), fmt.Errorf("getting expiration time: %w", err)
	}

	return signed, exp.Time, nil
}

func (d *JWTDriver) Parse(tokenStr string) (AccessClaims, error) {
	tok, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (any, error) {
		return d.hmacKey, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenSignatureInvalid), errors.Is(err, jwt.ErrTokenMalformed):
			return AccessClaims{}, ErrJWTInvalid
		case errors.Is(err, jwt.ErrTokenExpired):
			return AccessClaims{}, ErrJWTExpired
		default:
			return AccessClaims{}, fmt.Errorf("parsing token: %w", err)
		}
	}

	ac := tok.Claims.(*JWTClaims)

	claims := AccessClaims{
		Issuer:    ac.Issuer,
		Subject:   ac.Subject,
		Audience:  ac.Audience,
		ExpiresAt: ac.ExpiresAt.Time,
		// NotBefore: ac.NotBefore.Time,
		IssuedAt: ac.IssuedAt.Time,
		ID:       ac.ID,
		UserID:   ac.UserID,
	}

	return claims, nil
}

func (d *JWTDriver) IssueToken(ctx context.Context, user auth.User) (any, error) {
	uid, ok := user.UserID().(uuid.UUID)
	if !ok {
		return nil, errors.New("invalid user id")
	}

	jwtStr, exp, err := d.Sign(uid)
	if err != nil {
		return nil, fmt.Errorf("sign jwt: %w", err)
	}

	refreshSess, err := d.refreshSessionSvc.CreateRefresh(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("create refresh session: %w", err)
	}

	accessToken := AccessToken{
		AccessToken:  jwtStr,
		RefreshToken: refreshSess.ID.String(),
		ExpiresIn:    int(time.Until(exp).Seconds()),
	}

	return accessToken, nil
}

func (d *JWTDriver) Login(ctx context.Context, user auth.User) (any, error) {
	return d.IssueToken(ctx, user)
}

func (d *JWTDriver) Validate(_ context.Context, payload any) (auth.Verified, error) {
	token, ok := payload.(string)
	if !ok {
		return auth.Verified{}, errors.New("invalid jwt payload")
	}

	claims, err := d.Parse(token)
	if err != nil {
		return auth.Verified{}, fmt.Errorf("parse token: %w", err)
	}

	return auth.Verified{
		User: user.User{ID: claims.UserID},
	}, nil
}

func (d *JWTDriver) RefreshToken(ctx context.Context, refreshToken string) (any, error) {
	refreshTokenID, err := uuid.Parse(refreshToken)
	if err != nil {
		return nil, ErrRefreshTokenInvalid
	}

	currRefresh, err := d.refreshSessionSvc.GetRefresh(ctx, refreshTokenID)
	if err != nil {
		return nil, fmt.Errorf("get refresh session: %w", err)
	}

	newRefresh, err := d.refreshSessionSvc.ExchangeRefresh(ctx, currRefresh)
	if err != nil {
		return nil, fmt.Errorf("exchange refresh: %w", err)
	}

	jwtStr, exp, err := d.Sign(newRefresh.UserID)
	if err != nil {
		return nil, fmt.Errorf("sign jwt: %w", err)
	}

	accessToken := AccessToken{
		AccessToken:  jwtStr,
		RefreshToken: newRefresh.ID.String(),
		ExpiresIn:    int(time.Until(exp).Seconds()),
	}

	return accessToken, nil
}
