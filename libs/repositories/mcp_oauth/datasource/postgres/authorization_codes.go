package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	entity "github.com/twirapp/twir/libs/entities/mcp_oauth"
	mcpOAuth "github.com/twirapp/twir/libs/repositories/mcp_oauth"
)

func (repository *Pgx) CreateAuthorizationCode(
	ctx context.Context,
	input mcpOAuth.CreateAuthorizationCodeInput,
) (entity.AuthorizationCode, error) {
	model, err := scanAuthorizationCode(repository.pool.QueryRow(
		ctx,
		`INSERT INTO mcp_oauth_authorization_codes (
			code_hash, client_id, channel_id, user_id, redirect_uri, pkce_challenge,
			scopes, resource, expires_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING `+authorizationCodeColumns,
		input.CodeHash.Bytes(),
		input.ClientID,
		input.ChannelID,
		input.UserID,
		input.RedirectURI,
		input.PKCEChallenge,
		scopesToStrings(input.Scopes),
		input.Resource,
		input.ExpiresAt,
	))
	if err != nil {
		return entity.NilAuthorizationCode, fmt.Errorf("create MCP OAuth authorization code: %w", err)
	}

	code, err := mapAuthorizationCode(model)
	if err != nil {
		return entity.NilAuthorizationCode, err
	}
	return code, nil
}

func (repository *Pgx) ConsumeAuthorizationCode(
	ctx context.Context,
	codeHash entity.CredentialHash,
) (entity.AuthorizationCode, error) {
	model, err := scanAuthorizationCode(repository.pool.QueryRow(
		ctx,
		`DELETE FROM mcp_oauth_authorization_codes
		WHERE code_hash = $1 AND expires_at > NOW()
		RETURNING `+authorizationCodeColumns,
		codeHash.Bytes(),
	))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.NilAuthorizationCode, mcpOAuth.ErrAuthorizationCodeNotFound
		}
		return entity.NilAuthorizationCode, fmt.Errorf("consume MCP OAuth authorization code: %w", err)
	}

	code, err := mapAuthorizationCode(model)
	if err != nil {
		return entity.NilAuthorizationCode, err
	}
	return code, nil
}

func scanAuthorizationCode(row pgx.Row) (authorizationCodeModel, error) {
	var model authorizationCodeModel
	err := row.Scan(
		&model.ID,
		&model.CodeHash,
		&model.ClientID,
		&model.ChannelID,
		&model.UserID,
		&model.RedirectURI,
		&model.PKCEChallenge,
		&model.Scopes,
		&model.Resource,
		&model.ExpiresAt,
		&model.CreatedAt,
	)
	return model, err
}
