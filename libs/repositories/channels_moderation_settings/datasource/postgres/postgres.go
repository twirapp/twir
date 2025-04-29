package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/channels_moderation_settings"
	"github.com/twirapp/twir/libs/repositories/channels_moderation_settings/model"
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

var _ channels_moderation_settings.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

var selectColumns = []string{
	"id",
	"name",
	"channel_id",
	"type",
	"enabled",
	"max_warnings",
	"ban_time",
	"ban_message",
	"warning_message",
	"check_clips",
	"trigger_length",
	"max_percentage",
	"deny_list",
	"excluded_roles",
	"denied_chat_languages",
	"deny_list_regexp_enabled",
	"deny_list_word_boundary_enabled",
	"deny_list_sensitivity_enabled",
	"created_at",
	"updated_at",
}

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (model.ChannelModerationSettings, error) {
	query, args, err := sq.Select(selectColumns...).
		From("channels_moderation_settings").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return model.ChannelModerationSettings{}, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return model.ChannelModerationSettings{}, err
	}

	result, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model.ChannelModerationSettings],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ChannelModerationSettings{}, channels_moderation_settings.ErrNotFound
		}
		return model.ChannelModerationSettings{}, err
	}

	return result, nil
}

func (c *Pgx) Create(
	ctx context.Context,
	input channels_moderation_settings.CreateOrUpdateInput,
) (model.ChannelModerationSettings, error) {
	query, args, err := sq.Insert("channels_moderation_settings").
		SetMap(makeCreateOrUpdateMap(input)).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return model.ChannelModerationSettings{}, err
	}

	var id uuid.UUID
	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	if err := conn.QueryRow(ctx, query, args...).Scan(&id); err != nil {
		return model.ChannelModerationSettings{}, err
	}

	return c.GetByID(ctx, id)
}

func (c *Pgx) Update(
	ctx context.Context,
	id uuid.UUID,
	input channels_moderation_settings.CreateOrUpdateInput,
) (model.ChannelModerationSettings, error) {
	updateInput := makeCreateOrUpdateMap(input)
	updateInput["updated_at"] = time.Now()

	query, args, err := sq.Update("channels_moderation_settings").
		SetMap(updateInput).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return model.ChannelModerationSettings{}, err
	}

	var idResult uuid.UUID
	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	if err := conn.QueryRow(ctx, query, args...).Scan(&idResult); err != nil {
		return model.ChannelModerationSettings{}, err
	}

	return c.GetByID(ctx, idResult)
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
DELETE FROM channels_moderation_settings
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, id)
	return err
}

func (c *Pgx) GetByChannelID(
	ctx context.Context,
	channelID string,
) ([]model.ChannelModerationSettings, error) {
	query, args, err := sq.Select(selectColumns...).
		From("channels_moderation_settings").
		Where(squirrel.Eq{"channel_id": channelID}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.ChannelModerationSettings])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}
