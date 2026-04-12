package pgx

import (
	"context"
	"errors"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	entity "github.com/twirapp/twir/libs/entities/kick_bot"
	"github.com/twirapp/twir/libs/repositories/kick_bots"
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

var _ kick_bots.Repository = (*Pgx)(nil)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

type dbModel struct {
	ID                  uuid.UUID `db:"id"`
	Type                string    `db:"type"`
	AccessToken         string    `db:"access_token"`
	RefreshToken        string    `db:"refresh_token"`
	Scopes              []string  `db:"scopes"`
	ExpiresIn           int       `db:"expires_in"`
	ObtainmentTimestamp time.Time `db:"obtainment_timestamp"`
	KickUserID          string    `db:"kick_user_id"`
	KickUserLogin       string    `db:"kick_user_login"`
	CreatedAt           time.Time `db:"created_at"`

	isNil bool
}

func (m dbModel) IsNil() bool { return m.isNil }

var Nil = &dbModel{isNil: true}

func modelToEntity(m dbModel) entity.KickBot {
	return entity.KickBot{
		ID:                  m.ID.String(),
		Type:                m.Type,
		AccessToken:         m.AccessToken,
		RefreshToken:        m.RefreshToken,
		Scopes:              m.Scopes,
		ExpiresIn:           m.ExpiresIn,
		ObtainmentTimestamp: m.ObtainmentTimestamp,
		KickUserID:          m.KickUserID,
		KickUserLogin:       m.KickUserLogin,
		CreatedAt:           m.CreatedAt,
	}
}

const selectColumns = `
id, type, access_token, refresh_token, scopes, expires_in,
obtainment_timestamp, kick_user_id, kick_user_login, created_at
`

func (r *Pgx) GetDefault(ctx context.Context) (entity.KickBot, error) {
	query := `SELECT ` + selectColumns + ` FROM kick_bots WHERE type = 'DEFAULT' LIMIT 1`

	conn := r.getter.DefaultTrOrDB(ctx, r.pool)
	rows, err := conn.Query(ctx, query)
	if err != nil {
		return entity.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[dbModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Nil, kick_bots.ErrNotFound
		}
		return entity.Nil, err
	}

	return modelToEntity(result), nil
}

func (r *Pgx) GetByID(ctx context.Context, id uuid.UUID) (entity.KickBot, error) {
	query := `SELECT ` + selectColumns + ` FROM kick_bots WHERE id = $1`

	conn := r.getter.DefaultTrOrDB(ctx, r.pool)
	rows, err := conn.Query(ctx, query, id)
	if err != nil {
		return entity.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[dbModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Nil, kick_bots.ErrNotFound
		}
		return entity.Nil, err
	}

	return modelToEntity(result), nil
}

func (r *Pgx) Create(ctx context.Context, input kick_bots.CreateInput) (entity.KickBot, error) {
	query := `
INSERT INTO kick_bots (
	type, access_token, refresh_token, scopes, expires_in, obtainment_timestamp,
	kick_user_id, kick_user_login
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING ` + selectColumns

	conn := r.getter.DefaultTrOrDB(ctx, r.pool)
	rows, err := conn.Query(
		ctx,
		query,
		input.Type,
		input.AccessToken,
		input.RefreshToken,
		input.Scopes,
		input.ExpiresIn,
		input.ObtainmentTimestamp,
		input.KickUserID,
		input.KickUserLogin,
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

func (r *Pgx) UpdateToken(ctx context.Context, id uuid.UUID, input kick_bots.UpdateTokenInput) (entity.KickBot, error) {
	query := `
UPDATE kick_bots SET
	access_token = $2,
	refresh_token = $3,
	scopes = $4,
	expires_in = $5,
	obtainment_timestamp = $6
WHERE id = $1
RETURNING ` + selectColumns

	conn := r.getter.DefaultTrOrDB(ctx, r.pool)
	rows, err := conn.Query(
		ctx,
		query,
		id,
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
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Nil, kick_bots.ErrNotFound
		}
		return entity.Nil, err
	}

	return modelToEntity(result), nil
}
