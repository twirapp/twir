package pgx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/events"
	"github.com/twirapp/twir/libs/repositories/events/model"
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

var _ events.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetManyByChannelID(ctx context.Context, channelID string) ([]model.Event, error) {
	query := `
SELECT 
    e.id, 
    e."channelId", 
    e.type, 
    e."rewardId", 
    e."commandId", 
    e."keywordId", 
    e.description, 
    e.enabled, 
    e.online_only,
    COALESCE(
        (
            SELECT json_agg(
                json_build_object(
                    'id', op.id,
                    'type', op.type,
                    'input', op.input,
                    'delay', op.delay,
                    'repeat', op.repeat,
                    'useAnnounce', op."useAnnounce",
                    'timeoutTime', op."timeoutTime",
                    'timeoutMessage', op."timeoutMessage",
                    'target', op.target,
                    'enabled', op.enabled,
                    'filters', (
                        SELECT COALESCE(
                            json_agg(
                                json_build_object(
                                    'id', f.id,
                                    'type', f.type,
                                    'left', f.left,
                                    'right', f.right
                                )
                            ),
                            '[]'::json
                        )
                        FROM channels_events_operations_filters f
                        WHERE f."operationId" = op.id
                    )
                )
            )
            FROM channels_events_operations op
            WHERE op."eventId" = e.id
        ),
        '[]'::json
    ) as operations
FROM channels_events e
WHERE e."channelId" = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, channelID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var events []model.Event
	for rows.Next() {
		var event model.Event
		var operationsJSON []byte
		
		err := rows.Scan(
			&event.ID,
			&event.ChannelID,
			&event.Type,
			&event.RewardID,
			&event.CommandID,
			&event.KeywordID,
			&event.Description,
			&event.Enabled,
			&event.OnlineOnly,
			&operationsJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		if err := json.Unmarshal(operationsJSON, &event.Operations); err != nil {
			return nil, fmt.Errorf("unmarshal operations: %w", err)
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return events, nil
}

func (c *Pgx) GetByID(ctx context.Context, id string) (model.Event, error) {
	query := `
SELECT 
    e.id, 
    e."channelId", 
    e.type, 
    e."rewardId", 
    e."commandId", 
    e."keywordId", 
    e.description, 
    e.enabled, 
    e.online_only,
    COALESCE(
        (
            SELECT json_agg(
                json_build_object(
                    'id', op.id,
                    'type', op.type,
                    'input', op.input,
                    'delay', op.delay,
                    'repeat', op.repeat,
                    'useAnnounce', op."useAnnounce",
                    'timeoutTime', op."timeoutTime",
                    'timeoutMessage', op."timeoutMessage",
                    'target', op.target,
                    'enabled', op.enabled,
                    'filters', (
                        SELECT COALESCE(
                            json_agg(
                                json_build_object(
                                    'id', f.id,
                                    'type', f.type,
                                    'left', f.left,
                                    'right', f.right
                                )
                            ),
                            '[]'::json
                        )
                        FROM channels_events_operations_filters f
                        WHERE f."operationId" = op.id
                    )
                )
            )
            FROM channels_events_operations op
            WHERE op."eventId" = e.id
        ),
        '[]'::json
    ) as operations
FROM channels_events e
WHERE e.id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	row := conn.QueryRow(ctx, query, id)

	var event model.Event
	var operationsJSON []byte
	
	err := row.Scan(
		&event.ID,
		&event.ChannelID,
		&event.Type,
		&event.RewardID,
		&event.CommandID,
		&event.KeywordID,
		&event.Description,
		&event.Enabled,
		&event.OnlineOnly,
		&operationsJSON,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, events.ErrNotFound
		}
		return model.Nil, fmt.Errorf("scan: %w", err)
	}

	if err := json.Unmarshal(operationsJSON, &event.Operations); err != nil {
		return model.Nil, fmt.Errorf("unmarshal operations: %w", err)
	}

	return event, nil
}

func (c *Pgx) Create(ctx context.Context, input events.CreateInput) (model.Event, error) {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return model.Nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Create event
	eventID := uuid.New().String()
	query := `
INSERT INTO channels_events (id, "channelId", type, "rewardId", "commandId", "keywordId", description, enabled, online_only)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id
`

	_, err = tx.Exec(
		ctx,
		query,
		eventID,
		input.ChannelID,
		input.Type,
		input.RewardID,
		input.CommandID,
		input.KeywordID,
		input.Description,
		input.Enabled,
		input.OnlineOnly,
	)
	if err != nil {
		return model.Nil, fmt.Errorf("insert event: %w", err)
	}

	// Create operations
	for _, op := range input.Operations {
		operationID := uuid.New().String()
		query := `
INSERT INTO channels_events_operations (id, "eventId", type, input, delay, repeat, "useAnnounce", "timeoutTime", "timeoutMessage", target, enabled)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id
`

		_, err = tx.Exec(
			ctx,
			query,
			operationID,
			eventID,
			op.Type,
			op.Input,
			op.Delay,
			op.Repeat,
			op.UseAnnounce,
			op.TimeoutTime,
			op.TimeoutMessage,
			op.Target,
			op.Enabled,
		)
		if err != nil {
			return model.Nil, fmt.Errorf("insert operation: %w", err)
		}

		// Create filters
		for _, filter := range op.Filters {
			filterID := uuid.New().String()
			query := `
INSERT INTO channels_events_operations_filters (id, "operationId", type, left, right)
VALUES ($1, $2, $3, $4, $5)
`

			_, err = tx.Exec(
				ctx,
				query,
				filterID,
				operationID,
				filter.Type,
				filter.Left,
				filter.Right,
			)
			if err != nil {
				return model.Nil, fmt.Errorf("insert filter: %w", err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Nil, fmt.Errorf("commit transaction: %w", err)
	}

	// Return the created event
	return c.GetByID(ctx, eventID)
}

func (c *Pgx) Update(ctx context.Context, id string, input events.UpdateInput) (model.Event, error) {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return model.Nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Update event
	updateBuilder := sq.Update("channels_events").Where(squirrel.Eq{"id": id})

	if input.Type != nil {
		updateBuilder = updateBuilder.Set("type", *input.Type)
	}
	if input.RewardID != nil {
		updateBuilder = updateBuilder.Set(`"rewardId"`, *input.RewardID)
	}
	if input.CommandID != nil {
		updateBuilder = updateBuilder.Set(`"commandId"`, *input.CommandID)
	}
	if input.KeywordID != nil {
		updateBuilder = updateBuilder.Set(`"keywordId"`, *input.KeywordID)
	}
	if input.Description != nil {
		updateBuilder = updateBuilder.Set("description", *input.Description)
	}
	if input.Enabled != nil {
		updateBuilder = updateBuilder.Set("enabled", *input.Enabled)
	}
	if input.OnlineOnly != nil {
		updateBuilder = updateBuilder.Set("online_only", *input.OnlineOnly)
	}

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return model.Nil, fmt.Errorf("build update query: %w", err)
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return model.Nil, fmt.Errorf("update event: %w", err)
	}

	// If operations are provided, delete existing ones and create new ones
	if input.Operations != nil {
		// Delete existing operations (cascade will delete filters)
		_, err = tx.Exec(ctx, `DELETE FROM channels_events_operations WHERE "eventId" = $1`, id)
		if err != nil {
			return model.Nil, fmt.Errorf("delete operations: %w", err)
		}

		// Create new operations
		for _, op := range *input.Operations {
			operationID := uuid.New().String()
			query := `
INSERT INTO channels_events_operations (id, "eventId", type, input, delay, repeat, "useAnnounce", "timeoutTime", "timeoutMessage", target, enabled)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id
`

			_, err = tx.Exec(
				ctx,
				query,
				operationID,
				id,
				op.Type,
				op.Input,
				op.Delay,
				op.Repeat,
				op.UseAnnounce,
				op.TimeoutTime,
				op.TimeoutMessage,
				op.Target,
				op.Enabled,
			)
			if err != nil {
				return model.Nil, fmt.Errorf("insert operation: %w", err)
			}

			// Create filters
			for _, filter := range op.Filters {
				filterID := uuid.New().String()
				query := `
INSERT INTO channels_events_operations_filters (id, "operationId", type, left, right)
VALUES ($1, $2, $3, $4, $5)
`

				_, err = tx.Exec(
					ctx,
					query,
					filterID,
					operationID,
					filter.Type,
					filter.Left,
					filter.Right,
				)
				if err != nil {
					return model.Nil, fmt.Errorf("insert filter: %w", err)
				}
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Nil, fmt.Errorf("commit transaction: %w", err)
	}

	// Return the updated event
	return c.GetByID(ctx, id)
}

func (c *Pgx) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM channels_events WHERE id = $1`
	
	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete event: %w", err)
	}
	
	return nil
}
