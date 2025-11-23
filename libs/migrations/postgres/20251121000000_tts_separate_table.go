package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upTtsSeparateTable, downTtsSeparateTable)
}

func upTtsSeparateTable(ctx context.Context, tx *sql.Tx) error {
	tablesCreateQuery := `
CREATE TABLE channels_modules_tts (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	channel_id TEXT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
	user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	enabled BOOLEAN,
	rate INTEGER NOT NULL DEFAULT 50,
	volume INTEGER NOT NULL DEFAULT 100,
	pitch INTEGER NOT NULL DEFAULT 50,
	voice TEXT NOT NULL DEFAULT '',
	allow_users_choose_voice_in_main_command BOOLEAN NOT NULL DEFAULT false,
	max_symbols INTEGER NOT NULL DEFAULT 500,
	disallowed_voices TEXT[] NOT NULL DEFAULT '{}',
	do_not_read_emoji BOOLEAN NOT NULL DEFAULT false,
	do_not_read_twitch_emotes BOOLEAN NOT NULL DEFAULT false,
	do_not_read_links BOOLEAN NOT NULL DEFAULT false,
	read_chat_messages BOOLEAN NOT NULL DEFAULT false,
	read_chat_messages_nicknames BOOLEAN NOT NULL DEFAULT false,
	UNIQUE(channel_id, user_id)
);

CREATE INDEX IF NOT EXISTS channels_modules_tts_channel_id_idx ON channels_modules_tts (channel_id);
CREATE UNIQUE INDEX IF NOT EXISTS channels_modules_tts_channel_id_null_user_id_idx
	ON channels_modules_tts (channel_id) WHERE user_id IS NULL;
`

	if _, err := tx.ExecContext(ctx, tablesCreateQuery); err != nil {
		return fmt.Errorf("create tables: %w", err)
	}

	// Migrate existing data
	type OldTTSSettings struct {
		Enabled                            *bool    `json:"enabled"`
		Rate                               int      `json:"rate"`
		Volume                             int      `json:"volume"`
		Pitch                              int      `json:"pitch"`
		Voice                              string   `json:"voice"`
		AllowUsersChooseVoiceInMainCommand bool     `json:"allow_users_choose_voice_in_main_command"`
		MaxSymbols                         int      `json:"max_symbols"`
		DisallowedVoices                   []string `json:"disallowed_voices"`
		DoNotReadEmoji                     bool     `json:"do_not_read_emoji"`
		DoNotReadTwitchEmotes              bool     `json:"do_not_read_twitch_emotes"`
		DoNotReadLinks                     bool     `json:"do_not_read_links"`
		ReadChatMessages                   bool     `json:"read_chat_messages"`
		ReadChatMessagesNicknames          bool     `json:"read_chat_messages_nicknames"`
	}

	findQuery := `
SELECT id, settings, "channelId", "userId" FROM channels_modules_settings WHERE type = 'tts'
`

	rows, err := tx.QueryContext(ctx, findQuery)
	if err != nil {
		return fmt.Errorf("query old tts settings: %w", err)
	}
	defer rows.Close()

	type oldRow struct {
		ID        string
		Settings  []byte
		ChannelID string
		UserID    *string
	}

	var oldRows []oldRow
	for rows.Next() {
		var row oldRow
		if err := rows.Scan(&row.ID, &row.Settings, &row.ChannelID, &row.UserID); err != nil {
			return fmt.Errorf("scan old tts row: %w", err)
		}
		oldRows = append(oldRows, row)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate old tts rows: %w", err)
	}

	// Insert into new table
	insertQuery := `
INSERT INTO channels_modules_tts (
	channel_id, user_id, enabled, rate, volume, pitch, voice,
	allow_users_choose_voice_in_main_command, max_symbols, disallowed_voices,
	do_not_read_emoji, do_not_read_twitch_emotes, do_not_read_links,
	read_chat_messages, read_chat_messages_nicknames
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
`

	for _, row := range oldRows {
		var settings OldTTSSettings
		if err := json.Unmarshal(row.Settings, &settings); err != nil {
			return fmt.Errorf("unmarshal tts settings for channel %s: %w", row.ChannelID, err)
		}

		// Set defaults
		if settings.Rate == 0 {
			settings.Rate = 50
		}
		if settings.Volume == 0 {
			settings.Volume = 100
		}
		if settings.Pitch == 0 {
			settings.Pitch = 50
		}
		if settings.MaxSymbols == 0 {
			settings.MaxSymbols = 500
		}
		if settings.DisallowedVoices == nil {
			settings.DisallowedVoices = []string{}
		}

		_, err := tx.ExecContext(
			ctx,
			insertQuery,
			row.ChannelID,
			row.UserID,
			settings.Enabled,
			settings.Rate,
			settings.Volume,
			settings.Pitch,
			settings.Voice,
			settings.AllowUsersChooseVoiceInMainCommand,
			settings.MaxSymbols,
			settings.DisallowedVoices,
			settings.DoNotReadEmoji,
			settings.DoNotReadTwitchEmotes,
			settings.DoNotReadLinks,
			settings.ReadChatMessages,
			settings.ReadChatMessagesNicknames,
		)
		if err != nil {
			return fmt.Errorf("insert tts settings for channel %s: %w", row.ChannelID, err)
		}
	}

	// Delete old settings
	deleteQuery := `DELETE FROM channels_modules_settings WHERE type = 'tts'`
	if _, err := tx.ExecContext(ctx, deleteQuery); err != nil {
		return fmt.Errorf("delete old tts settings: %w", err)
	}

	return nil
}

func downTtsSeparateTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS channels_modules_tts;`)
	return err
}
