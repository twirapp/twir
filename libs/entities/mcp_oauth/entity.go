package mcp_oauth

import (
	"crypto/sha256"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Scope string

const (
	// ScopeRead and ScopeWrite are legacy compatibility aliases for the full read and edit scope sets.
	ScopeRead  Scope = "read"
	ScopeWrite Scope = "write"
)

type CredentialHash [sha256.Size]byte

func (hash CredentialHash) Bytes() []byte {
	return hash[:]
}

type Client struct {
	ID                      uuid.UUID
	ClientID                string
	Metadata                json.RawMessage
	RedirectURIs            []string
	GrantTypes              []string
	ResponseTypes           []string
	TokenEndpointAuthMethod string
	Scopes                  []Scope
	CreatedAt               time.Time
	UpdatedAt               time.Time

	isNil bool
}

func (client Client) IsNil() bool {
	return client.isNil
}

var NilClient = Client{isNil: true}

type AuthorizationCode struct {
	ID            uuid.UUID
	CodeHash      CredentialHash
	ClientID      string
	ChannelID     uuid.UUID
	UserID        uuid.UUID
	RedirectURI   string
	PKCEChallenge string
	Scopes        []Scope
	Resource      string
	ExpiresAt     time.Time
	CreatedAt     time.Time

	isNil bool
}

func (code AuthorizationCode) IsNil() bool {
	return code.isNil
}

var NilAuthorizationCode = AuthorizationCode{isNil: true}

type Token struct {
	ID               uuid.UUID
	FamilyID         uuid.UUID
	ClientID         string
	ChannelID        uuid.UUID
	UserID           uuid.UUID
	AccessTokenHash  CredentialHash
	RefreshTokenHash CredentialHash
	Scopes           []Scope
	Resource         string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
	RevokedAt        *time.Time
	ReplacedByID     *uuid.UUID
	CreatedAt        time.Time
	UpdatedAt        time.Time

	isNil bool
}

func (token Token) IsNil() bool {
	return token.isNil
}

var NilToken = Token{isNil: true}
