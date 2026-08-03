package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	entity "github.com/twirapp/twir/libs/entities/mcp_oauth"
	mcpOAuth "github.com/twirapp/twir/libs/repositories/mcp_oauth"
)

func (repository *Pgx) CreateClient(ctx context.Context, input mcpOAuth.CreateClientInput) (entity.Client, error) {
	model, err := scanClient(repository.pool.QueryRow(
		ctx,
		`INSERT INTO mcp_oauth_clients (
			client_id, metadata, redirect_uris, grant_types, response_types,
			token_endpoint_auth_method, scopes
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING `+clientColumns,
		input.ClientID,
		input.Metadata,
		input.RedirectURIs,
		input.GrantTypes,
		input.ResponseTypes,
		input.TokenEndpointAuthMethod,
		scopesToStrings(input.Scopes),
	))
	if err != nil {
		if isClientIDConflict(err) {
			return entity.NilClient, mcpOAuth.ErrClientAlreadyExists
		}
		return entity.NilClient, fmt.Errorf("create MCP OAuth client: %w", err)
	}

	return mapClient(model), nil
}

func (repository *Pgx) GetClient(ctx context.Context, clientID string) (entity.Client, error) {
	model, err := scanClient(repository.pool.QueryRow(
		ctx,
		`SELECT `+clientColumns+` FROM mcp_oauth_clients WHERE client_id = $1`,
		clientID,
	))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.NilClient, mcpOAuth.ErrClientNotFound
		}
		return entity.NilClient, fmt.Errorf("get MCP OAuth client: %w", err)
	}

	return mapClient(model), nil
}

func scanClient(row pgx.Row) (clientModel, error) {
	var model clientModel
	err := row.Scan(
		&model.ID,
		&model.ClientID,
		&model.Metadata,
		&model.RedirectURIs,
		&model.GrantTypes,
		&model.ResponseTypes,
		&model.TokenEndpointAuthMethod,
		&model.Scopes,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	return model, err
}

func isClientIDConflict(err error) bool {
	var databaseError *pgconn.PgError
	return errors.As(err, &databaseError) && databaseError.ConstraintName == "mcp_oauth_clients_client_id_key"
}
