package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"github.com/twirapp/twir/libs/repositories/chat_wall"
	"github.com/twirapp/twir/libs/repositories/chat_wall/model"
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

var _ chat_wall.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetByID(ctx context.Context, id ulid.ULID) (model.ChatWall, error) {
	query := `
SELECT
	id,
	channel_id,
	created_at,
	updated_at,
	phrase,
	enabled,
	action,
	duration_seconds,
	timeout_duration_seconds,
	(SELECT COUNT(*) FROM channels_chat_wall_log WHERE wall_id = channels_chat_wall.id) AS affected_messages
FROM channels_chat_wall
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, id.String())
	if err != nil {
		return model.ChatWallNil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.ChatWall])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ChatWallNil, nil
		}
		return model.ChatWallNil, err
	}

	return result, nil
}

func (c *Pgx) UpdateChannelSettings(
	ctx context.Context,
	input chat_wall.UpdateChannelSettingsInput,
) error {
	builder := sq.Insert("channels_chat_wall_settings").
		Columns("channel_id", "mute_subscribers", "mute_vips").
		Values(input.ChannelID, input.MuteSubscribers, input.MuteVips).
		Suffix("ON CONFLICT (channel_id) DO UPDATE SET mute_subscribers = EXCLUDED.mute_subscribers, mute_vips = EXCLUDED.mute_vips")

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute insert: %w", err)
	}

	return nil
}

func (c *Pgx) CreateManyLogs(ctx context.Context, inputs []chat_wall.CreateLogInput) error {
	insertBuilder := sq.Insert("channels_chat_wall_log").
		Columns("wall_id", "user_id", "text")

	for _, input := range inputs {
		insertBuilder = insertBuilder.Values(
			input.WallID.String(),
			input.UserID,
			input.Text,
		)
	}

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute insert: %w", err)
	}

	return nil
}

func (c *Pgx) CreateLog(ctx context.Context, input chat_wall.CreateLogInput) error {
	query := `
INSERT INTO channels_chat_wall_log (wall_id, user_id, text)
VALUES ($1, $2, $3)
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, input.WallID.String(), input.UserID, input.Text)
	return err
}

func (c *Pgx) GetChannelSettings(ctx context.Context, channelID string) (
	model.ChatWallSettings,
	error,
) {
	query := `
SELECT id, channel_id, created_at, updated_at, mute_subscribers, mute_vips
FROM channels_chat_wall_settings
WHERE channel_id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, channelID)
	if err != nil {
		return model.ChatWallSettingsNil, err
	}

	result, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model.ChatWallSettings],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ChatWallSettingsNil, chat_wall.ErrSettingsNotFound
		}
		return model.ChatWallSettingsNil, err
	}

	return result, nil
}

func (c *Pgx) GetMany(ctx context.Context, input chat_wall.GetManyInput) ([]model.ChatWall, error) {
	query := `
SELECT
	id,
	channel_id,
	created_at,
	updated_at,
	phrase,
	enabled,
	action,
	duration_seconds,
	timeout_duration_seconds,
	(SELECT COUNT(*) FROM channels_chat_wall_log WHERE wall_id = channels_chat_wall.id) AS affected_messages
FROM channels_chat_wall
WHERE channel_id = $1
ORDER BY created_at DESC
`

	queryArgs := []any{input.ChannelID}

	if input.Enabled != nil {
		query += " AND enabled = $2"
		queryArgs = append(queryArgs, input.Enabled)
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.ChatWall])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (c *Pgx) GetLogs(ctx context.Context, wallID ulid.ULID) ([]model.ChatWallLog, error) {
	query := `
SELECT id, wall_id, user_id, text, created_at
FROM channels_chat_wall_log
WHERE wall_id = $1
ORDER BY created_at DESC
`
	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, wallID.String())
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.ChatWallLog])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (c *Pgx) Create(ctx context.Context, input chat_wall.CreateInput) (model.ChatWall, error) {
	query := `
INSERT INTO channels_chat_wall (channel_id, phrase, enabled, action, duration_seconds, timeout_duration_seconds)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, channel_id, created_at, updated_at, phrase, enabled, action, duration_seconds, timeout_duration_seconds, 0 as affected_messages
`

	var timeoutDuration *int
	if input.TimeoutDuration != nil {
		newDurationFloat := input.TimeoutDuration.Seconds()
		newDurationInt := int(newDurationFloat)
		timeoutDuration = &newDurationInt
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(
		ctx,
		query,
		input.ChannelID,
		input.Phrase,
		input.Enabled,
		input.Action,
		input.Duration.Seconds(),
		timeoutDuration,
	)

	if err != nil {
		return model.ChatWall{}, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.ChatWall])
	if err != nil {
		return model.ChatWall{}, err
	}

	return result, nil
}

func (c *Pgx) Update(ctx context.Context, id ulid.ULID, input chat_wall.UpdateInput) (
	model.ChatWall,
	error,
) {
	queryBuilder := sq.Update("channels_chat_wall").
		Where(squirrel.Eq{"id": id.String()}).
		Suffix("RETURNING id, channel_id, created_at, updated_at, phrase, enabled, action, duration_seconds, timeout_duration_seconds, 0 as affected_messages")

	if input.Phrase != nil {
		queryBuilder = queryBuilder.Set("phrase", *input.Phrase)
	}

	if input.Enabled != nil {
		queryBuilder = queryBuilder.Set("enabled", *input.Enabled)
	}

	if input.Action != nil {
		queryBuilder = queryBuilder.Set("action", *input.Action)
	}

	if input.Duration != nil {
		queryBuilder = queryBuilder.Set("duration_seconds", *input.Duration)
	}

	if input.TimeoutDuration != nil {
		newDurationFloat := input.TimeoutDuration.Seconds()
		newDurationInt := int(newDurationFloat)

		queryBuilder = queryBuilder.Set("timeout_duration_seconds", newDurationInt)
	}

	queryBuilder = queryBuilder.Set("updated_at", squirrel.Expr("NOW()"))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return model.ChatWall{}, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return model.ChatWall{}, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.ChatWall])
	if err != nil {
		return model.ChatWall{}, err
	}

	return result, nil
}

func (c *Pgx) Delete(ctx context.Context, id ulid.ULID) error {
	query := `
DELETE FROM channels_chat_wall
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, id.String())
	return err
}
