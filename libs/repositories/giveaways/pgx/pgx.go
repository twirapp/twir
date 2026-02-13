package pgx

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	channels_giveaways "github.com/twirapp/twir/libs/entities/channels_giveaways"
	"github.com/twirapp/twir/libs/repositories"
	"github.com/twirapp/twir/libs/repositories/giveaways"
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

var (
	_  giveaways.Repository = (*Pgx)(nil)
	sq                      = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

// scanModel is internal struct for scanning from database
type scanModel struct {
	ID                   uuid.UUID                       `db:"id"`
	ChannelID            string                          `db:"channel_id"`
	Type                 channels_giveaways.GiveawayType `db:"type"`
	CreatedAt            time.Time                       `db:"created_at"`
	Keyword              *string                         `db:"keyword"`
	MinWatchedTime       *int64                          `db:"min_watched_time"`
	MinMessages          *int32                          `db:"min_messages"`
	MinUsedChannelPoints *int64                          `db:"min_used_channel_points"`
	MinFollowDuration    *int64                          `db:"min_follow_duration"`
	RequireSubscription  bool                            `db:"require_subscription"`
	UpdatedAt            time.Time                       `db:"updated_at"`
	StartedAt            *time.Time                      `db:"started_at"`
	StoppedAt            *time.Time                      `db:"stopped_at"`
	CreatedByUserID      string                          `db:"created_by_user_id"`
}

// scanModelToEntity converts scanModel to entity
func scanModelToEntity(sm scanModel) channels_giveaways.Giveaway {
	return channels_giveaways.Giveaway{
		ID:                   sm.ID,
		ChannelID:            sm.ChannelID,
		Type:                 sm.Type,
		Keyword:              sm.Keyword,
		MinWatchedTime:       sm.MinWatchedTime,
		MinMessages:          sm.MinMessages,
		MinUsedChannelPoints: sm.MinUsedChannelPoints,
		MinFollowDuration:    sm.MinFollowDuration,
		RequireSubscription:  sm.RequireSubscription,
		CreatedAt:            sm.CreatedAt,
		UpdatedAt:            sm.UpdatedAt,
		StartedAt:            sm.StartedAt,
		StoppedAt:            sm.StoppedAt,
		CreatedByUserID:      sm.CreatedByUserID,
	}
}

func (p *Pgx) Create(
	ctx context.Context,
	input giveaways.CreateInput,
) (channels_giveaways.Giveaway, error) {
	query := `
INSERT INTO channels_giveaways (
	"channel_id",
	"type",
	"keyword",
	"min_watched_time",
	"min_messages",
	"min_used_channel_points",
	"min_follow_duration",
	"require_subscription",
	"created_by_user_id"
) VALUES (
	$1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING
	id,
	channel_id,
	type,
	keyword,
	min_watched_time,
	min_messages,
	min_used_channel_points,
	min_follow_duration,
	require_subscription,
	created_at,
	updated_at,
	started_at,
	stopped_at,
	created_by_user_id
	`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(
		ctx,
		query,
		input.ChannelID,
		input.Type,
		input.Keyword,
		input.MinWatchedTime,
		input.MinMessages,
		input.MinUsedChannelPoints,
		input.MinFollowDuration,
		input.RequireSubscription,
		input.CreatedByUserID,
	)
	if err != nil {
		return channels_giveaways.GiveawayNil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		return channels_giveaways.GiveawayNil, err
	}

	return scanModelToEntity(result), nil
}

func (p *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
DELETE FROM channels_giveaways WHERE id = $1
	`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Exec(ctx, query, id.String())
	if err != nil {
		return err
	}

	if rows.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row to be affected, got %d", rows.RowsAffected())
	}

	return nil
}

func (p *Pgx) GetByChannelIDAndKeyword(
	ctx context.Context,
	channelID, keyword string,
) (channels_giveaways.Giveaway, error) {
	query := `
SELECT
	id,
	channel_id,
	type,
	keyword,
	min_watched_time,
	min_messages,
	min_used_channel_points,
	min_follow_duration,
	require_subscription,
	created_at,
	updated_at,
	started_at,
	stopped_at,
	created_by_user_id
FROM channels_giveaways
WHERE channel_id = $1 AND keyword = $2 AND stopped_at IS NULL
ORDER BY created_at DESC;
	`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, channelID, keyword)
	if err != nil {
		return channels_giveaways.GiveawayNil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return channels_giveaways.GiveawayNil, giveaways.ErrNotFound
		}

		return channels_giveaways.GiveawayNil, err
	}

	return scanModelToEntity(result), nil
}

func (p *Pgx) GetByID(ctx context.Context, id uuid.UUID) (channels_giveaways.Giveaway, error) {
	query := `
SELECT
	id,
	channel_id,
	type,
	keyword,
	min_watched_time,
	min_messages,
	min_used_channel_points,
	min_follow_duration,
	require_subscription,
	created_at,
	updated_at,
	started_at,
	stopped_at,
	created_by_user_id
FROM channels_giveaways
WHERE id = $1
LIMIT 1;
	`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, id.String())
	if err != nil {
		return channels_giveaways.GiveawayNil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return channels_giveaways.GiveawayNil, giveaways.ErrNotFound
		}

		return channels_giveaways.GiveawayNil, err
	}

	return scanModelToEntity(result), nil
}

func (p *Pgx) GetManyByChannelID(
	ctx context.Context,
	channelID string,
) ([]channels_giveaways.Giveaway, error) {
	selectBuilder := sq.Select(
		"id",
		"channel_id",
		"type",
		"keyword",
		"min_watched_time",
		"min_messages",
		"min_used_channel_points",
		"min_follow_duration",
		"require_subscription",
		"created_at",
		"updated_at",
		"started_at",
		"stopped_at",
		"created_by_user_id",
	).
		From("channels_giveaways").
		OrderBy("created_at DESC").
		Where(squirrel.Eq{`"channel_id"`: channelID})

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	scanResults, err := pgx.CollectRows(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		return nil, err
	}

	result := make([]channels_giveaways.Giveaway, 0, len(scanResults))
	for _, sm := range scanResults {
		result = append(result, scanModelToEntity(sm))
	}

	return result, nil
}

func (p *Pgx) GetManyActiveByChannelID(
	ctx context.Context,
	channelID string,
) ([]channels_giveaways.Giveaway, error) {
	selectBuilder := sq.Select(
		"id",
		"channel_id",
		"type",
		"keyword",
		"min_watched_time",
		"min_messages",
		"min_used_channel_points",
		"min_follow_duration",
		"require_subscription",
		"created_at",
		"updated_at",
		"started_at",
		"stopped_at",
		"created_by_user_id",
	).
		From("channels_giveaways").
		OrderBy("created_at DESC").
		Where(squirrel.Eq{`"channel_id"`: channelID}).
		Where(squirrel.Expr("stopped_at IS NULL"))

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	scanResults, err := pgx.CollectRows(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		return nil, err
	}

	result := make([]channels_giveaways.Giveaway, 0, len(scanResults))
	for _, sm := range scanResults {
		result = append(result, scanModelToEntity(sm))
	}

	return result, nil
}

func (p *Pgx) UpdateStatuses(
	ctx context.Context,
	id uuid.UUID,
	input giveaways.UpdateStatusInput,
) (channels_giveaways.Giveaway, error) {
	updateBuilder := sq.Update("channels_giveaways").
		Where(squirrel.Eq{"id": id.String()}).
		Suffix(
			`RETURNING
			id,
			channel_id,
			type,
			keyword,
			min_watched_time,
			min_messages,
			min_used_channel_points,
			min_follow_duration,
			require_subscription,
			created_at,
			updated_at,
			started_at,
			stopped_at,
			created_by_user_id`,
		)

	if input.StartedAt.Valid {
		updateBuilder = updateBuilder.Set("started_at", input.StartedAt)
	} else {
		updateBuilder = updateBuilder.Set("started_at", nil)
	}

	if input.StoppedAt.Valid {
		updateBuilder = updateBuilder.Set("stopped_at", input.StoppedAt)
	} else {
		updateBuilder = updateBuilder.Set("stopped_at", nil)
	}

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return channels_giveaways.GiveawayNil, err
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return channels_giveaways.GiveawayNil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		return channels_giveaways.GiveawayNil, err
	}

	return scanModelToEntity(result), nil
}

func (p *Pgx) Update(
	ctx context.Context,
	id uuid.UUID,
	input giveaways.UpdateInput,
) (channels_giveaways.Giveaway, error) {
	updateBuilder := sq.Update("channels_giveaways").
		Where(squirrel.Eq{"id": id.String()}).
		Suffix(
			`RETURNING
			id,
			channel_id,
			type,
			keyword,
			min_watched_time,
			min_messages,
			min_used_channel_points,
			min_follow_duration,
			require_subscription,
			created_at,
			updated_at,
			started_at,
			stopped_at,
			created_by_user_id`,
		)
	updateBuilder = repositories.SquirrelApplyPatch(
		updateBuilder,
		map[string]any{
			"started_at":              input.StartedAt,
			"keyword":                 input.Keyword,
			"stopped_at":              input.StoppedAt,
			"min_watched_time":        input.MinWatchedTime,
			"min_messages":            input.MinMessages,
			"min_used_channel_points": input.MinUsedChannelPoints,
			"min_follow_duration":     input.MinFollowDuration,
			"require_subscription":    input.RequireSubscription,
		},
	)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return channels_giveaways.GiveawayNil, err
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return channels_giveaways.GiveawayNil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		return channels_giveaways.GiveawayNil, err
	}

	return scanModelToEntity(result), nil
}
