package pgx

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/alerts"
	"github.com/twirapp/twir/libs/repositories/alerts/model"
)

type Opts struct {
	PgxPool *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{
		pool: opts.PgxPool,
	}
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return New(Opts{PgxPool: pool})
}

var _ alerts.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool *pgxpool.Pool
}

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (model.Alert, error) {
	query := `
SELECT id, channel_id, audio_id, name, audio_volume, command_ids, reward_ids, greetings_ids, keywords_ids
FROM channels_alerts
WHERE id = $1
LIMIT 1
`

	rows, err := c.pool.Query(ctx, query, id)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Alert])
	if err != nil {
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) GetManyByChannelID(ctx context.Context, channelID string) ([]model.Alert, error) {
	query := `
SELECT id, channel_id, audio_id, name, audio_volume, command_ids, reward_ids, greetings_ids, keywords_ids
FROM channels_alerts
WHERE channel_id = $1
`

	rows, err := c.pool.Query(ctx, query, channelID)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Alert])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Pgx) Create(ctx context.Context, input alerts.CreateInput) (model.Alert, error) {
	query := `
INSERT INTO channels_alerts (channel_id, audio_id, name, audio_volume, command_ids, reward_ids, greetings_ids, keywords_ids)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, channel_id, audio_id, name, audio_volume, command_ids, reward_ids, greetings_ids, keywords_ids
`

	rows, err := c.pool.Query(
		ctx,
		query,
		input.ChannelID,
		input.AudioID,
		input.Name,
		input.AudioVolume,
		input.CommandIDS,
		input.RewardIDS,
		input.GreetingsIDS,
		input.KeywordsIDS,
	)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Alert])
	if err != nil {
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) Update(ctx context.Context, id uuid.UUID, input alerts.UpdateInput) (
	model.Alert,
	error,
) {
	updateBuilder := sq.
		Update("channels_alerts").
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id, channel_id, audio_id, name, audio_volume, command_ids, reward_ids, greetings_ids, keywords_ids")

	if input.Name != nil {
		updateBuilder = updateBuilder.Set("name", *input.Name)
	}

	if input.AudioID != nil {
		updateBuilder = updateBuilder.Set("audio_id", *input.AudioID)
	}

	if input.AudioVolume != nil {
		updateBuilder = updateBuilder.Set("audio_volume", *input.AudioVolume)
	}

	if input.CommandIDS != nil {
		updateBuilder = updateBuilder.Set("command_ids", input.CommandIDS)
	}

	if input.RewardIDS != nil {
		updateBuilder = updateBuilder.Set("reward_ids", input.RewardIDS)
	}

	if input.GreetingsIDS != nil {
		updateBuilder = updateBuilder.Set("greetings_ids", input.GreetingsIDS)
	}

	if input.KeywordsIDS != nil {
		updateBuilder = updateBuilder.Set("keywords_ids", input.KeywordsIDS)
	}

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return model.Nil, err
	}

	rows, err := c.pool.Query(ctx, query, args...)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Alert])
	if err != nil {
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
DELETE FROM channels_alerts
WHERE id = $1
`

	result, err := c.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() != 1 {
		return alerts.ErrNotFound
	}

	return nil
}
