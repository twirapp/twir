package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	entity "github.com/twirapp/twir/libs/entities/mcp_oauth"
	mcpOAuth "github.com/twirapp/twir/libs/repositories/mcp_oauth"
)

func (repository *Pgx) CreateToken(ctx context.Context, input mcpOAuth.CreateTokenInput) (entity.Token, error) {
	model, err := insertInitialToken(ctx, repository.pool, input)
	if err != nil {
		return entity.NilToken, fmt.Errorf("create MCP OAuth token: %w", err)
	}

	token, err := mapToken(model)
	if err != nil {
		return entity.NilToken, err
	}
	return token, nil
}

func (repository *Pgx) RotateRefreshToken(
	ctx context.Context,
	input mcpOAuth.RotateRefreshTokenInput,
) (entity.Token, error) {
	tx, err := repository.pool.Begin(ctx)
	if err != nil {
		return entity.NilToken, fmt.Errorf("begin MCP OAuth refresh token rotation: %w", err)
	}
	defer tx.Rollback(ctx)

	current, err := findRefreshTokenForUpdate(ctx, tx, input.PresentedRefreshTokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.NilToken, mcpOAuth.ErrRefreshTokenNotFound
		}
		return entity.NilToken, fmt.Errorf("lock MCP OAuth refresh token: %w", err)
	}

	if current.RevokedAt != nil || current.ReplacedByID != nil {
		return repository.revokeReusedFamily(ctx, tx, current.FamilyID)
	}

	nextTokenID, err := uuid.NewV7()
	if err != nil {
		return entity.NilToken, fmt.Errorf("generate rotated MCP OAuth token ID: %w", err)
	}

	next, err := insertRotatedToken(
		ctx,
		tx,
		nextTokenID,
		current.FamilyID,
		mcpOAuth.CreateTokenInput{
			ClientID:         current.ClientID,
			ChannelID:        current.ChannelID,
			UserID:           current.UserID,
			AccessTokenHash:  input.NextAccessTokenHash,
			RefreshTokenHash: input.NextRefreshTokenHash,
			Scopes:           input.Scopes,
			Resource:         current.Resource,
			AccessExpiresAt:  input.AccessExpiresAt,
			RefreshExpiresAt: input.RefreshExpiresAt,
		},
	)
	if err != nil {
		return entity.NilToken, fmt.Errorf("create rotated MCP OAuth token: %w", err)
	}

	updated, err := tx.Exec(
		ctx,
		`UPDATE mcp_oauth_tokens
		SET replaced_by_id = $2, updated_at = NOW()
		WHERE id = $1
			AND revoked_at IS NULL
			AND replaced_by_id IS NULL
			AND refresh_expires_at > NOW()`,
		current.ID,
		nextTokenID,
	)
	if err != nil {
		return entity.NilToken, fmt.Errorf("mark MCP OAuth refresh token rotated: %w", err)
	}
	if updated.RowsAffected() == 0 {
		if _, err := tx.Exec(ctx, `DELETE FROM mcp_oauth_tokens WHERE id = $1`, nextTokenID); err != nil {
			return entity.NilToken, fmt.Errorf("delete unrotated MCP OAuth token: %w", err)
		}
		return repository.revokeReusedFamily(ctx, tx, current.FamilyID)
	}

	if err := tx.Commit(ctx); err != nil {
		return entity.NilToken, fmt.Errorf("commit MCP OAuth refresh token rotation: %w", err)
	}

	token, err := mapToken(next)
	if err != nil {
		return entity.NilToken, err
	}
	return token, nil
}

func (repository *Pgx) RevokeToken(ctx context.Context, clientID string, tokenHash entity.CredentialHash) error {
	tx, err := repository.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin MCP OAuth token revocation: %w", err)
	}
	defer tx.Rollback(ctx)

	var familyID uuid.UUID
	err = tx.QueryRow(
		ctx,
		`SELECT family_id
		FROM mcp_oauth_tokens
		WHERE client_id = $1 AND (access_token_hash = $2 OR refresh_token_hash = $3)
		LIMIT 1
		FOR UPDATE`,
		clientID,
		tokenHash.Bytes(),
		tokenHash.Bytes(),
	).Scan(&familyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("lock MCP OAuth token family for revocation: %w", err)
	}

	if err := revokeFamilyInTx(ctx, tx, familyID); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit MCP OAuth token revocation: %w", err)
	}
	return nil
}

func (repository *Pgx) revokeReusedFamily(ctx context.Context, tx pgx.Tx, familyID uuid.UUID) (entity.Token, error) {
	if err := revokeFamilyInTx(ctx, tx, familyID); err != nil {
		return entity.NilToken, err
	}
	if err := tx.Commit(ctx); err != nil {
		return entity.NilToken, fmt.Errorf("commit MCP OAuth refresh token family revocation: %w", err)
	}
	return entity.NilToken, &mcpOAuth.RefreshTokenReuseError{FamilyID: familyID}
}

func findRefreshTokenForUpdate(
	ctx context.Context,
	tx pgx.Tx,
	refreshTokenHash entity.CredentialHash,
) (tokenModel, error) {
	return scanToken(tx.QueryRow(
		ctx,
		`SELECT `+tokenColumns+`
		FROM mcp_oauth_tokens
		WHERE refresh_token_hash = $1
		FOR UPDATE`,
		refreshTokenHash.Bytes(),
	))
}

func insertInitialToken(
	ctx context.Context,
	pool *pgxpool.Pool,
	input mcpOAuth.CreateTokenInput,
) (tokenModel, error) {
	return scanToken(pool.QueryRow(
		ctx,
		`INSERT INTO mcp_oauth_tokens (
			client_id, channel_id, user_id, access_token_hash, refresh_token_hash,
			scopes, resource, access_expires_at, refresh_expires_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING `+tokenColumns,
		input.ClientID,
		input.ChannelID,
		input.UserID,
		input.AccessTokenHash.Bytes(),
		input.RefreshTokenHash.Bytes(),
		scopesToStrings(input.Scopes),
		input.Resource,
		input.AccessExpiresAt,
		input.RefreshExpiresAt,
	))
}

func insertRotatedToken(
	ctx context.Context,
	tx pgx.Tx,
	id uuid.UUID,
	familyID uuid.UUID,
	input mcpOAuth.CreateTokenInput,
) (tokenModel, error) {
	return scanToken(tx.QueryRow(
		ctx,
		`INSERT INTO mcp_oauth_tokens (
			id, family_id, client_id, channel_id, user_id, access_token_hash,
			refresh_token_hash, scopes, resource, access_expires_at, refresh_expires_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING `+tokenColumns,
		id,
		familyID,
		input.ClientID,
		input.ChannelID,
		input.UserID,
		input.AccessTokenHash.Bytes(),
		input.RefreshTokenHash.Bytes(),
		scopesToStrings(input.Scopes),
		input.Resource,
		input.AccessExpiresAt,
		input.RefreshExpiresAt,
	))
}

func revokeFamilyInTx(ctx context.Context, tx pgx.Tx, familyID uuid.UUID) error {
	if _, err := tx.Exec(
		ctx,
		`UPDATE mcp_oauth_tokens
		SET revoked_at = NOW(), updated_at = NOW()
		WHERE family_id = $1 AND revoked_at IS NULL`,
		familyID,
	); err != nil {
		return fmt.Errorf("revoke MCP OAuth token family: %w", err)
	}
	return nil
}

func scanToken(row pgx.Row) (tokenModel, error) {
	var model tokenModel
	err := row.Scan(
		&model.ID,
		&model.FamilyID,
		&model.ClientID,
		&model.ChannelID,
		&model.UserID,
		&model.AccessTokenHash,
		&model.RefreshTokenHash,
		&model.Scopes,
		&model.Resource,
		&model.AccessExpiresAt,
		&model.RefreshExpiresAt,
		&model.RevokedAt,
		&model.ReplacedByID,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	return model, err
}
