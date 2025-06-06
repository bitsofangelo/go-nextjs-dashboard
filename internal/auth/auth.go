package auth

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrLoginNotSupported        = errors.New("login not supported")
	ErrLogoutNotSupported       = errors.New("logout not supported")
	ErrRefreshTokenNotSupported = errors.New("logout not supported")
	ErrPasswordIncorrect        = errors.New("password incorrect")
	ErrJWTExpired               = errors.New("JWTDriver is expired")
	ErrJWTInvalid               = errors.New("JWTDriver is invalid")
	ErrRefreshTokenExpired      = errors.New("refresh token is expired")
	ErrRefreshTokenUserMismatch = errors.New("token user does not match")
	ErrRefreshTokenUsed         = errors.New("refresh token is used")
	ErrRefreshTokenInvalid      = errors.New("refresh token is used")
)

type GuardOption struct {
	Driver       Driver
	UserProvider UserProvider
}

type guard struct {
	userProvider UserProvider
	driver       Driver
}

var _ Guard = (*guard)(nil)

func (g guard) Authenticate(ctx context.Context, creds Credentials) (User, error) {
	user, err := g.userProvider.FindByCredentials(ctx, creds)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (g guard) Login(ctx context.Context, user User) (any, error) {
	if d, ok := g.driver.(LoginHandler); ok {
		return d.Login(ctx, user)
	}
	return nil, ErrLoginNotSupported
}

func (g guard) Logout(ctx context.Context, id any) error {
	if d, ok := g.driver.(LogoutHandler); ok {
		return d.Logout(ctx, id)
	}
	return ErrLogoutNotSupported
}

func (g guard) Check(ctx context.Context, payload any) (Verified, error) {
	verified, err := g.driver.Verify(ctx, payload)
	if err != nil {
		return Verified{}, err
	}

	verified.ctx = WithUserCtx(ctx, verified.User)

	return verified, nil
}

func (g guard) RefreshToken(ctx context.Context, refreshToken string) (any, error) {
	if d, ok := g.driver.(TokenRefresher); ok {
		return d.RefreshToken(ctx, refreshToken)
	}
	return nil, ErrRefreshTokenNotSupported
}

type Provider struct {
	mu           sync.RWMutex
	guards       map[string]*guard
	defaultGuard string
}

var _ Auth = (*Provider)(nil)

func New() *Provider {
	return &Provider{
		guards:       make(map[string]*guard),
		defaultGuard: "",
	}
}

func (a *Provider) Authenticate(ctx context.Context, creds Credentials) (User, error) {
	return a.MustGuard(a.defaultGuard).Authenticate(ctx, creds)
}

func (a *Provider) Login(ctx context.Context, user User) (any, error) {
	return a.MustGuard(a.defaultGuard).Login(ctx, user)
}

func (a *Provider) Logout(ctx context.Context, id any) error {
	return a.MustGuard(a.defaultGuard).Logout(ctx, id)
}

func (a *Provider) Check(ctx context.Context, payload any) (Verified, error) {
	return a.MustGuard(a.defaultGuard).Check(ctx, payload)
}

func (a *Provider) RefreshToken(ctx context.Context, refreshToken string) (any, error) {
	return a.MustGuard(a.defaultGuard).RefreshToken(ctx, refreshToken)
}

func (a *Provider) Guard(name string) (Guard, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if g, ok := a.guards[name]; ok {
		return g, nil
	}
	return nil, fmt.Errorf("guard '%s' not found", name)
}

func (a *Provider) MustGuard(name string) Guard {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if g, ok := a.guards[name]; ok {
		return g
	}

	panic(fmt.Sprintf("guard '%s' not found", name))
}

func (a *Provider) Extend(name string, option GuardOption) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if option.UserProvider == nil || option.Driver == nil {
		return errors.New("UserProvider or Driver cannot be nil")
	}

	if a.guards == nil {
		a.guards = make(map[string]*guard)
	}

	a.guards[name] = &guard{
		userProvider: option.UserProvider,
		driver:       option.Driver,
	}

	return nil
}

func (a *Provider) SetDefaultGuard(name string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, ok := a.guards[name]; !ok {
		return fmt.Errorf("guard '%s' not found", name)
	}

	a.defaultGuard = name
	return nil
}

type Token struct {
	refreshStore RefreshStore
	// logger       logger.Logger
}

func NewToken(refreshStore RefreshStore) *Token {
	return &Token{
		refreshStore: refreshStore,
		// logger:       logger.With("component", "auth"),
	}
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

type Verified struct {
	ctx    context.Context
	User   User
	Scopes []string
}

func (v Verified) Context() context.Context {
	return v.ctx
}

type ctxKey string

var userCtxKey = ctxKey("user")

func WithUserCtx(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

func UserFromCtx(ctx context.Context) (User, bool) {
	u, ok := ctx.Value(userCtxKey).(User)
	return u, ok
}
