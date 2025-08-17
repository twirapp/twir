package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/tokens"
	tokenmodel "github.com/twirapp/twir/libs/repositories/tokens/model"
)

type Opts struct {
	PgxPool *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{
		pool:   opts.PgxPool,
		getter: trmpgx.DefaultCtxGetter,
	}
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return New(Opts{PgxPool: pool})
}

var _ tokens.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (*tokenmodel.Token, error) {
	query := `
SELECT id, "accessToken", "refreshToken", "expiresIn", "obtainmentTimestamp", scopes
FROM tokens
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[tokenmodel.Token])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, tokens.ErrNotFound
		}
		return nil, err // Other error
	}

	return &result, nil
}

func (c *Pgx) GetByUserID(ctx context.Context, userID string) (*tokenmodel.Token, error) {
	query := `
SELECT token.id, "accessToken", "refreshToken", "expiresIn", "obtainmentTimestamp", scopes
FROM users
JOIN tokens AS token ON users."tokenId" = token.id
WHERE users.id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[tokenmodel.Token])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, tokens.ErrNotFound
		}
		return nil, err // Other error
	}

	return &result, nil
}

func (c *Pgx) GetByBotID(ctx context.Context, botID string) (*tokenmodel.Token, error) {
	query := `
SELECT token.id, "accessToken", "refreshToken", "expiresIn", "obtainmentTimestamp", scopes
FROM bots
JOIN tokens AS token ON bots."tokenId" = token.id
WHERE bots.id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, botID)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[tokenmodel.Token])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, tokens.ErrNotFound
		}
		return nil, err // Other error
	}

	return &result, nil
}

func (c *Pgx) CreateUserToken(ctx context.Context, input tokens.CreateInput) (
	*tokenmodel.Token,
	error,
) {
	query := `
INSERT INTO tokens ("accessToken", "refreshToken", "expiresIn", "obtainmentTimestamp", scopes)
VALUES ($1, $2, $3, $4, $5)
RETURNING id
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	var id uuid.UUID
	err := conn.QueryRow(
		ctx, query,
		input.AccessToken,
		input.RefreshToken,
		input.ExpiresIn,
		input.ObtainmentTimestamp,
		input.Scopes,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to create user token: %w", err)
	}

	newToken, err := c.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created token: %w", err)
	}

	return newToken, nil
}

func (c *Pgx) UpdateTokenByID(
	ctx context.Context,
	id uuid.UUID,
	input tokens.UpdateTokenInput,
) (*tokenmodel.Token, error) {
	updateBuilder := sq.Update("tokens")

	if input.AccessToken != nil {
		updateBuilder = updateBuilder.Set(`"accessToken"`, *input.AccessToken)
	}

	if input.RefreshToken != nil {
		updateBuilder = updateBuilder.Set(`"refreshToken"`, *input.RefreshToken)
	}

	if input.ExpiresIn != nil {
		updateBuilder = updateBuilder.Set(`"expiresIn"`, *input.ExpiresIn)
	}

	if input.ObtainmentTimestamp != nil {
		updateBuilder = updateBuilder.Set(`"obtainmentTimestamp"`, *input.ObtainmentTimestamp)
	}

	if len(input.Scopes) > 0 {
		updateBuilder = updateBuilder.Set("scopes", input.Scopes)
	}

	updateBuilder = updateBuilder.Where(squirrel.Eq{"id": id})
	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build update query: %w", err)
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update user token: %w", err)
	}

	updatedToken, err := c.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated token: %w", err)
	}

	return updatedToken, nil
}
