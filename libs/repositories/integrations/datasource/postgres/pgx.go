package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/integrations"
	"github.com/twirapp/twir/libs/repositories/integrations/model"
)

var _ integrations.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool *pgxpool.Pool
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return New(pool)
}

func New(pool *pgxpool.Pool) *Pgx {
	return &Pgx{
		pool: pool,
	}
}

func (p *Pgx) GetByService(ctx context.Context, service model.Service) (model.Integration, error) {
	query := `
SELECT
    id,
    service,
    "accessToken",
    "refreshToken",
    "clientId",
    "clientSecret",
    "apiKey",
    "redirectUrl"
FROM integrations
WHERE service = $1
LIMIT 1
`

	rows, err := p.pool.Query(ctx, query, service)
	if err != nil {
		return model.Nil, fmt.Errorf("GetByService: failed to execute query: %w", err)
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Integration])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, fmt.Errorf("integration not found for service %s", service)
		}
		return model.Nil, fmt.Errorf("GetByService: failed to collect row: %w", err)
	}

	return result, nil
}
