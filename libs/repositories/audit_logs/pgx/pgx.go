package pgx

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/audit_logs"
	"github.com/twirapp/twir/libs/repositories/audit_logs/model"
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
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool *pgxpool.Pool
}

func (p *Pgx) GetByID(ctx context.Context, id uuid.UUID) (model.AuditLog, error) {
	// TODO implement me
	panic("implement me")
}

func (p *Pgx) GetMany(ctx context.Context, input audit_logs.GetManyInput) (
	[]model.AuditLog,
	error,
) {
	selectBuilder := sq.
		Select(
			"id",
			"table_name",
			"operation_type",
			"old_value",
			"new_value",
			"object_id",
			"channel_id",
			"user_id",
			"created_at",
		).
		From("audit_logs").
		OrderBy("created_at DESC")

	perPage := input.Limit
	if perPage == 0 {
		perPage = 20
	}

	offset := input.Page * perPage

	selectBuilder = selectBuilder.Limit(uint64(perPage)).Offset(uint64(offset)).OrderBy("created_at DESC")

	if input.ChannelID != nil {
		selectBuilder = selectBuilder.Where("channel_id = ?", input.ChannelID)
	}

	if input.ActorID != nil {
		selectBuilder = selectBuilder.Where("user_id = ?", input.ActorID)
	}

	if input.ObjectID != nil {
		selectBuilder = selectBuilder.Where("object_id = ?", input.ObjectID)
	}

	if len(input.OperationTypes) > 0 {
		operations := make([]string, 0, len(input.OperationTypes))
		for _, operation := range input.OperationTypes {
			operations = append(operations, string(operation))
		}

		selectBuilder = selectBuilder.Where(
			squirrel.Eq{
				"operation_type": operations,
			},
		)
	}

	if len(input.Systems) > 0 {
		selectBuilder = selectBuilder.Where(
			squirrel.Eq{
				"system": input.Systems,
			},
		)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	logs, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.AuditLog])
	if err != nil {
		return nil, err
	}

	return logs, err
}

func (p *Pgx) Count(ctx context.Context, input audit_logs.GetCountInput) (int, error) {
	selectBuilder := sq.
		Select("COUNT(*)").
		From("audit_logs")

	if input.ChannelID != nil {
		selectBuilder = selectBuilder.Where("channel_id = ?", input.ChannelID)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var count int
	err = p.pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (p *Pgx) Create(ctx context.Context, input audit_logs.CreateInput) (model.AuditLog, error) {
	insertBuilder := sq.Insert("audit_logs").
		SetMap(
			map[string]any{
				"table_name":     input.Table,
				"operation_type": string(input.OperationType),
				"old_value":      input.OldValue,
				"new_value":      input.NewValue,
				"object_id":      input.ObjectID,
				"channel_id":     input.ChannelID,
				"user_id":        input.UserID,
			},
		).
		Suffix("RETURNING id")

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		return model.Nil, err
	}

	var id uuid.UUID
	err = p.pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return model.Nil, err
	}

	return p.GetByID(ctx, id)
}
