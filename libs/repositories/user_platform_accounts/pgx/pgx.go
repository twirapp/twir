package pgx

import (
	"context"
	"errors"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/entities/platform"
	entity "github.com/twirapp/twir/libs/entities/user_platform_account"
	"github.com/twirapp/twir/libs/repositories/user_platform_accounts"
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

var _ user_platform_accounts.Repository = (*Pgx)(nil)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

type dbModel struct {
	ID                  uuid.UUID `db:"id"`
	UserID              uuid.UUID `db:"user_id"`
	Platform            string    `db:"platform"`
	PlatformUserID      string    `db:"platform_user_id"`
	PlatformLogin       string    `db:"platform_login"`
	PlatformDisplayName string    `db:"platform_display_name"`
	PlatformAvatar      string    `db:"platform_avatar"`
	AccessToken         string    `db:"access_token"`
	RefreshToken        string    `db:"refresh_token"`
	Scopes              []string  `db:"scopes"`
	ExpiresIn           int       `db:"expires_in"`
	ObtainmentTimestamp time.Time `db:"obtainment_timestamp"`
}

func modelToEntity(m dbModel) entity.UserPlatformAccount {
	return entity.UserPlatformAccount{
		ID:                  m.ID,
		UserID:              m.UserID,
		Platform:            platform.Platform(m.Platform),
		PlatformUserID:      m.PlatformUserID,
		PlatformLogin:       m.PlatformLogin,
		PlatformDisplayName: m.PlatformDisplayName,
		PlatformAvatar:      m.PlatformAvatar,
		AccessToken:         m.AccessToken,
		RefreshToken:        m.RefreshToken,
		Scopes:              m.Scopes,
		ExpiresIn:           m.ExpiresIn,
		ObtainmentTimestamp: m.ObtainmentTimestamp,
	}
}

const selectColumns = `
id, user_id, platform::text AS platform, platform_user_id, platform_login,
platform_display_name, platform_avatar, access_token, refresh_token,
scopes, expires_in, obtainment_timestamp
`

func (r *Pgx) GetByUserIDAndPlatform(ctx context.Context, userID uuid.UUID, plat platform.Platform) (entity.UserPlatformAccount, error) {
	query := `SELECT ` + selectColumns + ` FROM user_platform_accounts WHERE user_id = $1 AND platform = $2`

	conn := r.getter.DefaultTrOrDB(ctx, r.pool)
	rows, err := conn.Query(ctx, query, userID, plat)
	if err != nil {
		return entity.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[dbModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Nil, user_platform_accounts.ErrNotFound
		}
		return entity.Nil, err
	}

	return modelToEntity(result), nil
}

func (r *Pgx) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]entity.UserPlatformAccount, error) {
	query := `SELECT ` + selectColumns + ` FROM user_platform_accounts WHERE user_id = $1`

	conn := r.getter.DefaultTrOrDB(ctx, r.pool)
	rows, err := conn.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbModel])
	if err != nil {
		return nil, err
	}

	result := make([]entity.UserPlatformAccount, len(models))
	for i, m := range models {
		result[i] = modelToEntity(m)
	}

	return result, nil
}

func (r *Pgx) GetByPlatformUserID(ctx context.Context, plat platform.Platform, platformUserID string) (entity.UserPlatformAccount, error) {
	query := `SELECT ` + selectColumns + ` FROM user_platform_accounts WHERE platform = $1 AND platform_user_id = $2`

	conn := r.getter.DefaultTrOrDB(ctx, r.pool)
	rows, err := conn.Query(ctx, query, plat, platformUserID)
	if err != nil {
		return entity.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[dbModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Nil, user_platform_accounts.ErrNotFound
		}
		return entity.Nil, err
	}

	return modelToEntity(result), nil
}

func (r *Pgx) Upsert(ctx context.Context, input user_platform_accounts.UpsertInput) (entity.UserPlatformAccount, error) {
	query := `
INSERT INTO user_platform_accounts (
	user_id, platform, platform_user_id, platform_login, platform_display_name,
	platform_avatar, access_token, refresh_token, scopes, expires_in, obtainment_timestamp
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
ON CONFLICT (platform, platform_user_id) DO UPDATE SET
	access_token         = EXCLUDED.access_token,
	refresh_token        = EXCLUDED.refresh_token,
	scopes               = EXCLUDED.scopes,
	expires_in           = EXCLUDED.expires_in,
	obtainment_timestamp = EXCLUDED.obtainment_timestamp,
	platform_login       = EXCLUDED.platform_login,
	platform_display_name = EXCLUDED.platform_display_name,
	platform_avatar      = EXCLUDED.platform_avatar
RETURNING ` + selectColumns

	conn := r.getter.DefaultTrOrDB(ctx, r.pool)
	rows, err := conn.Query(
		ctx,
		query,
		input.UserID,
		input.Platform,
		input.PlatformUserID,
		input.PlatformLogin,
		input.PlatformDisplayName,
		input.PlatformAvatar,
		input.AccessToken,
		input.RefreshToken,
		input.Scopes,
		input.ExpiresIn,
		input.ObtainmentTimestamp,
	)
	if err != nil {
		return entity.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[dbModel])
	if err != nil {
		return entity.Nil, err
	}

	return modelToEntity(result), nil
}

func (r *Pgx) GetAllByPlatform(ctx context.Context, plat platform.Platform) ([]entity.UserPlatformAccount, error) {
	query := `SELECT ` + selectColumns + ` FROM user_platform_accounts WHERE platform = $1::platform`

	conn := r.getter.DefaultTrOrDB(ctx, r.pool)
	rows, err := conn.Query(ctx, query, plat)
	if err != nil {
		return nil, err
	}

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbModel])
	if err != nil {
		return nil, err
	}

	result := make([]entity.UserPlatformAccount, len(models))
	for i, m := range models {
		result[i] = modelToEntity(m)
	}

	return result, nil
}

func (r *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM user_platform_accounts WHERE id = $1`

	conn := r.getter.DefaultTrOrDB(ctx, r.pool)
	_, err := conn.Exec(ctx, query, id)
	return err
}
