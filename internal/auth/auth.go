package auth

import (
	"errors"
)

var (
	ErrPasswordIncorrect        = errors.New("password incorrect")
	ErrJWTExpired               = errors.New("JWT is expired")
	ErrJWTInvalid               = errors.New("JWT is invalid")
	ErrRefreshTokenExpired      = errors.New("refresh token is expired")
	ErrRefreshTokenUserMismatch = errors.New("token user does not match")
	ErrRefreshTokenUsed         = errors.New("refresh token is used")
	ErrRefreshTokenInvalid      = errors.New("refresh token is used")
)
