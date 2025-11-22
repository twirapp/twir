package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upTTSSeparateTable, downTTSSeparateTable)
}

func upTTSSeparateTable(ctx context.Context, tx *sql.Tx) error {
	tablesCreateQuery := `
CREATE TABLE channels_overlays_tts (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	channel_id TEXT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	enabled BOOLEAN NOT NULL DEFAULT false,
	voice TEXT NOT NULL DEFAULT 'alan',
	disallowed_voices TEXT[] NOT NULL DEFAULT '{}',
	pitch INTEGER NOT NULL DEFAULT 50,
	rate INTEGER NOT NULL DEFAULT 50,
	volume INTEGER NOT NULL DEFAULT 30,
	do_not_read_twitch_emotes BOOLEAN NOT NULL DEFAULT true,
	do_not_read_emoji BOOLEAN NOT NULL DEFAULT true,
	do_not_read_links BOOLEAN NOT NULL DEFAULT true,
	allow_users_choose_voice_in_main_command BOOLEAN NOT NULL DEFAULT false,
	max_symbols INTEGER NOT NULL DEFAULT 0,
	read_chat_messages BOOLEAN NOT NULL DEFAULT false,
	read_chat_messages_nicknames BOOLEAN NOT NULL DEFAULT false
);

CREATE UNIQUE INDEX IF NOT EXISTS channels_overlays_tts_channel_id_unique ON channels_overlays_tts (channel_id);
`

	if _, err := tx.ExecContext(ctx, tablesCreateQuery); err != nil {
		return fmt.Errorf("create tables: %w", err)
	}

	type OldTTSSettings struct {
		Enabled                          bool     `json:"enabled"`
		Voice                            string   `json:"voice"`
		DisallowedVoices                 []string `json:"disallowedVoices"`
		Pitch                            int32    `json:"pitch"`
		Rate                             int32    `json:"rate"`
		Volume                           int32    `json:"volume"`
		DoNotReadTwitchEmotes            bool     `json:"doNotReadTwitchEmotes"`
		DoNotReadEmoji                   bool     `json:"doNotReadEmoji"`
		DoNotReadLinks                   bool     `json:"doNotReadLinks"`
		AllowUsersChooseVoiceInMainCommand bool     `json:"allowUsersChooseVoiceInMainCommand"`
		MaxSymbols                       int32    `json:"maxSymbols"`
		ReadChatMessages                 bool     `json:"readChatMessages"`
		ReadChatMessagesNicknames        bool     `json:"readChatMessagesNicknames"`
	}

	findQuery := `
SELECT id, settings, "channelId" FROM channels_modules_settings WHERE type = 'tts' AND "userId" IS NULL
`

	rows, err := tx.QueryContext(ctx, findQuery)
	if err != nil {
		return fmt.Errorf("find query: %w", err)
	}
	defer rows.Close()

	type migrationData struct {
		id        string
		channelID string
		settings  OldTTSSettings
	}

	data := make([]migrationData, 0)
	for rows.Next() {
		var id string
		var channelID string
		var settingsRaw []byte

		if err := rows.Scan(&id, &settingsRaw, &channelID); err != nil {
			return fmt.Errorf("scan: %w", err)
		}

		var settings OldTTSSettings
		if err := json.Unmarshal(settingsRaw, &settings); err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}

		data = append(
			data, migrationData{
				id:        id,
				channelID: channelID,
				settings:  settings,
			},
		)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows: %w", err)
	}

	for _, d := range data {
		insertQuery := `
INSERT INTO channels_overlays_tts (
	channel_id,
	enabled,
	voice,
	disallowed_voices,
	pitch,
	rate,
	volume,
	do_not_read_twitch_emotes,
	do_not_read_emoji,
	do_not_read_links,
	allow_users_choose_voice_in_main_command,
	max_symbols,
	read_chat_messages,
	read_chat_messages_nicknames
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);
		`

		if _, err := tx.ExecContext(
			ctx,
			insertQuery,
			d.channelID,
			d.settings.Enabled,
			d.settings.Voice,
			d.settings.DisallowedVoices,
			d.settings.Pitch,
			d.settings.Rate,
			d.settings.Volume,
			d.settings.DoNotReadTwitchEmotes,
			d.settings.DoNotReadEmoji,
			d.settings.DoNotReadLinks,
			d.settings.AllowUsersChooseVoiceInMainCommand,
			d.settings.MaxSymbols,
			d.settings.ReadChatMessages,
			d.settings.ReadChatMessagesNicknames,
		); err != nil {
			return fmt.Errorf("insert: %w", err)
		}
	}

	deleteQuery := `
DELETE FROM channels_modules_settings WHERE type = 'tts' AND "userId" IS NULL
	`

	if _, err := tx.ExecContext(ctx, deleteQuery); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func downTTSSeparateTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}

