package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	usermodel "github.com/twirapp/twir/libs/repositories/users/model"
	userstatsmodel "github.com/twirapp/twir/libs/repositories/users_stats/model"
	"github.com/twirapp/twir/libs/repositories/userswithstats"
	"github.com/twirapp/twir/libs/repositories/userswithstats/model"
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

var _ userswithstats.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

type userStatScanJsonRow struct {
	ID             uuid.UUID `json:"id"`
	UserID         string    `json:"user_id"`
	ChannelID      string    `json:"channel_id"`
	Messages       int       `json:"messages"`
	Emotes         int       `json:"emotes"`
	Watched        int       `json:"watched"`
	UsedChannelPts int       `json:"used_channel_emotes"`
	IsMod          bool      `json:"is_mod"`
	IsVip          bool      `json:"is_vip"`
	IsSubscriber   bool      `json:"is_subscriber"`
	Reputation     int       `json:"reputation"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (c *Pgx) GetByUserAndChannelID(
	ctx context.Context,
	input userswithstats.GetByUserAndChannelIDInput,
) (model.UserWithStats, error) {
	query := `
SELECT
	u.id, u."isBotAdmin", u."tokenId", u."apiKey", u.hide_on_landing_page, u.is_banned, u.created_at,
	CASE
    WHEN us.id IS NULL THEN NULL
    ELSE JSON_BUILD_OBJECT(
      'id', us.id,
      'user_id', us."userId",
      'channel_id', us."channelId",
      'messages', us."messages",
      'emotes', us."emotes",
      'watched', us."watched",
      'used_channel_emotes', us."usedChannelPoints",
      'is_mod', us."is_mod",
      'isVip', us."is_vip",
      'is_subscriber', us."is_subscriber",
      'reputation', us."reputation",
      'created_at', us."created_at",
      'updated_at', us."updated_at"
    )
  END AS stats
FROM users as u
LEFT JOIN users_stats us ON us."userId" = u.id AND us."channelId" = $2
WHERE u.id = $1
LIMIT 1;
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	row := conn.QueryRow(ctx, query, input.UserID, input.ChannelID)

	var userModel usermodel.User
	var userStatsRawJson sql.Null[[]byte]

	if err := row.Scan(
		&userModel.ID,
		&userModel.IsBotAdmin,
		&userModel.TokenID,
		&userModel.ApiKey,
		&userModel.HideOnLandingPage,
		&userModel.IsBanned,
		&userModel.CreatedAt,
		&userStatsRawJson,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.UserWithStats{}, userswithstats.ErrNotFound
		}

		return model.UserWithStats{}, fmt.Errorf("userswithstats: pgx scan row: %w", err)
	}

	var userStats *userstatsmodel.UserStat
	if userStatsRawJson.Valid {
		var statsRow userStatScanJsonRow
		if err := json.Unmarshal(userStatsRawJson.V, &statsRow); err != nil {
			return model.UserWithStats{}, fmt.Errorf("userswithstats: unmarshal user stats json: %w", err)
		}
		userStats = &userstatsmodel.UserStat{
			ID:                statsRow.ID,
			UserID:            statsRow.UserID,
			ChannelID:         statsRow.ChannelID,
			Messages:          int32(statsRow.Messages),
			Emotes:            statsRow.Emotes,
			Watched:           int64(statsRow.Watched),
			UsedChannelPoints: int64(statsRow.UsedChannelPts),
			IsMod:             statsRow.IsMod,
			IsVip:             statsRow.IsVip,
			IsSubscriber:      statsRow.IsSubscriber,
			Reputation:        int64(statsRow.Reputation),
			CreatedAt:         statsRow.CreatedAt,
			UpdatedAt:         statsRow.UpdatedAt,
		}
	}

	return model.UserWithStats{
		User:  userModel,
		Stats: userStats,
	}, nil
}
