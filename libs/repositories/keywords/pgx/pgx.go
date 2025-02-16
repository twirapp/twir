package pgx

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories"
	"github.com/twirapp/twir/libs/repositories/keywords"
	"github.com/twirapp/twir/libs/repositories/keywords/model"
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

var _ keywords.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool *pgxpool.Pool
}

func (c *Pgx) GetAllByChannelID(ctx context.Context, channelID string) ([]model.Keyword, error) {
	query := `
SELECT id, "channelId", text, response, enabled, cooldown, "cooldownExpireAt", "isReply", "isRegular", usages
FROM channels_keywords
WHERE "channelId" = $1
`

	rows, err := c.pool.Query(ctx, query, channelID)
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
	err := c.pool.QueryRow(ctx, query, channelID).Scan(&count)
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

	rows, err := c.pool.Query(ctx, query, id)
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

	rows, err := c.pool.Query(
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
	updateBuilder := sq.Update("channels_keywords")
	updateBuilder = repositories.SquirrelApplyPatch(
		updateBuilder, map[string]any{
			"text":               input.Text,
			"response":           input.Response,
			"enabled":            input.Enabled,
			"cooldown":           input.Cooldown,
			`"cooldownExpireAt"`: input.CooldownExpireAt,
			`"isReply"`:          input.IsReply,
			`"isRegular"`:        input.IsRegular,
			"usages":             input.Usages,
		},
	)

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
DELETE FROM channels_keywords
WHERE id = $1
`

	rows, err := c.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if rows.RowsAffected() != 1 {
		return keywords.ErrKeywordNotFound
	}

	return nil
}
