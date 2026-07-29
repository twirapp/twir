package oauth

import (
	"context"
	"fmt"
	"time"
)

type Provider string
type CredentialID string

type CredentialKey struct {
	Provider Provider
	ID       CredentialID
}

type Credential struct {
	Provider     Provider
	ID           CredentialID
	AccessToken  string
	RefreshToken string
	Scopes       []string
	ObtainedAt   time.Time
	ExpiresIn    time.Duration
}

func (c Credential) Key() CredentialKey { return CredentialKey{Provider: c.Provider, ID: c.ID} }
func (k CredentialKey) Validate() error {
	if k.Provider == "" || k.ID == "" {
		return fmt.Errorf("%w: provider and credential ID are required", ErrInvalidCredential)
	}
	return nil
}
func (c Credential) Validate() error {
	if err := c.Key().Validate(); err != nil {
		return err
	}
	if c.AccessToken == "" || c.ObtainedAt.IsZero() || c.ExpiresIn <= 0 {
		return fmt.Errorf("%w: token and positive expiry are required", ErrInvalidCredential)
	}
	return nil
}
func (c Credential) Expired(now time.Time, skew time.Duration) bool {
	if c.ObtainedAt.IsZero() || c.ExpiresIn <= 0 {
		return true
	}
	if skew < 0 {
		skew = 0
	}
	return !now.Before(c.ObtainedAt.Add(c.ExpiresIn).Add(-skew))
}

type Store interface {
	Load(context.Context, CredentialKey) (Credential, error)
	Commit(context.Context, Credential) error
}
type Refresher interface {
	Refresh(context.Context, Credential) (RefreshResult, error)
}
type RefreshResult struct {
	AccessToken  string
	RefreshToken *string
	Scopes       []string
	ExpiresIn    time.Duration
}
type Locker interface {
	Acquire(context.Context, CredentialKey) (Lease, error)
}
type Lease interface {
	Context() context.Context
	Lost() <-chan struct{}
	Release(context.Context) error
}
type Clock interface{ Now() time.Time }
type systemClock struct{}

func (systemClock) Now() time.Time { return time.Now() }

type Observer interface{ Observe(context.Context, Event) }
type Event struct {
	Provider     Provider
	CredentialID CredentialID
	Operation    string
	Duration     time.Duration
	ErrorClass   string
}
