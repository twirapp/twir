package pgx

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
	"github.com/twirapp/twir/libs/repositories/modules_tts"
	"github.com/twirapp/twir/libs/repositories/modules_tts/model"
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

var _ modules_tts.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) scanRow(rows pgx.Rows) (model.TTS, error) {
	var tts model.TTS
	err := rows.Scan(
		&tts.ID,
		&tts.ChannelID,
		&tts.UserID,
		&tts.CreatedAt,
		&tts.UpdatedAt,
		&tts.Enabled,
		&tts.Rate,
		&tts.Volume,
		&tts.Pitch,
		&tts.Voice,
		&tts.AllowUsersChooseVoiceInMainCommand,
		&tts.MaxSymbols,
		&tts.DisallowedVoices,
		&tts.DoNotReadEmoji,
		&tts.DoNotReadTwitchEmotes,
		&tts.DoNotReadLinks,
		&tts.ReadChatMessages,
		&tts.ReadChatMessagesNicknames,
	)
	if err != nil {
		return model.Nil, err
	}
	return tts, nil
}

// Channel-level operations
func (c *Pgx) GetByChannelID(ctx context.Context, channelID string) (model.TTS, error) {
	query := `
SELECT
	id, channel_id, user_id, created_at, updated_at, enabled, rate, volume, pitch, voice,
	allow_users_choose_voice_in_main_command, max_symbols, disallowed_voices,
	do_not_read_emoji, do_not_read_twitch_emotes, do_not_read_links,
	read_chat_messages, read_chat_messages_nicknames
FROM channels_modules_tts
WHERE channel_id = $1 AND user_id IS NULL
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, channelID)
	if err != nil {
		return model.Nil, fmt.Errorf("query tts by channel id: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return model.Nil, err
		}
		return model.Nil, modules_tts.ErrNotFound
	}

	return c.scanRow(rows)
}

func (c *Pgx) CreateForChannel(
	ctx context.Context,
	input modules_tts.CreateInput,
) (model.TTS, error) {
	query := `
INSERT INTO channels_modules_tts (
	channel_id, user_id, enabled, rate, volume, pitch, voice,
	allow_users_choose_voice_in_main_command, max_symbols, disallowed_voices,
	do_not_read_emoji, do_not_read_twitch_emotes, do_not_read_links,
	read_chat_messages, read_chat_messages_nicknames
) VALUES ($1, NULL, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING
	id, channel_id, user_id, created_at, updated_at, enabled, rate, volume, pitch, voice,
	allow_users_choose_voice_in_main_command, max_symbols, disallowed_voices,
	do_not_read_emoji, do_not_read_twitch_emotes, do_not_read_links,
	read_chat_messages, read_chat_messages_nicknames
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(
		ctx,
		query,
		input.ChannelID,
		input.Enabled,
		input.Rate,
		input.Volume,
		input.Pitch,
		input.Voice,
		input.AllowUsersChooseVoiceInMainCommand,
		input.MaxSymbols,
		pq.Array(input.DisallowedVoices),
		input.DoNotReadEmoji,
		input.DoNotReadTwitchEmotes,
		input.DoNotReadLinks,
		input.ReadChatMessages,
		input.ReadChatMessagesNicknames,
	)
	if err != nil {
		return model.Nil, fmt.Errorf("create tts for channel: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return model.Nil, fmt.Errorf("no rows returned after insert")
	}

	return c.scanRow(rows)
}

func (c *Pgx) UpdateForChannel(
	ctx context.Context,
	channelID string,
	input modules_tts.UpdateInput,
) (model.TTS, error) {
	query := `
UPDATE channels_modules_tts
SET
	enabled = $2,
	rate = $3,
	volume = $4,
	pitch = $5,
	voice = $6,
	allow_users_choose_voice_in_main_command = $7,
	max_symbols = $8,
	disallowed_voices = $9,
	do_not_read_emoji = $10,
	do_not_read_twitch_emotes = $11,
	do_not_read_links = $12,
	read_chat_messages = $13,
	read_chat_messages_nicknames = $14,
	updated_at = now()
WHERE channel_id = $1 AND user_id IS NULL
RETURNING
	id, channel_id, user_id, created_at, updated_at, enabled, rate, volume, pitch, voice,
	allow_users_choose_voice_in_main_command, max_symbols, disallowed_voices,
	do_not_read_emoji, do_not_read_twitch_emotes, do_not_read_links,
	read_chat_messages, read_chat_messages_nicknames
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(
		ctx,
		query,
		channelID,
		input.Enabled,
		input.Rate,
		input.Volume,
		input.Pitch,
		input.Voice,
		input.AllowUsersChooseVoiceInMainCommand,
		input.MaxSymbols,
		pq.Array(input.DisallowedVoices),
		input.DoNotReadEmoji,
		input.DoNotReadTwitchEmotes,
		input.DoNotReadLinks,
		input.ReadChatMessages,
		input.ReadChatMessagesNicknames,
	)
	if err != nil {
		return model.Nil, fmt.Errorf("update tts for channel: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return model.Nil, err
		}
		return model.Nil, modules_tts.ErrNotFound
	}

	return c.scanRow(rows)
}

func (c *Pgx) DeleteForChannel(ctx context.Context, channelID string) error {
	query := `DELETE FROM channels_modules_tts WHERE channel_id = $1 AND user_id IS NULL`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, channelID)
	if err != nil {
		return fmt.Errorf("delete tts for channel: %w", err)
	}

	return nil
}

// User-specific operations
func (c *Pgx) GetByChannelIDAndUserID(
	ctx context.Context,
	channelID, userID string,
) (model.TTS, error) {
	query := `
SELECT
	id, channel_id, user_id, created_at, updated_at, enabled, rate, volume, pitch, voice,
	allow_users_choose_voice_in_main_command, max_symbols, disallowed_voices,
	do_not_read_emoji, do_not_read_twitch_emotes, do_not_read_links,
	read_chat_messages, read_chat_messages_nicknames
FROM channels_modules_tts
WHERE channel_id = $1 AND user_id = $2
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, channelID, userID)
	if err != nil {
		return model.Nil, fmt.Errorf("query tts by channel and user id: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return model.Nil, err
		}
		return model.Nil, modules_tts.ErrNotFound
	}

	return c.scanRow(rows)
}

func (c *Pgx) GetAllUsersByChannelID(ctx context.Context, channelID string) ([]model.TTS, error) {
	query := `
SELECT
	id, channel_id, user_id, created_at, updated_at, enabled, rate, volume, pitch, voice,
	allow_users_choose_voice_in_main_command, max_symbols, disallowed_voices,
	do_not_read_emoji, do_not_read_twitch_emotes, do_not_read_links,
	read_chat_messages, read_chat_messages_nicknames
FROM channels_modules_tts
WHERE channel_id = $1 AND user_id IS NOT NULL
ORDER BY created_at DESC
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, channelID)
	if err != nil {
		return nil, fmt.Errorf("query all user tts by channel id: %w", err)
	}
	defer rows.Close()

	var result []model.TTS
	for rows.Next() {
		tts, err := c.scanRow(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, tts)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Pgx) CreateForUser(
	ctx context.Context,
	input modules_tts.CreateInput,
) (model.TTS, error) {
	if input.UserID == nil {
		return model.Nil, fmt.Errorf("user id is required")
	}

	query := `
INSERT INTO channels_modules_tts (
	channel_id, user_id, enabled, rate, volume, pitch, voice,
	allow_users_choose_voice_in_main_command, max_symbols, disallowed_voices,
	do_not_read_emoji, do_not_read_twitch_emotes, do_not_read_links,
	read_chat_messages, read_chat_messages_nicknames
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING
	id, channel_id, user_id, created_at, updated_at, enabled, rate, volume, pitch, voice,
	allow_users_choose_voice_in_main_command, max_symbols, disallowed_voices,
	do_not_read_emoji, do_not_read_twitch_emotes, do_not_read_links,
	read_chat_messages, read_chat_messages_nicknames
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(
		ctx,
		query,
		input.ChannelID,
		*input.UserID,
		input.Enabled,
		input.Rate,
		input.Volume,
		input.Pitch,
		input.Voice,
		input.AllowUsersChooseVoiceInMainCommand,
		input.MaxSymbols,
		pq.Array(input.DisallowedVoices),
		input.DoNotReadEmoji,
		input.DoNotReadTwitchEmotes,
		input.DoNotReadLinks,
		input.ReadChatMessages,
		input.ReadChatMessagesNicknames,
	)
	if err != nil {
		return model.Nil, fmt.Errorf("create tts for user: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return model.Nil, fmt.Errorf("no rows returned after insert")
	}

	return c.scanRow(rows)
}

func (c *Pgx) UpdateForUser(
	ctx context.Context,
	channelID, userID string,
	input modules_tts.UpdateInput,
) (model.TTS, error) {
	query := `
UPDATE channels_modules_tts
SET
	enabled = $3,
	rate = $4,
	volume = $5,
	pitch = $6,
	voice = $7,
	allow_users_choose_voice_in_main_command = $8,
	max_symbols = $9,
	disallowed_voices = $10,
	do_not_read_emoji = $11,
	do_not_read_twitch_emotes = $12,
	do_not_read_links = $13,
	read_chat_messages = $14,
	read_chat_messages_nicknames = $15,
	updated_at = now()
WHERE channel_id = $1 AND user_id = $2
RETURNING
	id, channel_id, user_id, created_at, updated_at, enabled, rate, volume, pitch, voice,
	allow_users_choose_voice_in_main_command, max_symbols, disallowed_voices,
	do_not_read_emoji, do_not_read_twitch_emotes, do_not_read_links,
	read_chat_messages, read_chat_messages_nicknames
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(
		ctx,
		query,
		channelID,
		userID,
		input.Enabled,
		input.Rate,
		input.Volume,
		input.Pitch,
		input.Voice,
		input.AllowUsersChooseVoiceInMainCommand,
		input.MaxSymbols,
		pq.Array(input.DisallowedVoices),
		input.DoNotReadEmoji,
		input.DoNotReadTwitchEmotes,
		input.DoNotReadLinks,
		input.ReadChatMessages,
		input.ReadChatMessagesNicknames,
	)
	if err != nil {
		return model.Nil, fmt.Errorf("update tts for user: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return model.Nil, err
		}
		return model.Nil, modules_tts.ErrNotFound
	}

	return c.scanRow(rows)
}

func (c *Pgx) DeleteForUser(ctx context.Context, channelID, userID string) error {
	query := `DELETE FROM channels_modules_tts WHERE channel_id = $1 AND user_id = $2`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, channelID, userID)
	if err != nil {
		return fmt.Errorf("delete tts for user: %w", err)
	}

	return nil
}

func (c *Pgx) DeleteUsersForChannel(
	ctx context.Context,
	channelID string,
	userIDs []string,
) error {
	if len(userIDs) == 0 {
		return nil
	}

	query, args, err := sq.Delete("channels_modules_tts").
		Where(squirrel.Eq{"channel_id": channelID, "user_id": userIDs}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build delete query: %w", err)
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete tts users for channel: %w", err)
	}

	return nil
}
