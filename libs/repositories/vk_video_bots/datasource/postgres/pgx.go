package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	entity "github.com/twirapp/twir/libs/entities/vk_video_bot"
	vkvideobots "github.com/twirapp/twir/libs/repositories/vk_video_bots"
)

const singletonLockQuery = `SELECT pg_advisory_xact_lock(761247388)`

const selectColumns = `
	id,
	encrypted_access_token,
	encrypted_refresh_token,
	scopes,
	expires_in,
	obtainment_timestamp,
	vk_user_id,
	created_at,
	updated_at`

type Opts struct {
	PgxPool *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{
		pool:            opts.PgxPool,
		transactionPool: opts.PgxPool,
		getter:          trmpgx.DefaultCtxGetter,
	}
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return New(Opts{PgxPool: pool})
}

var _ vkvideobots.Repository = (*Pgx)(nil)

type database interface {
	Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

type Pgx struct {
	pool            database
	transactionPool *pgxpool.Pool
	getter          *trmpgx.CtxGetter
}

type dbModel struct {
	ID                    uuid.UUID `db:"id"`
	EncryptedAccessToken  string    `db:"encrypted_access_token"`
	EncryptedRefreshToken string    `db:"encrypted_refresh_token"`
	Scopes                []string  `db:"scopes"`
	ExpiresIn             int       `db:"expires_in"`
	ObtainmentTimestamp   time.Time `db:"obtainment_timestamp"`
	VKUserID              uuid.UUID `db:"vk_user_id"`
	CreatedAt             time.Time `db:"created_at"`
	UpdatedAt             time.Time `db:"updated_at"`

	isNil bool
}

func (m dbModel) IsNil() bool { return m.isNil }

var Nil = dbModel{isNil: true}

func (r *Pgx) Get(ctx context.Context) (entity.VKVideoBot, error) {
	query := `
		WITH locked AS (` + singletonLockQuery + `)
		SELECT ` + selectColumns + `
		FROM vk_video_bots
		CROSS JOIN locked
		WHERE singleton
		LIMIT 1`

	return r.queryOne(ctx, query)
}

func (r *Pgx) Lock(ctx context.Context) error {
	if _, err := r.connection(ctx).Exec(ctx, singletonLockQuery); err != nil {
		return fmt.Errorf("lock VK Video bot singleton: %w", err)
	}

	return nil
}

func (r *Pgx) Upsert(ctx context.Context, input vkvideobots.UpsertInput) (entity.VKVideoBot, error) {
	query := `
		WITH locked AS (` + singletonLockQuery + `)
		INSERT INTO vk_video_bots (
			singleton,
			encrypted_access_token,
			encrypted_refresh_token,
			scopes,
			expires_in,
			obtainment_timestamp,
			vk_user_id
		)
		SELECT TRUE, $1, $2, $3, $4, $5, $6
		FROM locked
		ON CONFLICT (singleton) DO UPDATE SET
			encrypted_access_token = EXCLUDED.encrypted_access_token,
			encrypted_refresh_token = EXCLUDED.encrypted_refresh_token,
			scopes = EXCLUDED.scopes,
			expires_in = EXCLUDED.expires_in,
			obtainment_timestamp = EXCLUDED.obtainment_timestamp,
			vk_user_id = EXCLUDED.vk_user_id,
			updated_at = NOW()
		RETURNING ` + selectColumns

	return r.queryOne(
		ctx,
		query,
		input.EncryptedAccessToken,
		input.EncryptedRefreshToken,
		input.Scopes,
		input.ExpiresIn,
		input.ObtainmentTimestamp,
		input.VKUserID,
	)
}

func (r *Pgx) Update(ctx context.Context, input vkvideobots.UpdateInput) (entity.VKVideoBot, error) {
	query := `
		WITH locked AS (` + singletonLockQuery + `)
		UPDATE vk_video_bots
		SET
			encrypted_access_token = $1,
			encrypted_refresh_token = $2,
			scopes = $3,
			expires_in = $4,
			obtainment_timestamp = $5,
			vk_user_id = $6,
			updated_at = NOW()
		FROM locked
		WHERE singleton
		RETURNING ` + selectColumns

	return r.queryOne(
		ctx,
		query,
		input.EncryptedAccessToken,
		input.EncryptedRefreshToken,
		input.Scopes,
		input.ExpiresIn,
		input.ObtainmentTimestamp,
		input.VKUserID,
	)
}

func (r *Pgx) queryOne(ctx context.Context, query string, args ...any) (entity.VKVideoBot, error) {
	rows, err := r.connection(ctx).Query(ctx, query, args...)
	if err != nil {
		return entity.Nil, fmt.Errorf("query VK Video bot singleton: %w", err)
	}

	model, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[dbModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Nil, vkvideobots.ErrNotFound
		}
		return entity.Nil, fmt.Errorf("collect VK Video bot singleton: %w", err)
	}

	return mapToEntity(model), nil
}

func (r *Pgx) connection(ctx context.Context) database {
	if r.transactionPool != nil {
		return r.getter.DefaultTrOrDB(ctx, r.transactionPool)
	}

	return r.pool
}

func mapToEntity(model dbModel) entity.VKVideoBot {
	if model.IsNil() {
		return entity.Nil
	}

	return entity.VKVideoBot{
		ID:                    model.ID,
		EncryptedAccessToken:  model.EncryptedAccessToken,
		EncryptedRefreshToken: model.EncryptedRefreshToken,
		Scopes:                model.Scopes,
		ExpiresIn:             model.ExpiresIn,
		ObtainmentTimestamp:   model.ObtainmentTimestamp,
		VKUserID:              model.VKUserID,
		CreatedAt:             model.CreatedAt,
		UpdatedAt:             model.UpdatedAt,
	}
}
