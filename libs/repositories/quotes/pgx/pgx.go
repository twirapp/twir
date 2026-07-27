package pgx

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories"
	"github.com/twirapp/twir/libs/repositories/quotes"
	"github.com/twirapp/twir/libs/repositories/quotes/model"
)

type Opts struct {
	PgxPool *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{
		pool: opts.PgxPool,
	}
}

func NewFx(pgxpool *pgxpool.Pool) *Pgx {
	return New(Opts{PgxPool: pgxpool})
}

var _ quotes.Repository = (*Pgx)(nil)

var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool *pgxpool.Pool
}

func (c *Pgx) GetAllByChannelID(ctx context.Context, channelID uuid.UUID) ([]model.Quote, error) {
	query := `
SELECT id, "channelId", number, text, "creatorId", "creatorName", "gameId", "gameName", "createdAt", "updatedAt"
FROM channels_quotes
WHERE "channelId" = $1
`

	rows, err := c.pool.Query(ctx, query, channelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Quote])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (model.Quote, error) {
	query := `
SELECT id, "channelId", number, text, "creatorId", "creatorName", "gameId", "gameName", "createdAt", "updatedAt"
FROM channels_quotes
WHERE id = $1
LIMIT 1;
`

	rows, err := c.pool.Query(ctx, query, id)
	if err != nil {
		return model.Nil, err
	}
	defer rows.Close()

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Quote])
	if err != nil {
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) Create(ctx context.Context, input quotes.CreateInput) (model.Quote, error) {
	query := `
INSERT INTO channels_quotes ("channelId", number, text, "creatorId", "creatorName", "gameId", "gameName")
SELECT $1, COALESCE(MAX(number), 0) + 1, $2, $3, $4, $5, $6
FROM channels_quotes
WHERE "channelId" = $1
RETURNING id, "channelId", number, text, "creatorId", "creatorName", "gameId", "gameName", "createdAt", "updatedAt"
`

	rows, err := c.pool.Query(
		ctx,
		query,
		input.ChannelID,
		input.Text,
		input.CreatorID,
		input.CreatorName,
		input.GameID,
		input.GameName,
	)
	if err != nil {
		return model.Nil, err
	}
	defer rows.Close()

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Quote])
	if err != nil {
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) Update(ctx context.Context, id uuid.UUID, input quotes.UpdateInput) (model.Quote, error) {
	updateBuilder := sq.Update("channels_quotes")
	updateBuilder = repositories.SquirrelApplyPatch(
		updateBuilder, map[string]any{
			"text": input.Text,
		},
	)
	updateBuilder = updateBuilder.Set(`"updatedAt"`, squirrel.Expr("now()"))
	updateBuilder = updateBuilder.Where(squirrel.Eq{"id": id})

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return model.Nil, err
	}

	_, err = c.pool.Query(ctx, query, args...)
	if err != nil {
		return model.Nil, err
	}

	return c.GetByID(ctx, id)
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
DELETE FROM channels_quotes
WHERE id = $1
`

	rows, err := c.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if rows.RowsAffected() != 1 {
		return quotes.ErrQuoteNotFound
	}

	return nil
}
