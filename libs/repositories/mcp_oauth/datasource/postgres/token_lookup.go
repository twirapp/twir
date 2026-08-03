package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	entity "github.com/twirapp/twir/libs/entities/mcp_oauth"
	mcpOAuth "github.com/twirapp/twir/libs/repositories/mcp_oauth"
)

func (repository *Pgx) GetTokenByAccessTokenHash(ctx context.Context, hash entity.CredentialHash) (entity.Token, error) {
	return repository.getTokenByHash(ctx, "access_token_hash", hash, mcpOAuth.ErrAccessTokenNotFound)
}

func (repository *Pgx) GetTokenByRefreshTokenHash(ctx context.Context, hash entity.CredentialHash) (entity.Token, error) {
	return repository.getTokenByHash(ctx, "refresh_token_hash", hash, mcpOAuth.ErrRefreshTokenNotFound)
}

func (repository *Pgx) getTokenByHash(ctx context.Context, column string, hash entity.CredentialHash, notFound error) (entity.Token, error) {
	model, err := scanToken(repository.pool.QueryRow(ctx, `SELECT `+tokenColumns+` FROM mcp_oauth_tokens WHERE `+column+` = $1`, hash.Bytes()))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.NilToken, notFound
		}
		return entity.NilToken, fmt.Errorf("get MCP OAuth token: %w", err)
	}
	return mapToken(model)
}
