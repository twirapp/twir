package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/seventv_integration"
	"github.com/twirapp/twir/libs/repositories/seventv_integration/model"
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

var _ seventv_integration.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) Create(ctx context.Context, input seventv_integration.CreateInput) error {
	query := `
INSERT INTO channels_integrations_seventv (channel_id, reward_id_for_add_emote, reward_id_for_remove_emote, delete_emotes_only_added_by_app)
VALUES ($1, $2, $3, $4)
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(
		ctx,
		query,
		input.ChannelID,
		input.RewardIdForAddEmote,
		input.RewardIdForRemoveEmote,
		input.DeleteEmotesOnlyAddedByApp,
	)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}

func (c *Pgx) GetByChannelID(ctx context.Context, channelID string) (
	model.SevenTvIntegration,
	error,
) {
	query := `
SELECT id, channel_id, reward_id_for_add_emote, reward_id_for_remove_emote, delete_emotes_only_added_by_app, added_emotes
FROM channels_integrations_seventv
WHERE channel_id = $1
LIMIT 1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, channelID)
	if err != nil {
		return model.Nil, fmt.Errorf("query: %w", err)
	}

	data, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.SevenTvIntegration])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, nil
		}

		return model.Nil, fmt.Errorf("collect: %w", err)
	}

	return data, nil
}

func (c *Pgx) Update(
	ctx context.Context,
	id uuid.UUID,
	input seventv_integration.UpdateInput,
) error {
	updateBuilder := sq.Update("channels_integrations_seventv").
		Where(squirrel.Eq{"id": id})

	if input.RewardIdForAddEmote != nil {
		updateBuilder = updateBuilder.Set("reward_id_for_add_emote", *input.RewardIdForAddEmote)
	}

	if input.RewardIdForRemoveEmote != nil {
		updateBuilder = updateBuilder.Set("reward_id_for_remove_emote", *input.RewardIdForRemoveEmote)
	}

	updateBuilder = updateBuilder.Set(
		"delete_emotes_only_added_by_app",
		input.DeleteEmotesOnlyAddedByApp,
	)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil

}
