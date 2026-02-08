package pgx

import (
	"context"
	"errors"
	"strings"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	discordsendednotifications "github.com/twirapp/twir/libs/repositories/discord_sended_notifications"
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

var _ discordsendednotifications.Repository = (*Pgx)(nil)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

var selectColumns = []string{
	"id",
	"guild_id",
	"message_id",
	"channel_id",
	"discord_channel_id",
	"created_at",
	"updated_at",
	"last_updated_at",
}

var selectColumnsStr = strings.Join(selectColumns, ", ")

func (p *Pgx) CreateMany(ctx context.Context, input []discordsendednotifications.CreateInput) error {
	batch := &pgx.Batch{}

	for _, i := range input {
		batch.Queue(
			"INSERT INTO discord_sended_notifications (guild_id, message_id, channel_id, discord_channel_id) VALUES ($1, $2, $3, $4)",
			i.GuildID,
			i.MessageID,
			i.TwitchChannelID,
			i.DiscordChannelID,
		)
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)

	br := conn.SendBatch(ctx, batch)
	defer br.Close()

	for range input {
		_, err := br.Exec()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Pgx) Create(ctx context.Context, input discordsendednotifications.CreateInput) error {
	query := `
INSERT INTO discord_sended_notifications (
	guild_id,
	message_id,
	channel_id,
	discord_channel_id
)
VALUES (@guild_id, @message_id, @channel_id, @discord_channel_id)
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	_, err := conn.Exec(
		ctx,
		query,
		pgx.NamedArgs{
			"guild_id":           input.GuildID,
			"message_id":         input.MessageID,
			"channel_id":         input.TwitchChannelID,
			"discord_channel_id": input.DiscordChannelID,
		},
	)
	return err
}

func (p *Pgx) GetByMessageID(
	ctx context.Context,
	messageID string,
) (discordsendednotifications.DiscordSendedNotification, error) {
	query := `
SELECT ` + selectColumnsStr + `
FROM discord_sended_notifications
WHERE message_id = @message_id
LIMIT 1
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, pgx.NamedArgs{"message_id": messageID})
	if err != nil {
		return discordsendednotifications.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[discordsendednotifications.DiscordSendedNotification],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return discordsendednotifications.Nil, nil
		}
		return discordsendednotifications.Nil, err
	}

	return result, nil
}

func (p *Pgx) GetByChannelID(
	ctx context.Context,
	channelID string,
) ([]discordsendednotifications.DiscordSendedNotification, error) {
	query := `
SELECT ` + selectColumnsStr + `
FROM discord_sended_notifications
WHERE channel_id = @channel_id
ORDER BY created_at
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, pgx.NamedArgs{"channel_id": channelID})
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByName[discordsendednotifications.DiscordSendedNotification],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (p *Pgx) GetByGuildID(
	ctx context.Context,
	guildID string,
) ([]discordsendednotifications.DiscordSendedNotification, error) {
	query := `
SELECT ` + selectColumnsStr + `
FROM discord_sended_notifications
WHERE guild_id = @guild_id
ORDER BY created_at
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, pgx.NamedArgs{"guild_id": guildID})
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByName[discordsendednotifications.DiscordSendedNotification],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (p *Pgx) GetAll(ctx context.Context) ([]discordsendednotifications.DiscordSendedNotification, error) {
	query := `
SELECT ` + selectColumnsStr + `
FROM discord_sended_notifications
ORDER BY created_at
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByName[discordsendednotifications.DiscordSendedNotification],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (p *Pgx) DeleteByMessageID(ctx context.Context, messageID string) error {
	query := `DELETE FROM discord_sended_notifications WHERE message_id = @message_id`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	_, err := conn.Exec(ctx, query, pgx.NamedArgs{"message_id": messageID})
	return err
}

func (p *Pgx) DeleteByChannelID(ctx context.Context, channelID string) error {
	query := `DELETE FROM discord_sended_notifications WHERE channel_id = @channel_id`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	_, err := conn.Exec(ctx, query, pgx.NamedArgs{"channel_id": channelID})
	return err
}

func (p *Pgx) DeleteByGuildID(ctx context.Context, guildID string) error {
	query := `DELETE FROM discord_sended_notifications WHERE guild_id = @guild_id`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	_, err := conn.Exec(ctx, query, pgx.NamedArgs{"guild_id": guildID})
	return err
}

func (p *Pgx) UpdateLastUpdatedAt(ctx context.Context, messageID string) error {
	query := `UPDATE discord_sended_notifications SET last_updated_at = now() WHERE message_id = @message_id`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	_, err := conn.Exec(ctx, query, pgx.NamedArgs{"message_id": messageID})
	return err
}
