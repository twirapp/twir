package pgx

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"github.com/twirapp/twir/libs/repositories"
	"github.com/twirapp/twir/libs/repositories/giveaways_participants"
	"github.com/twirapp/twir/libs/repositories/giveaways_participants/model"
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

var _ giveaways_participants.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (p *Pgx) ResetWinners(
	ctx context.Context, participantsIds ...ulid.ULID,
) error {
	query := `
UPDATE channels_giveaways_participants
SET is_winner = FALSE
WHERE id = ANY($1)
	`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)

	_, err := conn.Exec(ctx, query, participantsIds)

	return err
}

func (p *Pgx) Create(
	ctx context.Context,
	input giveaways_participants.CreateInput,
) (model.ChannelGiveawayParticipant, error) {
	query := `
INSERT INTO channels_giveaways_participants("giveaway_id", "display_name", "user_id") VALUES ($1, $2, $3)
RETURNING id, giveaway_id, is_winner, display_name, user_id
	`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(
		ctx,
		query,
		input.GiveawayID,
		input.DisplayName,
		input.UserID,
	)
	if err != nil {
		return model.ChannelGiveawayParticipantNil, err
	}

	result, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model.ChannelGiveawayParticipant],
	)
	if err != nil {
		return model.ChannelGiveawayParticipantNil, err
	}

	return result, nil
}

func (p *Pgx) Delete(ctx context.Context, id ulid.ULID) error {
	query := `
DELETE FROM channels_giveaways_participants WHERE id = $1
	`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if rows.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row to be affected, got %d", rows.RowsAffected())
	}

	return nil
}

func (p *Pgx) GetByID(ctx context.Context, id ulid.ULID) (model.ChannelGiveawayParticipant, error) {
	query := `
SELECT id, giveaway_id, is_winner, display_name, user_id FROM channels_giveaways_participants WHERE id = $1
	`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, id)
	if err != nil {
		return model.ChannelGiveawayParticipantNil, err
	}

	result, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model.ChannelGiveawayParticipant],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ChannelGiveawayParticipantNil, giveaways_participants.ErrNotFound
		}

		return model.ChannelGiveawayParticipantNil, err
	}

	return result, nil
}

func (p *Pgx) GetManyByGiveawayID(
	ctx context.Context,
	giveawayID ulid.ULID,
) ([]model.ChannelGiveawayParticipant, error) {
	selectBuilder := sq.Select().
		From("channels_giveaways_participants").
		Where(squirrel.Eq{`"giveaway_id"`: giveawayID})

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.ChannelGiveawayParticipant])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *Pgx) GetWinnerForGiveaway(
	ctx context.Context,
	giveawayID ulid.ULID,
) (model.ChannelGiveawayParticipant, error) {
	query := `
SELECT id, giveaway_id, is_winner, display_name, user_id FROM channels_giveaways_participants WHERE giveaway_id = $1 AND is_winner = true
	`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, giveawayID)
	if err != nil {
		return model.ChannelGiveawayParticipantNil, err
	}

	result, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model.ChannelGiveawayParticipant],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ChannelGiveawayParticipantNil, giveaways_participants.ErrNotFound
		}

		return model.ChannelGiveawayParticipantNil, err
	}

	return result, nil
}

func (p *Pgx) Update(
	ctx context.Context,
	id ulid.ULID,
	input giveaways_participants.UpdateInput,
) (model.ChannelGiveawayParticipant, error) {
	updateBuilder := sq.Update("channels_giveaways").
		Where(squirrel.Eq{"id": id}).
		Suffix(`RETURNING id, giveaway_id, is_winner, display_name, user_id`)
	updateBuilder = repositories.SquirrelApplyPatch(
		updateBuilder,
		map[string]any{
			"is_winner": input.IsWinner,
		},
	)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return model.ChannelGiveawayParticipantNil, err
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return model.ChannelGiveawayParticipantNil, err
	}

	result, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model.ChannelGiveawayParticipant],
	)
	if err != nil {
		return model.ChannelGiveawayParticipantNil, err
	}

	return result, nil
}
