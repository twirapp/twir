package pgx

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/channels_storage"
	"github.com/twirapp/twir/libs/repositories/channels_storage/model"
)

type Opts struct {
	Pgx *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{
		pool: opts.Pgx,
	}
}

func NewFx(pgxpool *pgxpool.Pool) *Pgx {
	return New(
		Opts{
			Pgx: pgxpool,
		},
	)
}

var (
	_ channels_storage.Repository = (*Pgx)(nil)
)

type Pgx struct {
	pool *pgxpool.Pool
}

func (c *Pgx) GetAllByChannelID(ctx context.Context, channelID string) (
	[]model.ChannelStorage,
	error,
) {
	query := `
SELECT id, channel_id, key, value, created_at, updated_at
FROM channels_storage
WHERE channel_id = $1
ORDER BY key
`
	rows, err := c.pool.Query(ctx, query, channelID)
	if err != nil {
		return nil, err
	}

	entries, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.ChannelStorage])
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (c *Pgx) GetByKey(ctx context.Context, channelID string, key string) (
	model.ChannelStorage,
	error,
) {
	query := `
SELECT id, channel_id, key, value, created_at, updated_at
FROM channels_storage
WHERE channel_id = $1 AND key = $2
LIMIT 1
`
	rows, err := c.pool.Query(ctx, query, channelID, key)
	if err != nil {
		return model.Nil, err
	}

	entry, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.ChannelStorage])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, channels_storage.ErrNotFound
		}
		return model.Nil, err
	}

	return entry, nil
}

func (c *Pgx) Set(ctx context.Context, input channels_storage.SetInput) (
	model.ChannelStorage,
	error,
) {
	query := `
INSERT INTO channels_storage (channel_id, key, value)
VALUES ($1, $2, $3)
ON CONFLICT (channel_id, key) DO UPDATE SET value = EXCLUDED.value, updated_at = now()
RETURNING id, channel_id, key, value, created_at, updated_at
`
	rows, err := c.pool.Query(ctx, query, input.ChannelID, input.Key, input.Value)
	if err != nil {
		return model.Nil, err
	}

	entry, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.ChannelStorage])
	if err != nil {
		return model.Nil, err
	}

	return entry, nil
}

func (c *Pgx) Delete(ctx context.Context, channelID string, key string) error {
	query := `
DELETE FROM channels_storage
WHERE channel_id = $1 AND key = $2
`
	tag, err := c.pool.Exec(ctx, query, channelID, key)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return channels_storage.ErrNotFound
	}

	return nil
}

func (c *Pgx) DeleteAllByChannelID(ctx context.Context, channelID string) error {
	query := `
DELETE FROM channels_storage
WHERE channel_id = $1
`
	_, err := c.pool.Exec(ctx, query, channelID)
	return err
}

func (c *Pgx) GetTotalSizeByChannelID(ctx context.Context, channelID string) (int64, error) {
	query := `
SELECT COALESCE(SUM(pg_column_size(value)), 0)
FROM channels_storage
WHERE channel_id = $1
`
	var size int64
	err := c.pool.QueryRow(ctx, query, channelID).Scan(&size)
	if err != nil {
		return 0, err
	}

	return size, nil
}

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (model.ChannelStorage, error) {
	query := `
SELECT id, channel_id, key, value, created_at, updated_at
FROM channels_storage
WHERE id = $1
LIMIT 1
`
	rows, err := c.pool.Query(ctx, query, id)
	if err != nil {
		return model.Nil, err
	}

	entry, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.ChannelStorage])
	if err != nil {
		return model.Nil, err
	}

	return entry, nil
}

func (c *Pgx) Update(ctx context.Context, id uuid.UUID, value json.RawMessage) (
	model.ChannelStorage,
	error,
) {
	query := `
UPDATE channels_storage
SET value = $2, updated_at = now()
WHERE id = $1
RETURNING id, channel_id, key, value, created_at, updated_at
`
	rows, err := c.pool.Query(ctx, query, id, value)
	if err != nil {
		return model.Nil, err
	}

	entry, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.ChannelStorage])
	if err != nil {
		return model.Nil, err
	}

	return entry, nil
}
