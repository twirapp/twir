package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/channels_files"
	"github.com/twirapp/twir/libs/repositories/channels_files/model"
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

var _ channels_files.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (model.ChannelFile, error) {
	query := `
SELECT id, channel_id, mime_type, file_name, size
FROM channels_files
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, id)
	if err != nil {
		return model.ChannelFileNil, fmt.Errorf("getById query: %w", err)
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.ChannelFile])
	if err != nil {
		return model.ChannelFileNil, fmt.Errorf("getById scan: %w", err)
	}

	return result, nil
}

func (c *Pgx) GetMany(ctx context.Context, input channels_files.GetManyInput) (
	[]model.ChannelFile,
	error,
) {
	query := `
SELECT id, channel_id, mime_type, file_name, size
FROM channels_files
WHERE channel_id = $1
ORDER BY size DESC;
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, input.ChannelID)
	if err != nil {
		return nil, fmt.Errorf("getMany query: %w", err)
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.ChannelFile])
	if err != nil {
		return nil, fmt.Errorf("getMany scan: %w", err)
	}

	return result, nil
}

func (c *Pgx) Create(ctx context.Context, input channels_files.CreateInput) (
	model.ChannelFile,
	error,
) {
	query := `
INSERT INTO channels_files (channel_id, file_name, mime_type, size)
VALUES ($1, $2, $3, $4)
RETURNING id, channel_id, file_name, mime_type, size
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(
		ctx,
		query,
		input.ChannelID,
		input.FileName,
		input.MimeType,
		input.Size,
	)
	if err != nil {
		return model.ChannelFile{}, fmt.Errorf("create query: %w", err)
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.ChannelFile])
	if err != nil {
		return model.ChannelFile{}, fmt.Errorf("create scan: %w", err)
	}

	return result, nil
}

func (c *Pgx) DeleteByID(ctx context.Context, id uuid.UUID) error {
	query := `
DELETE FROM channels_files
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, id)
	return err
}

func (c *Pgx) GetTotalChannelUploadedSizeBytes(ctx context.Context, channelID string) (
	int64,
	error,
) {
	query := `
SELECT COALESCE(SUM(size), 0)
FROM channels_files
WHERE channel_id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	var total int64
	err := conn.QueryRow(ctx, query, channelID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("getTotalChannelUploadedSizeBytes: %w", err)
	}

	return total, nil
}
