package oauth

import (
	"context"
	"fmt"
	"time"
)

type AppTokenKey struct {
	Provider Provider
	ID       CredentialID
}

func (key AppTokenKey) Validate() error {
	if key.Provider == "" || key.ID == "" {
		return fmt.Errorf("%w: app provider and token ID are required", ErrInvalidCredential)
	}
	return nil
}

type AppToken struct {
	AccessToken string
	ObtainedAt  time.Time
	ExpiresIn   time.Duration
}

func (token AppToken) Validate() error {
	if token.AccessToken == "" || token.ObtainedAt.IsZero() || token.ExpiresIn <= 0 {
		return fmt.Errorf("%w: app access token and positive expiry are required", ErrInvalidCredential)
	}
	return nil
}

func (token AppToken) Expired(now time.Time, skew time.Duration) bool {
	if skew < 0 {
		skew = 0
	}
	return token.Validate() != nil || !now.Before(token.ObtainedAt.Add(token.ExpiresIn).Add(-skew))
}

type AppTokenStore interface {
	LoadAppToken(context.Context, AppTokenKey) (AppToken, error)
	CommitAppToken(context.Context, AppTokenKey, AppToken) error
}

type AppTokenFetcher interface {
	FetchAppToken(context.Context, AppTokenKey) (AppToken, error)
}

type AppTokenLocker interface {
	AcquireAppToken(context.Context, AppTokenKey) (Lease, error)
}

type AppTokenDependencies struct {
	Store   AppTokenStore
	Fetcher AppTokenFetcher
	Locker  AppTokenLocker
}

type AppTokenSourceOptions struct {
	Clock       Clock
	Skew        time.Duration
	WorkTimeout time.Duration
}
