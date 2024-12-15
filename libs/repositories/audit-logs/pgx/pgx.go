package pgx

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	audit_logs "github.com/twirapp/twir/libs/repositories/audit-logs"
	"github.com/twirapp/twir/libs/repositories/audit-logs/model"
)

type Opts struct {
	PgxPool *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{
		pool: opts.PgxPool,
	}
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return New(Opts{PgxPool: pool})
}

var _ audit_logs.Repository = (*Pgx)(nil)

type Pgx struct {
	pool *pgxpool.Pool
}

func (p Pgx) GetByChannelID(ctx context.Context, channelID string, limit int) (
	[]model.AuditLog,
	error,
) {
	query := `
SELECT id, table_name, operation_type, old_value, new_value, object_id, channel_id, user_id, created_at
FROM audit_logs
WHERE channel_id = $1
ORDER BY created_at DESC
LIMIT $2
`

	rows, err := p.pool.Query(ctx, query, channelID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.AuditLog])
	if err != nil {
		return nil, err
	}

	return logs, err
}
