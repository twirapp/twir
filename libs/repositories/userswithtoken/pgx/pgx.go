package pgx

import (
	"context"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	tokenmodel "github.com/twirapp/twir/libs/repositories/tokens/model"
	"github.com/twirapp/twir/libs/repositories/userswithtoken"
	"github.com/twirapp/twir/libs/repositories/userswithtoken/model"
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

var _ userswithtoken.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetByID(ctx context.Context, userID string) (model.UserWithToken, error) {
	query := `
SELECT u.id, u."isTester", u."isBotAdmin", u."tokenId", u."apiKey", u.hide_on_landing_page, u.is_banned, t.id token_id, t."accessToken" token_access_token, t."refreshToken" token_refresh_token, t."obtainmentTimestamp" token_obtainment_timestamp, t.scopes token_scopes, t."expiresIn" token_expires_in
FROM users u
LEFT JOIN tokens t ON u."tokenId" = t.id
WHERE u.id = $1
`

	var user model.UserWithToken
	var (
		tokenID                             uuid.UUID
		tokenAccessToken, tokenRefreshToken string
		tokenObtainmentTimestamp            time.Time
		tokenScopes                         []string
		tokenExpiresIn                      int
	)

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	err := conn.QueryRow(ctx, query, userID).Scan(
		&user.User.ID,
		&user.User.IsTester,
		&user.User.IsBotAdmin,
		&user.User.TokenID,
		&user.User.ApiKey,
		&user.User.HideOnLandingPage,
		&user.User.IsBanned,
		&tokenID,
		&tokenAccessToken,
		&tokenRefreshToken,
		&tokenObtainmentTimestamp,
		&tokenScopes,
		&tokenExpiresIn,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.UserWithToken{}, userswithtoken.ErrNotFound
		}
		return user, err
	}

	if tokenID != uuid.Nil {
		user.Token = &tokenmodel.Token{
			ID:                  tokenID,
			AccessToken:         tokenAccessToken,
			RefreshToken:        tokenRefreshToken,
			ObtainmentTimestamp: tokenObtainmentTimestamp,
			Scopes:              tokenScopes,
			ExpiresIn:           tokenExpiresIn,
		}
	}

	return user, err
}
