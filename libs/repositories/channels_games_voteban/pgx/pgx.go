package pgx

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	votebanentity "github.com/twirapp/twir/libs/entities/voteban"
	"github.com/twirapp/twir/libs/repositories"
	"github.com/twirapp/twir/libs/repositories/channels_games_voteban"
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

var (
	_  channels_games_voteban.Repository = (*Pgx)(nil)
	sq                                   = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

const selectColumns = `id, channel_id, enabled, timeout_seconds, timeout_moderators,
	init_message, ban_message, ban_message_moderators, survive_message, survive_message_moderators,
	needed_votes, vote_duration, voting_mode, chat_votes_words_positive, chat_votes_words_negative`

type scanModel struct {
	ID                       uuid.UUID `db:"id"`
	ChannelID                string    `db:"channel_id"`
	Enabled                  bool      `db:"enabled"`
	TimeoutSeconds           int       `db:"timeout_seconds"`
	TimeoutModerators        bool      `db:"timeout_moderators"`
	InitMessage              string    `db:"init_message"`
	BanMessage               string    `db:"ban_message"`
	BanMessageModerators     string    `db:"ban_message_moderators"`
	SurviveMessage           string    `db:"survive_message"`
	SurviveMessageModerators string    `db:"survive_message_moderators"`
	NeededVotes              int       `db:"needed_votes"`
	VoteDuration             int       `db:"vote_duration"`
	VotingMode               string    `db:"voting_mode"`
	ChatVotesWordsPositive   []string  `db:"chat_votes_words_positive"`
	ChatVotesWordsNegative   []string  `db:"chat_votes_words_negative"`
}

func (c scanModel) toEntity() votebanentity.Voteban {
	return votebanentity.Voteban{
		ID:                       c.ID,
		ChannelID:                c.ChannelID,
		Enabled:                  c.Enabled,
		TimeoutSeconds:           c.TimeoutSeconds,
		TimeoutModerators:        c.TimeoutModerators,
		InitMessage:              c.InitMessage,
		BanMessage:               c.BanMessage,
		BanMessageModerators:     c.BanMessageModerators,
		SurviveMessage:           c.SurviveMessage,
		SurviveMessageModerators: c.SurviveMessageModerators,
		NeededVotes:              c.NeededVotes,
		VoteDuration:             c.VoteDuration,
		VotingMode:               votebanentity.VotingMode(c.VotingMode),
		ChatVotesWordsPositive:   c.ChatVotesWordsPositive,
		ChatVotesWordsNegative:   c.ChatVotesWordsNegative,
	}
}

func (c *Pgx) GetByChannelID(ctx context.Context, channelID string) (votebanentity.Voteban, error) {
	query := `
SELECT ` + selectColumns + `
FROM channels_games_voteban
WHERE channel_id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, channelID)
	if err != nil {
		return votebanentity.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return votebanentity.Nil, channels_games_voteban.ErrNotFound
		}
		return votebanentity.Nil, err
	}

	return result.toEntity(), nil
}

func (c *Pgx) GetOrCreateByChannelID(
	ctx context.Context,
	channelID string,
	input channels_games_voteban.CreateInput,
) (votebanentity.Voteban, error) {
	result, err := c.GetByChannelID(ctx, channelID)
	if err == nil {
		return result, nil
	}

	if !errors.Is(err, channels_games_voteban.ErrNotFound) {
		return votebanentity.Nil, err
	}

	insertBuilder := sq.
		Insert("channels_games_voteban").
		SetMap(
			map[string]any{
				"channel_id":                 channelID,
				"enabled":                    input.Enabled,
				"timeout_seconds":            input.TimeoutSeconds,
				"timeout_moderators":         input.TimeoutModerators,
				"init_message":               input.InitMessage,
				"ban_message":                input.BanMessage,
				"ban_message_moderators":     input.BanMessageModerators,
				"survive_message":            input.SurviveMessage,
				"survive_message_moderators": input.SurviveMessageModerators,
				"needed_votes":               input.NeededVotes,
				"vote_duration":              input.VoteDuration,
				"voting_mode":                input.VotingMode,
				"chat_votes_words_positive":  input.ChatVotesWordsPositive,
				"chat_votes_words_negative":  input.ChatVotesWordsNegative,
			},
		).
		Suffix("RETURNING " + selectColumns)

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		return votebanentity.Nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return votebanentity.Nil, err
	}

	dbResult, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		return votebanentity.Nil, err
	}

	return dbResult.toEntity(), nil
}

func (c *Pgx) Update(
	ctx context.Context,
	id uuid.UUID,
	input channels_games_voteban.UpdateInput,
) (votebanentity.Voteban, error) {
	updateBuilder := sq.
		Update("channels_games_voteban").
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING " + selectColumns)

	updateBuilder = repositories.SquirrelApplyPatch(
		updateBuilder,
		map[string]any{
			"enabled":                    input.Enabled,
			"timeout_seconds":            input.TimeoutSeconds,
			"timeout_moderators":         input.TimeoutModerators,
			"init_message":               input.InitMessage,
			"ban_message":                input.BanMessage,
			"ban_message_moderators":     input.BanMessageModerators,
			"survive_message":            input.SurviveMessage,
			"survive_message_moderators": input.SurviveMessageModerators,
			"needed_votes":               input.NeededVotes,
			"vote_duration":              input.VoteDuration,
			"voting_mode":                input.VotingMode,
		},
	)

	// Handle string arrays separately since they're not pointers
	if input.ChatVotesWordsPositive != nil {
		updateBuilder = updateBuilder.Set("chat_votes_words_positive", input.ChatVotesWordsPositive)
	}
	if input.ChatVotesWordsNegative != nil {
		updateBuilder = updateBuilder.Set("chat_votes_words_negative", input.ChatVotesWordsNegative)
	}

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return votebanentity.Nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return votebanentity.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return votebanentity.Nil, channels_games_voteban.ErrNotFound
		}
		return votebanentity.Nil, err
	}

	return result.toEntity(), nil
}
