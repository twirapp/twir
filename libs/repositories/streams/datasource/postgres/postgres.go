package postgres

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories"
	"github.com/twirapp/twir/libs/repositories/streams"
	"github.com/twirapp/twir/libs/repositories/streams/model"
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

var _ streams.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetByChannelID(ctx context.Context, channelID string) (model.Stream, error) {
	query := `
SELECT id, "userId", "userLogin", "userName", "gameId", "gameName", "communityIds", type, title, "viewerCount", "startedAt", "language", "thumbnailUrl", "tagIds", tags, "isMature"
FROM channels_streams
WHERE "userId" = $1
LIMIT 1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, channelID)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Stream])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, nil
		}
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) GetList(ctx context.Context) ([]model.Stream, error) {
	query := `
SELECT id, "userId", "userLogin", "userName", "gameId", "gameName", "communityIds", type, title, "viewerCount", "startedAt", "language", "thumbnailUrl", "tagIds", tags, "isMature"
FROM channels_streams
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Stream])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Pgx) Update(ctx context.Context, channelID string, input streams.UpdateInput) error {
	updateBuilder := sq.Update("channels_streams").
		Where(squirrel.Eq{`"userId"`: channelID})

	updateBuilder = repositories.SquirrelApplyPatch(
		updateBuilder, map[string]interface{}{
			`"gameId"`:       input.GameId,
			`"gameName"`:     input.GameName,
			`"communityIds"`: input.CommunityIds,
			`"type"`:         input.Type,
			`"title"`:        input.Title,
			`"viewerCount"`:  input.ViewerCount,
			`"startedAt"`:    input.StartedAt,
			`"language"`:     input.Language,
			`"thumbnailUrl"`: input.ThumbnailUrl,
			`"tagIds"`:       input.TagIds,
			`"tags"`:         input.Tags,
			`"isMature"`:     input.IsMature,
		},
	)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err = conn.Exec(ctx, query, args...)
	return err
}
