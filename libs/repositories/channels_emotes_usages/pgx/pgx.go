package pgx

import (
	"context"
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/channels_emotes_usages"
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

var _ channels_emotes_usages.Repository = (*Pgx)(nil)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (p Pgx) CreateMany(
	ctx context.Context,
	inputs []channels_emotes_usages.ChannelEmoteUsageInput,
) error {
	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	_, err := conn.CopyFrom(
		ctx,
		pgx.Identifier{"channels_emotes_usages"},
		[]string{"id", "channelId", "userId", "createdAt", "emote"},
		pgx.CopyFromSlice(
			len(inputs),
			func(i int) ([]any, error) {
				return []any{
					inputs[i].ID,
					inputs[i].ChannelID,
					inputs[i].UserID,
					inputs[i].CreatedAt,
					inputs[i].Emote,
				}, nil
			},
		),
	)
	if err != nil {
		return fmt.Errorf("failed to copy from: %w", err)
	}

	return nil
}
