package pgx

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/entities/platform"
	channelplatforms "github.com/twirapp/twir/libs/repositories/channel_platforms"
	"github.com/twirapp/twir/libs/repositories/channel_platforms/model"
)

type Opts struct {
	PgxPool *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{pool: opts.PgxPool}
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return New(Opts{PgxPool: pool})
}

var _ channelplatforms.Repository = (*Pgx)(nil)

type Pgx struct {
	pool *pgxpool.Pool
}

const selectColumns = `
	id,
	channel_id,
	platform,
	user_id,
	platform_channel_id,
	enabled,
	bot_user_id,
	bot_config,
	created_at,
	updated_at`

func (r *Pgx) Create(
	ctx context.Context,
	input channelplatforms.CreateInput,
) (model.ChannelPlatform, error) {
	query := `
		INSERT INTO channel_platforms (
			channel_id,
			platform,
			user_id,
			platform_channel_id,
			enabled,
			bot_user_id,
			bot_config
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING ` + selectColumns

	rows, err := r.pool.Query(
		ctx,
		query,
		input.ChannelID,
		input.Platform,
		input.UserID,
		input.PlatformChannelID,
		input.Enabled,
		input.BotUserID,
		input.BotConfig,
	)
	if err != nil {
		return model.Nil, fmt.Errorf("create channel platform binding: %w", err)
	}

	binding, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.ChannelPlatform])
	if err != nil {
		return model.Nil, fmt.Errorf("collect created channel platform binding: %w", err)
	}

	return binding, nil
}

func (r *Pgx) GetByChannelAndPlatform(
	ctx context.Context,
	channelID uuid.UUID,
	platform platform.Platform,
) (model.ChannelPlatform, error) {
	return r.getOne(
		ctx,
		`
			SELECT `+selectColumns+`
			FROM channel_platforms
			WHERE channel_id = $1 AND platform = $2
			LIMIT 1`,
		channelID,
		platform,
	)
}

func (r *Pgx) GetByPlatformChannelID(
	ctx context.Context,
	platform platform.Platform,
	platformChannelID string,
) (model.ChannelPlatform, error) {
	return r.getOne(
		ctx,
		`
			SELECT `+selectColumns+`
			FROM channel_platforms
			WHERE platform = $1 AND platform_channel_id = $2
			LIMIT 1`,
		platform,
		platformChannelID,
	)
}

func (r *Pgx) ListByChannelID(
	ctx context.Context,
	channelID uuid.UUID,
) ([]model.ChannelPlatform, error) {
	query := `
		SELECT ` + selectColumns + `
		FROM channel_platforms
		WHERE channel_id = $1
		ORDER BY platform`

	rows, err := r.pool.Query(ctx, query, channelID)
	if err != nil {
		return nil, fmt.Errorf("list channel platform bindings: %w", err)
	}

	bindings, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.ChannelPlatform])
	if err != nil {
		return nil, fmt.Errorf("collect channel platform bindings: %w", err)
	}

	return bindings, nil
}

func (r *Pgx) Update(
	ctx context.Context,
	id uuid.UUID,
	input channelplatforms.UpdateInput,
) (model.ChannelPlatform, error) {
	query := `
		UPDATE channel_platforms
		SET
			user_id = $2,
			platform_channel_id = $3,
			enabled = $4,
			bot_user_id = $5,
			bot_config = $6,
			updated_at = NOW()
		WHERE id = $1
		RETURNING ` + selectColumns

	rows, err := r.pool.Query(
		ctx,
		query,
		id,
		input.UserID,
		input.PlatformChannelID,
		input.Enabled,
		input.BotUserID,
		input.BotConfig,
	)
	if err != nil {
		return model.Nil, fmt.Errorf("update channel platform binding: %w", err)
	}

	binding, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.ChannelPlatform])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, channelplatforms.ErrNotFound
		}
		return model.Nil, fmt.Errorf("collect updated channel platform binding: %w", err)
	}

	return binding, nil
}

func (r *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	commandTag, err := r.pool.Exec(
		ctx,
		`DELETE FROM channel_platforms WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete channel platform binding: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return channelplatforms.ErrNotFound
	}

	return nil
}

func (r *Pgx) getOne(
	ctx context.Context,
	query string,
	args ...any,
) (model.ChannelPlatform, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return model.Nil, fmt.Errorf("query channel platform binding: %w", err)
	}

	binding, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.ChannelPlatform])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, channelplatforms.ErrNotFound
		}
		return model.Nil, fmt.Errorf("collect channel platform binding: %w", err)
	}

	return binding, nil
}
