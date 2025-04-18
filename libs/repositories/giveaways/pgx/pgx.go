package pgx

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"github.com/twirapp/twir/libs/repositories"
	"github.com/twirapp/twir/libs/repositories/giveaways"
	"github.com/twirapp/twir/libs/repositories/giveaways/model"
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

var _ giveaways.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (p *Pgx) Create(
	ctx context.Context,
	input giveaways.CreateInput,
) (model.ChannelGiveaway, error) {
	query := `
INSERT INTO channels_giveaways ("channel_id", "keyword", "created_by_user_id") VALUES (
	$1, $2, $3
) RETURNING id, channel_id, keyword, created_at, updated_at, started_at, ended_at, stopped_at, created_by_user_id, archived_at
	`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(
		ctx,
		query,
		input.ChannelID,
		input.Keyword,
		input.CreatedByUserID,
	)
	if err != nil {
		return model.ChannelGiveawayNil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.ChannelGiveaway])
	if err != nil {
		return model.ChannelGiveawayNil, err
	}

	return result, nil
}

func (p *Pgx) Delete(ctx context.Context, id ulid.ULID) error {
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
) (model.ChannelGiveaway, error) {
	query := `
SELECT id, channel_id, created_at, updated_at, started_at, ended_at, stopped_at, keyword, created_by_user_id, archived_at
FROM channels_giveaways
WHERE channel_id = $1 AND keyword = $2 AND archived_at IS NULL
	`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, channelID, keyword)
	if err != nil {
		return model.ChannelGiveawayNil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.ChannelGiveaway])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ChannelGiveawayNil, giveaways.ErrNotFound
		}

		return model.ChannelGiveawayNil, err
	}

	return result, nil
}

func (p *Pgx) GetByID(ctx context.Context, id ulid.ULID) (model.ChannelGiveaway, error) {
	query := `
SELECT id, channel_id, created_at, updated_at, started_at, ended_at, stopped_at, keyword, created_by_user_id, archived_at
FROM channels_giveaways
WHERE id = $1
	`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, id.String())
	if err != nil {
		return model.ChannelGiveawayNil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.ChannelGiveaway])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ChannelGiveawayNil, giveaways.ErrNotFound
		}

		return model.ChannelGiveawayNil, err
	}

	return result, nil
}

func (p *Pgx) GetManyByChannelID(
	ctx context.Context,
	channelID string,
) ([]model.ChannelGiveaway, error) {
	selectBuilder := sq.Select("id", "channel_id", "created_at", "keyword", "updated_at", "started_at", "ended_at", "created_by_user_id", "stopped_at", "archived_at").
		From("channels_giveaways").
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

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.ChannelGiveaway])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *Pgx) GetManyActiveByChannelID(
	ctx context.Context,
	channelID string,
) ([]model.ChannelGiveaway, error) {
	selectBuilder := sq.Select("id", "channel_id", "created_at", "keyword", "updated_at", "started_at", "ended_at", "created_by_user_id", "stopped_at", "archived_at").
		From("channels_giveaways").
		Where(squirrel.Eq{`"channel_id"`: channelID}).
		Where(squirrel.Expr("archived_at IS NULL"))

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.ChannelGiveaway])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *Pgx) UpdateStatuses(
	ctx context.Context,
	id ulid.ULID,
	input giveaways.UpdateStatusInput,
) (model.ChannelGiveaway, error) {
	updateBuilder := sq.Update("channels_giveaways").
		Where(squirrel.Eq{"id": id.String()}).
		Suffix(`RETURNING id, channel_id, created_at, updated_at, started_at, ended_at, keyword, created_by_user_id, archived_at, stopped_at`)

	if input.StartedAt.Valid {
		updateBuilder = updateBuilder.Set("started_at", input.StartedAt)
	} else {
		updateBuilder = updateBuilder.Set("started_at", nil)
	}

	if input.EndedAt.Valid {
		updateBuilder = updateBuilder.Set("ended_at", input.EndedAt)
	} else {
		updateBuilder = updateBuilder.Set("ended_at", nil)
	}

	if input.ArchivedAt.Valid {
		updateBuilder = updateBuilder.Set("archived_at", input.ArchivedAt)
	} else {
		updateBuilder = updateBuilder.Set("archived_at", nil)
	}

	if input.StoppedAt.Valid {
		updateBuilder = updateBuilder.Set("stopped_at", input.StoppedAt)
	} else {
		updateBuilder = updateBuilder.Set("stopped_at", nil)
	}

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return model.ChannelGiveawayNil, err
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return model.ChannelGiveawayNil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.ChannelGiveaway])
	if err != nil {
		return model.ChannelGiveawayNil, err
	}

	return result, nil
}

func (p *Pgx) Update(
	ctx context.Context,
	id ulid.ULID,
	input giveaways.UpdateInput,
) (model.ChannelGiveaway, error) {
	updateBuilder := sq.Update("channels_giveaways").
		Where(squirrel.Eq{"id": id.String()}).
		Suffix(`RETURNING id, channel_id, created_at, updated_at, started_at, ended_at, keyword, created_by_user_id, archived_at, stopped_at`)
	updateBuilder = repositories.SquirrelApplyPatch(
		updateBuilder,
		map[string]any{
			"started_at":  input.StartedAt,
			"ended_at":    input.EndedAt,
			"keyword":     input.Keyword,
			"archived_at": input.ArchivedAt,
			"stopped_at":  input.StoppedAt,
		},
	)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return model.ChannelGiveawayNil, err
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return model.ChannelGiveawayNil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.ChannelGiveaway])
	if err != nil {
		return model.ChannelGiveawayNil, err
	}

	return result, nil
}
