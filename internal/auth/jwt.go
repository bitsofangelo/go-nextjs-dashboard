package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/gelozr/go-dash/internal/config"
)

type JWTClaims struct {
	UserID uuid.UUID
	jwt.RegisteredClaims
}

type GOJWT struct {
	hmacKey []byte
}

var _ JWT = (*GOJWT)(nil)

func NewGOJWT(cfg *config.Config) *GOJWT {
	return &GOJWT{
		hmacKey: []byte(cfg.JWTHmacKey),
	}
}

func (g *GOJWT) Sign(uid uuid.UUID) (string, time.Time, error) {
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

	signed, err := tok.SignedString(g.hmacKey)
	if err != nil {
		return "", time.Now(), fmt.Errorf("signing token: %w", err)
	}

	exp, err := tok.Claims.GetExpirationTime()
	if err != nil {
		return "", time.Now(), fmt.Errorf("getting expiration time: %w", err)
	}

	return signed, exp.Time, nil
}

func (g *GOJWT) Parse(tokenStr string) (AccessClaims, error) {
	tok, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (any, error) {
		return g.hmacKey, nil
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
