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

const quoteFields = `id, channel_id, number, text, creator_id, creator_name, game_id, game_name, created_at, updated_at`

type Pgx struct {
	pool *pgxpool.Pool
}

func (c *Pgx) GetAllByChannelID(ctx context.Context, channelID uuid.UUID) ([]model.Quote, error) {
	query := `SELECT ` + quoteFields + ` FROM channels_quotes WHERE channel_id = $1`

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
	query := `SELECT ` + quoteFields + ` FROM channels_quotes WHERE id = $1 LIMIT 1`

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

func (c *Pgx) GetByChannelIDAndNumber(
	ctx context.Context,
	channelID uuid.UUID,
	number int,
) (model.Quote, error) {
	query := `SELECT ` + quoteFields + ` FROM channels_quotes WHERE channel_id = $1 AND number = $2 LIMIT 1`

	rows, err := c.pool.Query(ctx, query, channelID, number)
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

func (c *Pgx) GetRandomByChannelID(ctx context.Context, channelID uuid.UUID) (model.Quote, error) {
	query := `SELECT ` + quoteFields + ` FROM channels_quotes WHERE channel_id = $1 ORDER BY RANDOM() LIMIT 1`

	rows, err := c.pool.Query(ctx, query, channelID)
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
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return model.Nil, err
	}
	defer tx.Rollback(ctx)

	// Serialize per-channel number assignment: concurrent inserts must not
	// read the same MAX(number).
	_, err = tx.Exec(ctx, "SELECT pg_advisory_xact_lock(hashtext($1))", input.ChannelID.String())
	if err != nil {
		return model.Nil, err
	}

	query := `
INSERT INTO channels_quotes (channel_id, number, text, creator_id, creator_name, game_id, game_name)
SELECT $1, COALESCE(MAX(number), 0) + 1, $2, $3, $4, $5, $6
FROM channels_quotes
WHERE channel_id = $1
RETURNING ` + quoteFields

	rows, err := tx.Query(
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

	if err := tx.Commit(ctx); err != nil {
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
	updateBuilder = updateBuilder.Set("updated_at", squirrel.Expr("now()"))
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
	query := `DELETE FROM channels_quotes WHERE id = $1`

	rows, err := c.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if rows.RowsAffected() != 1 {
		return quotes.ErrQuoteNotFound
	}

	return nil
}

func (c *Pgx) DeleteByChannelIDAndNumber(ctx context.Context, channelID uuid.UUID, number int) error {
	query := `DELETE FROM channels_quotes WHERE channel_id = $1 AND number = $2`

	rows, err := c.pool.Exec(ctx, query, channelID, number)
	if err != nil {
		return err
	}

	if rows.RowsAffected() != 1 {
		return quotes.ErrQuoteNotFound
	}

	return nil
}
