package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	obsentity "github.com/twirapp/twir/libs/entities/obs"
	channelsmodulesobswebsocket "github.com/twirapp/twir/libs/repositories/channels_modules_obs_websocket"
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

type model struct {
	ID             int
	ChannelID      string
	ServerPort     int
	ServerAddress  string
	ServerPassword string
	Scenes         []string
	Sources        []string
	AudioSources   []string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func modelToEntity(m model) obsentity.ObsWebsocketData {
	return obsentity.ObsWebsocketData{
		ID:             m.ID,
		ChannelID:      m.ChannelID,
		ServerPort:     m.ServerPort,
		ServerAddress:  m.ServerAddress,
		ServerPassword: m.ServerPassword,
		Scenes:         m.Scenes,
		Sources:        m.Sources,
		AudioSources:   m.AudioSources,
	}
}

var _ channelsmodulesobswebsocket.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetByChannelID(ctx context.Context, channelID string) (
	obsentity.ObsWebsocketData,
	error,
) {
	query := `
SELECT id, channel_id, server_port, server_address, server_password, scenes, sources, audio_sources, created_at, updated_at
FROM channels_modules_obs_websocket
WHERE channel_id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, channelID)
	if err != nil {
		return obsentity.NilObsWebsocket, err
	}

	result, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return obsentity.NilObsWebsocket, nil
		}

		return obsentity.NilObsWebsocket, err
	}

	return modelToEntity(result), nil
}

func (c *Pgx) Update(
	ctx context.Context,
	id int,
	input channelsmodulesobswebsocket.UpdateInput,
) error {
	updateBuilder := sq.Update("channels_modules_obs_websocket").Where(squirrel.Eq{"id": id})

	if input.ServerPort != nil {
		updateBuilder = updateBuilder.Set("server_port", *input.ServerPort)
	}

	if input.ServerAddress != nil {
		updateBuilder = updateBuilder.Set("server_address", *input.ServerAddress)
	}

	if input.ServerPassword != nil {
		updateBuilder = updateBuilder.Set("server_password", *input.ServerPassword)
	}

	updateBuilder = updateBuilder.Set("updated_at", squirrel.Expr("NOW()"))

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err = conn.Exec(ctx, query, args...)
	return err
}

func (c *Pgx) Create(
	ctx context.Context,
	input channelsmodulesobswebsocket.CreateInput,
) (obsentity.ObsWebsocketData, error) {
	query := `
INSERT INTO channels_modules_obs_websocket (channel_id, server_port, server_address, server_password)
VALUES ($1, $2, $3, $4)
ON CONFLICT (channel_id)
DO UPDATE SET server_port = EXCLUDED.server_port, server_address = EXCLUDED.server_address, server_password = EXCLUDED.server_password, updated_at = NOW()
RETURNING id, channel_id, server_port, server_address, server_password, scenes, sources, audio_sources, created_at, updated_at
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(
		ctx,
		query,
		input.ChannelID,
		input.ServerPort,
		input.ServerAddress,
		input.ServerPassword,
	)
	if err != nil {
		return obsentity.NilObsWebsocket, err
	}

	result, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model],
	)
	if err != nil {
		return obsentity.NilObsWebsocket, err
	}

	return modelToEntity(result), nil
}

func (c *Pgx) UpdateSources(
	ctx context.Context,
	channelID string,
	input channelsmodulesobswebsocket.UpdateSourcesInput,
) error {
	query := `
INSERT INTO channels_modules_obs_websocket (channel_id, scenes, sources, audio_sources)
VALUES ($1, $2, $3, $4)
ON CONFLICT (channel_id)
DO UPDATE SET scenes = EXCLUDED.scenes, sources = EXCLUDED.sources, audio_sources = EXCLUDED.audio_sources, updated_at = NOW()
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, channelID, input.Scenes, input.Sources, input.AudioSources)
	return err
}

func (c *Pgx) Delete(ctx context.Context, id int) error {
	query := `
DELETE FROM channels_modules_obs_websocket
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, id)
	return err
}
