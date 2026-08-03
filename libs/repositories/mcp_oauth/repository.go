package mcp_oauth

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	entity "github.com/twirapp/twir/libs/entities/mcp_oauth"
)

type Repository interface {
	CreateClient(ctx context.Context, input CreateClientInput) (entity.Client, error)
	GetClient(ctx context.Context, clientID string) (entity.Client, error)
	CreateAuthorizationCode(ctx context.Context, input CreateAuthorizationCodeInput) (entity.AuthorizationCode, error)
	ConsumeAuthorizationCode(ctx context.Context, codeHash entity.CredentialHash) (entity.AuthorizationCode, error)
	CreateToken(ctx context.Context, input CreateTokenInput) (entity.Token, error)
	GetTokenByAccessTokenHash(ctx context.Context, tokenHash entity.CredentialHash) (entity.Token, error)
	GetTokenByRefreshTokenHash(ctx context.Context, tokenHash entity.CredentialHash) (entity.Token, error)
	RotateRefreshToken(ctx context.Context, input RotateRefreshTokenInput) (entity.Token, error)
	RevokeToken(ctx context.Context, clientID string, tokenHash entity.CredentialHash) error
}

type CreateClientInput struct {
	ClientID                string
	Metadata                json.RawMessage
	RedirectURIs            []string
	GrantTypes              []string
	ResponseTypes           []string
	TokenEndpointAuthMethod string
	Scopes                  []entity.Scope
}

type CreateAuthorizationCodeInput struct {
	CodeHash      entity.CredentialHash
	ClientID      string
	ChannelID     uuid.UUID
	UserID        uuid.UUID
	RedirectURI   string
	PKCEChallenge string
	Scopes        []entity.Scope
	Resource      string
	ExpiresAt     time.Time
}

type CreateTokenInput struct {
	ClientID         string
	ChannelID        uuid.UUID
	UserID           uuid.UUID
	AccessTokenHash  entity.CredentialHash
	RefreshTokenHash entity.CredentialHash
	Scopes           []entity.Scope
	Resource         string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
}

type RotateRefreshTokenInput struct {
	PresentedRefreshTokenHash entity.CredentialHash
	NextAccessTokenHash       entity.CredentialHash
	NextRefreshTokenHash      entity.CredentialHash
	Scopes                    []entity.Scope
	AccessExpiresAt           time.Time
	RefreshExpiresAt          time.Time
}
