package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upObsWebsocketSeparateTable, downObsWebsocketSeparateTable)
}

func upObsWebsocketSeparateTable(ctx context.Context, tx *sql.Tx) error {
	type obsWebsocketSettings20251203000000 struct {
		ServerPort     int    `json:"serverPort"`
		ServerAddress  string `json:"serverAddress"`
		ServerPassword string `json:"serverPassword"`
	}

	type channelModuleSettings20251203000000 struct {
		ID        string
		ChannelID string
		Type      string
		Settings  []byte
	}

	// Create new table
	createTableQuery := `
CREATE TABLE IF NOT EXISTS channels_modules_obs_websocket (
	id SERIAL PRIMARY KEY,
	channel_id TEXT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
	server_port INTEGER NOT NULL DEFAULT 4455,
	server_address TEXT NOT NULL DEFAULT '',
	server_password TEXT NOT NULL DEFAULT '',
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS channels_modules_obs_websocket_channel_id_idx ON channels_modules_obs_websocket(channel_id);
`
	if _, err := tx.ExecContext(ctx, createTableQuery); err != nil {
		return fmt.Errorf("create table: %w", err)
	}

	// Get all obs_websocket modules from old table
	rows, err := tx.QueryContext(
		ctx,
		`SELECT id, "channelId", type, settings FROM channels_modules_settings WHERE type = 'obs_websocket'`,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("get channels modules settings: %w", err)
	}

	if rows == nil {
		return nil
	}
	defer rows.Close()

	moduleSettings := []channelModuleSettings20251203000000{}
	for rows.Next() {
		module := channelModuleSettings20251203000000{}
		err := rows.Scan(
			&module.ID,
			&module.ChannelID,
			&module.Type,
			&module.Settings,
		)
		if err != nil {
			return fmt.Errorf("scan row: %w", err)
		}

		moduleSettings = append(moduleSettings, module)
	}

	if rows.Err() != nil {
		return fmt.Errorf("rows error: %w", rows.Err())
	}

	// Migrate data
	for _, module := range moduleSettings {
		var settings obsWebsocketSettings20251203000000
		if len(module.Settings) > 0 {
			err := json.Unmarshal(module.Settings, &settings)
			if err != nil {
				// Skip if we can't parse the settings
				continue
			}
		}

		insertQuery := `
INSERT INTO channels_modules_obs_websocket (
	channel_id,
	server_port,
	server_address,
	server_password
)
VALUES ($1, $2, $3, $4)
ON CONFLICT (channel_id) DO NOTHING
`

		_, err := tx.ExecContext(
			ctx,
			insertQuery,
			module.ChannelID,
			settings.ServerPort,
			settings.ServerAddress,
			settings.ServerPassword,
		)
		if err != nil {
			return fmt.Errorf("insert obs websocket settings for channel %s: %w", module.ChannelID, err)
		}
	}

	// Delete old records
	_, err = tx.ExecContext(
		ctx,
		`DELETE FROM channels_modules_settings WHERE type = 'obs_websocket'`,
	)
	if err != nil {
		return fmt.Errorf("delete old obs websocket settings: %w", err)
	}

	return nil
}

func downObsWebsocketSeparateTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "DROP TABLE IF EXISTS channels_modules_obs_websocket")
	return err
}
