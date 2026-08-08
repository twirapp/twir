package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	mcpOAuth "github.com/twirapp/twir/libs/repositories/mcp_oauth"
)

type Opts struct {
	PgxPool *pgxpool.Pool
}

type Pgx struct {
	pool *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{pool: opts.PgxPool}
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return New(Opts{PgxPool: pool})
}

var _ mcpOAuth.Repository = (*Pgx)(nil)
