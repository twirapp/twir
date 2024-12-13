package pgx

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/keywords"
	"github.com/twirapp/twir/libs/repositories/keywords/model"
)

type Opts struct {
	PgxPool *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{
		pgxpool: opts.PgxPool,
	}
}

func NewFx(pgxpool *pgxpool.Pool) *Pgx {
	return New(Opts{PgxPool: pgxpool})
}

var _ keywords.Repository = (*Pgx)(nil)

type Pgx struct {
	pgxpool *pgxpool.Pool
}

func (c *Pgx) GetAllByChannelID(ctx context.Context, channelID string) ([]model.Keyword, error) {
	query := `
SELECT id, "channelId", text, response, enabled, cooldown, "cooldownExpireAt", "isReply", "isRegular", usages
FROM channels_keywords
WHERE "channelId" = $1
`

	rows, err := c.pgxpool.Query(ctx, query, channelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Keyword])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Pgx) CountByChannelID(ctx context.Context, channelID string) (int, error) {
	query := `
SELECT COUNT(*)
FROM channels_keywords
WHERE "channelId" = $1
`

	var count int
	err := c.pgxpool.QueryRow(ctx, query, channelID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (model.Keyword, error) {
	query := `
SELECT id, "channelId", text, response, enabled, cooldown, "cooldownExpireAt", "isReply", "isRegular", usages
FROM channels_keywords
WHERE id = $1
LIMIT 1;
`

	rows, err := c.pgxpool.Query(ctx, query, id)
	if err != nil {
		return model.Nil, err
	}
	defer rows.Close()

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Keyword])
	if err != nil {
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) Create(ctx context.Context, input keywords.CreateInput) (model.Keyword, error) {
	query := `
INSERT INTO channels_keywords ("channelId", text, response, enabled, cooldown, "cooldownExpireAt", "isReply", "isRegular", usages)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, "channelId", text, response, enabled, cooldown, "cooldownExpireAt", "isReply", "isRegular", usages
`

	rows, err := c.pgxpool.Query(
		ctx,
		query,
		input.ChannelID,
		input.Text,
		input.Response,
		input.Enabled,
		input.Cooldown,
		input.CooldownExpireAt,
		input.IsReply,
		input.IsRegular,
		input.Usages,
	)
	if err != nil {
		return model.Nil, err
	}

	defer rows.Close()

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Keyword])
	if err != nil {
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) Update(ctx context.Context, id uuid.UUID, input keywords.UpdateInput) (
	model.Keyword,
	error,
) {
	// TODO implement me
	panic("implement me")
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
DELETE FROM channels_keywords
WHERE id = $1
`

	rows, err := c.pgxpool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if rows.RowsAffected() != 0 {
		return pgx.ErrNoRows
	}

	return nil
}
