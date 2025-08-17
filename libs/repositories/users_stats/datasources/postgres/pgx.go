package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories"
	usersstats "github.com/twirapp/twir/libs/repositories/users_stats"
	"github.com/twirapp/twir/libs/repositories/users_stats/model"
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

var _ usersstats.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

var selectFields = []string{
	`id`,
	`messages`,
	`watched`,
	`"channelId"`,
	`"userId"`,
	`"usedChannelPoints"`,
	`is_mod`,
	`is_vip`,
	`is_subscriber`,
	`reputation`,
	`emotes`,
	`created_at`,
	`updated_at`,
}

var selectFieldsJoined string

func init() {
	selectFieldsJoined = strings.Join(selectFields, ", ")
}

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (*model.UserStat, error) {
	query := `
SELECT id, messages, watched, "channelId", "userId", "usedChannelPoints", is_mod, is_vip, is_subscriber, reputation, emotes, created_at, updated_at
FROM users_stats
WHERE id = $1
LIMIT 1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	stat, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.UserStat])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, usersstats.ErrNotFound
		}
		return nil, err
	}

	return &stat, nil
}

func (c *Pgx) Create(ctx context.Context, input usersstats.CreateInput) (*model.UserStat, error) {
	insertBuilder := sq.Insert("users_stats").
		SetMap(
			map[string]any{
				`"userId"`:            input.UserID,
				`"channelId"`:         input.ChannelID,
				`messages`:            input.Messages,
				`watched`:             input.Watched,
				`"usedChannelPoints"`: input.UsedChannelPoints,
				`is_mod`:              input.IsMod,
				`is_vip`:              input.IsVip,
				`is_subscriber`:       input.IsSubscriber,
				`reputation`:          input.Reputation,
				`emotes`:              input.Emotes,
			},
		).
		Suffix(fmt.Sprintf("RETURNING %s", selectFieldsJoined))

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	stat, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.UserStat])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, usersstats.ErrNotFound
		}
		return nil, err
	}

	return &stat, nil
}

func (c *Pgx) Update(
	ctx context.Context,
	userID, channelID string,
	input usersstats.UpdateInput,
) (*model.UserStat, error) {
	updateBuilder := sq.Update("users_stats").
		Where(
			squirrel.Eq{
				`"userId"`:    userID,
				`"channelId"`: channelID,
			},
		)

	updateBuilder = repositories.SquirrelApplyPatch(
		updateBuilder,
		map[string]any{
			`is_mod`:        input.IsMod,
			`is_vip`:        input.IsVip,
			`is_subscriber`: input.IsSubscriber,
			`updated_at`:    squirrel.Expr("NOW()"),
		},
	)

	for name, value := range input.NumberFields {
		fieldName := fmt.Sprintf(`"%s"`, name)

		if value.IsIncrement {
			updateBuilder = updateBuilder.Set(
				fieldName,
				squirrel.Expr(fmt.Sprintf(`%s + %d`, fieldName, value.Count)),
			)
		} else {
			updateBuilder = updateBuilder.Set(fieldName, value.Count)
		}
	}

	updateBuilder = updateBuilder.Suffix(fmt.Sprintf("RETURNING %s", selectFieldsJoined))
	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build users update query: %w", err)
	}
	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute users update query: %w", err)
	}

	stat, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.UserStat])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, usersstats.ErrNotFound
		}
		return nil, fmt.Errorf("failed to collect users update result: %w", err)
	}

	return &stat, nil
}

func (c *Pgx) GetByUserAndChannelID(
	ctx context.Context,
	userID, channelID string,
) (*model.UserStat, error) {
	query := fmt.Sprintf(
		`
SELECT %s
FROM users_stats
WHERE "userId" = $1 AND "channelId" = $2
LIMIT 1
`, selectFieldsJoined,
	)

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, userID, channelID)
	if err != nil {
		return nil, err
	}

	stat, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.UserStat])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, usersstats.ErrNotFound
		}

		return nil, err
	}

	return &stat, nil
}
