package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var hmacKey = []byte(os.Getenv("JWT_HMAC_KEY"))

type AccessClaims struct {
	UserID uuid.UUID
	jwt.RegisteredClaims
}

type GOJWT struct{}

func NewGOJWT() *GOJWT {
	return &GOJWT{}
}

func (g *GOJWT) NewAccess(uid uuid.UUID) (string, time.Time, error) {
	claims := AccessClaims{
		UserID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(), // jti
			Issuer:    "myapi",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := tok.SignedString(hmacKey)
	if err != nil {
		return "", time.Now(), fmt.Errorf("signing token: %w", err)
	}

	exp, err := tok.Claims.GetExpirationTime()
	if err != nil {
		return "", time.Now(), fmt.Errorf("getting expiration time: %w", err)
	}

	return signed, exp.Time, nil
}

func (g *GOJWT) ParseAccess(tokenStr string) (Claims, error) {
	tok, err := jwt.ParseWithClaims(tokenStr, &AccessClaims{}, func(t *jwt.Token) (any, error) {
		return hmacKey, nil
	})

	if err != nil {
		return Claims{}, err // signature, exp, nbf, etc.
	}

	ac := tok.Claims.(*AccessClaims)

	claims := Claims{
		Issuer:    ac.Issuer,
		Subject:   ac.Subject,
		Audience:  ac.Audience,
		ExpiresAt: ac.ExpiresAt.Time,
		// NotBefore: ac.NotBefore.Time,
		IssuedAt: ac.IssuedAt.Time,
		ID:       ac.ID,
		UserID:   ac.UserID,
	}

	// OPTIONAL: check deny-list
	// if revoked, _ := redisClient.Get(ctx, "block:"+claims.ID).Result(); revoked == "1" {
	// 	return nil, fmt.Errorf("token revoked")
	// }

	return claims, nil
}
