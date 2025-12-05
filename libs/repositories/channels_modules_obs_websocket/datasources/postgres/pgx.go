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

func (c *Pgx) Upsert(
	ctx context.Context,
	input channelsmodulesobswebsocket.UpsertInput,
) (obsentity.ObsWebsocketData, error) {
	m, err := c.GetByChannelID(ctx, input.ChannelID)
	if err != nil {
		return obsentity.NilObsWebsocket, err
	}

	setMap := map[string]any{
		"channel_id": input.ChannelID,
	}

	if input.ServerPort != nil {
		setMap["server_port"] = *input.ServerPort
	}

	if input.ServerAddress != nil {
		setMap["server_address"] = *input.ServerAddress
	}

	if input.ServerPassword != nil {
		setMap["server_password"] = *input.ServerPassword
	}

	if input.Scenes != nil {
		setMap["scenes"] = *input.Scenes
	}

	if input.Sources != nil {
		setMap["sources"] = *input.Sources
	}

	if input.AudioSources != nil {
		setMap["audio_sources"] = *input.AudioSources
	}

	var (
		query  string
		args   []interface{}
		suffix = "RETURNING id, channel_id, server_port, server_address, server_password, scenes, sources, audio_sources, created_at, updated_at"
	)

	if m.IsNil() {
		query, args, _ = sq.Insert("channels_modules_obs_websocket").
			SetMap(setMap).
			Suffix(suffix).
			ToSql()
	} else {
		query, args, _ = sq.Update("channels_modules_obs_websocket").
			SetMap(setMap).
			Where(squirrel.Eq{"id": m.ID}).
			Suffix(suffix).
			ToSql()
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(
		ctx,
		query,
		args...,
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

func (c *Pgx) Delete(ctx context.Context, id int) error {
	query := `
DELETE FROM channels_modules_obs_websocket
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, id)
	return err
}
