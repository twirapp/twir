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
	goose.AddMigrationContext(upLastfmSeparateTable, downLastfmSeparateTable)
}

func upLastfmSeparateTable(ctx context.Context, tx *sql.Tx) error {
	type lastfmSeparateTableIntegrations20251202000000 struct {
		ID string
	}

	type lastfmSeparateTableChannelsIntegrations20251202000000 struct {
		ID        string
		ChannelID string
		Enabled   bool
		APIKey    *string // This is the session key for LastFM
		Data      []byte
	}

	type lastfmSeparateTableData20251202000000 struct {
		UserName *string `json:"username,omitempty"`
		Avatar   *string `json:"avatar,omitempty"`
	}

	// Get lastfm integration ID
	lastfmIntegration := lastfmSeparateTableIntegrations20251202000000{}
	err := tx.QueryRowContext(
		ctx,
		"SELECT id FROM integrations WHERE service = 'LASTFM'",
	).Scan(&lastfmIntegration.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("get lastfm integration: %w", err)
	}

	// Create new table first
	createTableQuery := `
CREATE TABLE IF NOT EXISTS channels_integrations_lastfm (
	id SERIAL PRIMARY KEY,
	channel_id TEXT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
	enabled BOOLEAN NOT NULL DEFAULT false,
	session_key TEXT,
	username TEXT,
	avatar TEXT,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS channels_integrations_lastfm_channel_id_idx ON channels_integrations_lastfm(channel_id);
`
	if _, err := tx.ExecContext(ctx, createTableQuery); err != nil {
		return fmt.Errorf("create table: %w", err)
	}

	// If no integration found, just exit
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	// Get all lastfm integrations from old table
	rows, err := tx.QueryContext(
		ctx,
		`SELECT id, "channelId", enabled, "apiKey", data FROM channels_integrations WHERE "integrationId" = $1`,
		lastfmIntegration.ID,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("get channels integrations: %w", err)
	}

	if rows == nil {
		return nil
	}
	defer rows.Close()

	channelsIntegrations := []lastfmSeparateTableChannelsIntegrations20251202000000{}
	for rows.Next() {
		channelIntegration := lastfmSeparateTableChannelsIntegrations20251202000000{}
		err := rows.Scan(
			&channelIntegration.ID,
			&channelIntegration.ChannelID,
			&channelIntegration.Enabled,
			&channelIntegration.APIKey,
			&channelIntegration.Data,
		)
		if err != nil {
			return fmt.Errorf("scan row: %w", err)
		}

		channelsIntegrations = append(channelsIntegrations, channelIntegration)
	}

	if rows.Err() != nil {
		return fmt.Errorf("rows error: %w", rows.Err())
	}

	// Migrate data
	for _, channelIntegration := range channelsIntegrations {
		var parsedData lastfmSeparateTableData20251202000000
		if len(channelIntegration.Data) > 0 {
			err := json.Unmarshal(channelIntegration.Data, &parsedData)
			if err != nil {
				// Skip if we can't parse the data, but still create entry with session key
				parsedData = lastfmSeparateTableData20251202000000{}
			}
		}

		// Only migrate if we have a session key (apiKey field)
		if channelIntegration.APIKey == nil || *channelIntegration.APIKey == "" {
			continue
		}

		insertQuery := `
INSERT INTO channels_integrations_lastfm (
	channel_id,
	enabled,
	session_key,
	username,
	avatar
)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (channel_id) DO NOTHING
`

		_, err := tx.ExecContext(
			ctx,
			insertQuery,
			channelIntegration.ChannelID,
			channelIntegration.Enabled,
			channelIntegration.APIKey,
			parsedData.UserName,
			parsedData.Avatar,
		)
		if err != nil {
			return fmt.Errorf("insert lastfm integration for channel %s: %w", channelIntegration.ChannelID, err)
		}
	}

	return nil
}

func downLastfmSeparateTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "DROP TABLE IF EXISTS channels_integrations_lastfm")
	return err
}
