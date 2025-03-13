package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/bots"
	"github.com/twirapp/twir/libs/repositories/bots/model"
	tokenmodel "github.com/twirapp/twir/libs/repositories/tokens/model"
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

var _ bots.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

type scanModel struct {
	botId      string    `db:"bot_id"`
	botType    string    `db:"bot_type"`
	botTokenId uuid.UUID `db:"bot_token_id"`

	tokenId                  uuid.UUID `db:"token_id"`
	tokenAccessToken         string    `db:"token_access_token"`
	tokenRefreshToken        string    `db:"token_refresh_token"`
	tokenExpiresIn           int       `db:"token_expires_in"`
	tokenObtainmentTimestamp time.Time `db:"token_obtainment_timestamp"`
	tokenScopes              []string  `db:"token_scopes"`
}

func (c *Pgx) GetDefault(ctx context.Context) (model.Bot, error) {
	query := `
SELECT b.id bot_id,
       b.type bot_type,
       b."tokenId" bot_token_id,
       t.id token_id,
       t."accessToken" token_access_token,
       t."refreshToken" token_refresh_token,
       t."expiresIn" token_expires_in,
       t."obtainmentTimestamp" token_obtainment_timestamp,
       t.scopes token_scopes
FROM bots b
LEFT JOIN tokens t ON b."tokenId" = t.id
WHERE b.type = 'DEFAULT'
LIMIT 1;
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query)
	if err != nil {
		return model.Nil, err
	}
	defer rows.Close()

	data := model.Bot{
		Token: &tokenmodel.Token{},
	}
	for rows.Next() {
		err := rows.Scan(
			&data.ID,
			&data.Type,
			&data.TokenID,

			&data.Token.ID,
			&data.Token.AccessToken,
			&data.Token.RefreshToken,
			&data.Token.ExpiresIn,
			&data.Token.ObtainmentTimestamp,
			&data.Token.Scopes,
		)

		if err != nil {
			return model.Nil, err
		}
	}

	if rows.Err() != nil {
		return model.Nil, fmt.Errorf("collect: %w", rows.Err())
	}

	return data, nil
}
