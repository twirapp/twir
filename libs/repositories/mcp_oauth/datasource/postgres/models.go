package postgres

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	entity "github.com/twirapp/twir/libs/entities/mcp_oauth"
)

const clientColumns = `
	id, client_id, metadata, redirect_uris, grant_types, response_types,
	token_endpoint_auth_method, scopes, created_at, updated_at`

const authorizationCodeColumns = `
	id, code_hash, client_id, channel_id, user_id, redirect_uri, pkce_challenge,
	scopes, resource, expires_at, created_at`

const tokenColumns = `
	id, family_id, client_id, channel_id, user_id, access_token_hash, refresh_token_hash,
	scopes, resource, access_expires_at, refresh_expires_at, revoked_at, replaced_by_id,
	created_at, updated_at`

type clientModel struct {
	ID                      uuid.UUID
	ClientID                string
	Metadata                []byte
	RedirectURIs            []string
	GrantTypes              []string
	ResponseTypes           []string
	TokenEndpointAuthMethod string
	Scopes                  []string
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type authorizationCodeModel struct {
	ID            uuid.UUID
	CodeHash      []byte
	ClientID      string
	ChannelID     uuid.UUID
	UserID        uuid.UUID
	RedirectURI   string
	PKCEChallenge string
	Scopes        []string
	Resource      string
	ExpiresAt     time.Time
	CreatedAt     time.Time
}

type tokenModel struct {
	ID               uuid.UUID
	FamilyID         uuid.UUID
	ClientID         string
	ChannelID        uuid.UUID
	UserID           uuid.UUID
	AccessTokenHash  []byte
	RefreshTokenHash []byte
	Scopes           []string
	Resource         string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
	RevokedAt        *time.Time
	ReplacedByID     *uuid.UUID
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func mapClient(model clientModel) entity.Client {
	return entity.Client{
		ID:                      model.ID,
		ClientID:                model.ClientID,
		Metadata:                model.Metadata,
		RedirectURIs:            model.RedirectURIs,
		GrantTypes:              model.GrantTypes,
		ResponseTypes:           model.ResponseTypes,
		TokenEndpointAuthMethod: model.TokenEndpointAuthMethod,
		Scopes:                  scopesFromStrings(model.Scopes),
		CreatedAt:               model.CreatedAt,
		UpdatedAt:               model.UpdatedAt,
	}
}

func mapAuthorizationCode(model authorizationCodeModel) (entity.AuthorizationCode, error) {
	codeHash, err := credentialHashFromBytes(model.CodeHash)
	if err != nil {
		return entity.NilAuthorizationCode, fmt.Errorf("map authorization code hash: %w", err)
	}

	return entity.AuthorizationCode{
		ID:            model.ID,
		CodeHash:      codeHash,
		ClientID:      model.ClientID,
		ChannelID:     model.ChannelID,
		UserID:        model.UserID,
		RedirectURI:   model.RedirectURI,
		PKCEChallenge: model.PKCEChallenge,
		Scopes:        scopesFromStrings(model.Scopes),
		Resource:      model.Resource,
		ExpiresAt:     model.ExpiresAt,
		CreatedAt:     model.CreatedAt,
	}, nil
}

func mapToken(model tokenModel) (entity.Token, error) {
	accessTokenHash, err := credentialHashFromBytes(model.AccessTokenHash)
	if err != nil {
		return entity.NilToken, fmt.Errorf("map access token hash: %w", err)
	}
	refreshTokenHash, err := credentialHashFromBytes(model.RefreshTokenHash)
	if err != nil {
		return entity.NilToken, fmt.Errorf("map refresh token hash: %w", err)
	}

	return entity.Token{
		ID:               model.ID,
		FamilyID:         model.FamilyID,
		ClientID:         model.ClientID,
		ChannelID:        model.ChannelID,
		UserID:           model.UserID,
		AccessTokenHash:  accessTokenHash,
		RefreshTokenHash: refreshTokenHash,
		Scopes:           scopesFromStrings(model.Scopes),
		Resource:         model.Resource,
		AccessExpiresAt:  model.AccessExpiresAt,
		RefreshExpiresAt: model.RefreshExpiresAt,
		RevokedAt:        model.RevokedAt,
		ReplacedByID:     model.ReplacedByID,
		CreatedAt:        model.CreatedAt,
		UpdatedAt:        model.UpdatedAt,
	}, nil
}

func scopesFromStrings(values []string) []entity.Scope {
	scopes := make([]entity.Scope, len(values))
	for index, value := range values {
		scopes[index] = entity.Scope(value)
	}
	return scopes
}

func scopesToStrings(values []entity.Scope) []string {
	scopes := make([]string, len(values))
	for index, value := range values {
		scopes[index] = string(value)
	}
	return scopes
}

func credentialHashFromBytes(value []byte) (entity.CredentialHash, error) {
	var hash entity.CredentialHash
	if len(value) != len(hash) {
		return entity.CredentialHash{}, fmt.Errorf("credential hash length %d, want %d", len(value), len(hash))
	}
	copy(hash[:], value)
	return hash, nil
}
