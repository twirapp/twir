package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	channelredemptionshistory "github.com/twirapp/twir/libs/repositories/channel_redemptions_history"
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

var _ channelredemptionshistory.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) CreateMany(
	ctx context.Context,
	inputs []channelredemptionshistory.CreateInput,
) error {
	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.CopyFrom(
		ctx,
		pgx.Identifier{"channel_redemptions_history"},
		[]string{
			"channel_id",
			"user_id",
			"reward_id",
			"reward_title",
			"reward_prompt",
			"reward_cost",
		},
		pgx.CopyFromSlice(
			len(inputs),
			func(i int) ([]any, error) {
				return []any{
					inputs[i].ChannelID,
					inputs[i].UserID,
					inputs[i].RewardID,
					inputs[i].RewardTitle,
					inputs[i].RewardPrompt,
					inputs[i].RewardCost,
				}, nil
			},
		),
	)
	if err != nil {
		return fmt.Errorf("failed to copy from: %w", err)
	}

	return nil
}
