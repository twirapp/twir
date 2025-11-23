package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upTtsSeparateUsersSettings, downTtsSeparateUsersSettings)
}

func upTtsSeparateUsersSettings(ctx context.Context, tx *sql.Tx) error {
	// Create user-specific TTS settings table
	createTableQuery := `
CREATE TABLE channels_overlays_tts_users (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	channel_id TEXT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
	user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	voice TEXT NOT NULL DEFAULT '',
	rate INTEGER NOT NULL DEFAULT 50,
	pitch INTEGER NOT NULL DEFAULT 50,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	CONSTRAINT channels_overlays_tts_users_channel_user_unique UNIQUE(channel_id, user_id)
);

CREATE INDEX IF NOT EXISTS channels_overlays_tts_users_channel_id_idx ON channels_overlays_tts_users (channel_id);
CREATE INDEX IF NOT EXISTS channels_overlays_tts_users_user_id_idx ON channels_overlays_tts_users (user_id);
CREATE INDEX IF NOT EXISTS channels_overlays_tts_users_channel_user_idx ON channels_overlays_tts_users (channel_id, user_id);
`

	if _, err := tx.ExecContext(ctx, createTableQuery); err != nil {
		return fmt.Errorf("create table: %w", err)
	}

	// Define old TTS user settings structure
	type OldUserTTSSettings struct {
		Voice string `json:"voice"`
		Rate  int    `json:"rate"`
		Pitch int    `json:"pitch"`
	}

	// Query user-specific settings from channels_modules_settings
	findQuery := `
SELECT id, settings, "channelId", "userId"
FROM channels_modules_settings
WHERE type = 'tts' AND "userId" IS NOT NULL
`

	rows, err := tx.QueryContext(ctx, findQuery)
	if err != nil {
		return fmt.Errorf("query user tts settings: %w", err)
	}
	defer rows.Close()

	type migrationData struct {
		id        string
		channelID string
		userID    string
		settings  OldUserTTSSettings
	}

	data := make([]migrationData, 0)
	for rows.Next() {
		var id string
		var channelID string
		var userID string
		var settingsRaw []byte

		if err := rows.Scan(&id, &settingsRaw, &channelID, &userID); err != nil {
			return fmt.Errorf("scan row: %w", err)
		}

		var settings OldUserTTSSettings
		if err := json.Unmarshal(settingsRaw, &settings); err != nil {
			return fmt.Errorf("unmarshal settings for user %s: %w", userID, err)
		}

		data = append(
			data, migrationData{
				id:        id,
				channelID: channelID,
				userID:    userID,
				settings:  settings,
			},
		)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows iteration: %w", err)
	}

	// Insert user settings into new table
	insertQuery := `
INSERT INTO channels_overlays_tts_users (channel_id, user_id, voice, rate, pitch)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (channel_id, user_id) DO NOTHING;
`

	for _, d := range data {
		// Set defaults for missing values
		voice := d.settings.Voice
		if voice == "" {
			voice = ""
		}

		rate := d.settings.Rate
		if rate == 0 {
			rate = 50
		}

		pitch := d.settings.Pitch
		if pitch == 0 {
			pitch = 50
		}

		if _, err := tx.ExecContext(
			ctx,
			insertQuery,
			d.channelID,
			d.userID,
			voice,
			rate,
			pitch,
		); err != nil {
			return fmt.Errorf(
				"insert user settings for channel %s user %s: %w",
				d.channelID,
				d.userID,
				err,
			)
		}
	}

	// Delete migrated user-specific settings from channels_modules_settings
	deleteQuery := `
DELETE FROM channels_modules_settings
WHERE type = 'tts' AND "userId" IS NOT NULL;
`

	if _, err := tx.ExecContext(ctx, deleteQuery); err != nil {
		return fmt.Errorf("delete old user settings: %w", err)
	}

	return nil
}

func downTtsSeparateUsersSettings(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS channels_overlays_tts_users;`)
	return err
}
