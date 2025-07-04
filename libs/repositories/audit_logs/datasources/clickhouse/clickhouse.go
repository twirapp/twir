package clickhouse

import (
	"context"

	"github.com/Masterminds/squirrel"
	twirclickhouse "github.com/twirapp/twir/libs/baseapp/clickhouse"
	"github.com/twirapp/twir/libs/repositories/audit_logs"
	"github.com/twirapp/twir/libs/repositories/audit_logs/model"
)

type Opts struct {
	Client *twirclickhouse.ClickhouseClient
}

func New(opts Opts) *Clickhouse {
	return &Clickhouse{
		client: opts.Client,
	}
}

func NewFx(client *twirclickhouse.ClickhouseClient) *Clickhouse {
	return New(Opts{Client: client})
}

var _ audit_logs.Repository = (*Clickhouse)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)

type Clickhouse struct {
	client *twirclickhouse.ClickhouseClient
}

func (c *Clickhouse) GetMany(ctx context.Context, input audit_logs.GetManyInput) (
	[]model.AuditLog,
	error,
) {
	selectBuilder := sq.
		Select(
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

	rows, err := c.client.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.AuditLog
	for rows.Next() {
		var log model.AuditLog
		err := rows.Scan(
			&log.TableName,
			&log.OperationType,
			&log.OldValue,
			&log.NewValue,
			&log.ObjectID,
			&log.ChannelID,
			&log.UserID,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, log)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}

func (c *Clickhouse) Count(ctx context.Context, input audit_logs.GetCountInput) (uint64, error) {
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

	var count uint64
	err = c.client.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *Clickhouse) Create(ctx context.Context, input audit_logs.CreateInput) error {
	query := `
INSERT INTO audit_logs (table_name, operation_type, old_value, new_value, object_id, channel_id, user_id)
VALUES (?, ?, ?, ?, ?, ?, ?);
`

	err := c.client.Exec(
		ctx,
		query,
		input.Table,
		string(input.OperationType),
		input.OldValue,
		input.NewValue,
		input.ObjectID,
		input.ChannelID,
		input.UserID,
	)

	return err
}
