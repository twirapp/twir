package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/twirapp/twir/libs/entities/platform"

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

const selectColumns = `id, channel_id, "userId", "userLogin", "userName", "gameId", "gameName", "communityIds", type, title, "viewerCount", "startedAt", "language", "thumbnailUrl", "tagIds", tags, "isMature", platform`

func (c *Pgx) GetByChannelID(
	ctx context.Context,
	channelID uuid.UUID,
	platform platform.Platform,
) (model.Stream, error) {
	query := `
SELECT ` + selectColumns + `
FROM channels_streams
WHERE channel_id = $1 AND platform = $2
LIMIT 1
`

	rows, err := c.getter.DefaultTrOrDB(ctx, c.pool).Query(ctx, query, channelID, platform)
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

func (c *Pgx) GetByUserID(
	ctx context.Context,
	userID string,
	platform platform.Platform,
) (model.Stream, error) {
	query := `
SELECT ` + selectColumns + `
FROM channels_streams
WHERE "userId" = $1 AND platform = $2
LIMIT 1
`

	rows, err := c.getter.DefaultTrOrDB(ctx, c.pool).Query(ctx, query, userID, platform)
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

func (c *Pgx) GetListByChannelID(ctx context.Context, channelID uuid.UUID) ([]model.Stream, error) {
	query := `
SELECT ` + selectColumns + `
FROM channels_streams
WHERE channel_id = $1
`

	rows, err := c.getter.DefaultTrOrDB(ctx, c.pool).Query(ctx, query, channelID)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Stream])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Pgx) GetList(ctx context.Context) ([]model.Stream, error) {
	query := `
SELECT ` + selectColumns + `
FROM channels_streams
`

	rows, err := c.getter.DefaultTrOrDB(ctx, c.pool).Query(ctx, query)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Stream])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Pgx) Save(ctx context.Context, input streams.SaveInput) error {
	query := `
INSERT INTO channels_streams (
	id, channel_id, "userId", "userLogin", "userName", "gameId", "gameName", "communityIds", type,
	title, "viewerCount", "startedAt", language, "thumbnailUrl", "tagIds", tags, "isMature", platform
) VALUES (
		COALESCE(NULLIF($1, ''), uuidv7()::text), $2, $3, $4, $5, $6, $7, $8, $9,
	$10, $11, $12, $13, $14, $15, $16, $17, $18
)
ON CONFLICT ("userId", platform) DO UPDATE SET
	channel_id = EXCLUDED.channel_id,
	"userLogin" = EXCLUDED."userLogin",
	"userName" = EXCLUDED."userName",
	"gameId" = EXCLUDED."gameId",
	"gameName" = EXCLUDED."gameName",
	"communityIds" = EXCLUDED."communityIds",
	type = EXCLUDED.type,
	title = EXCLUDED.title,
	"viewerCount" = EXCLUDED."viewerCount",
	"startedAt" = EXCLUDED."startedAt",
	language = EXCLUDED.language,
	"thumbnailUrl" = EXCLUDED."thumbnailUrl",
	"tagIds" = EXCLUDED."tagIds",
	tags = EXCLUDED.tags,
	"isMature" = EXCLUDED."isMature"
`

	_, err := c.getter.DefaultTrOrDB(ctx, c.pool).Exec(
		ctx,
		query,
		input.ID,
		input.ChannelID,
		input.UserId,
		input.UserLogin,
		input.UserName,
		input.GameId,
		input.GameName,
		pq.Array(input.CommunityIds),
		input.Type,
		input.Title,
		input.ViewerCount,
		input.StartedAt,
		input.Language,
		input.ThumbnailUrl,
		pq.Array(input.TagIds),
		pq.Array(input.Tags),
		input.IsMature,
		input.Platform,
	)
	return err
}

func (c *Pgx) DeleteByChannelID(ctx context.Context, channelID uuid.UUID, platform platform.Platform) error {
	query := `DELETE FROM channels_streams WHERE channel_id = $1 AND platform = $2`
	_, err := c.getter.DefaultTrOrDB(ctx, c.pool).Exec(ctx, query, channelID, platform)
	return err
}

func (c *Pgx) Update(
	ctx context.Context,
	channelID uuid.UUID,
	platform platform.Platform,
	input streams.UpdateInput,
) error {
	updateBuilder := sq.Update("channels_streams").
		Where(squirrel.Eq{`channel_id`: channelID, `platform`: platform.String()})

	updateBuilder = repositories.SquirrelApplyPatch(
		updateBuilder,
		map[string]any{
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

	_, err = c.getter.DefaultTrOrDB(ctx, c.pool).Exec(ctx, query, args...)
	return err
}

func (c *Pgx) Count(ctx context.Context) (uint64, error) {
	query := `SELECT COUNT(*) FROM channels_streams`

	var count uint64
	err := c.getter.DefaultTrOrDB(ctx, c.pool).QueryRow(ctx, query).Scan(&count)
	return count, err
}
