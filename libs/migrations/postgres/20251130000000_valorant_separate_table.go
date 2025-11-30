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
	goose.AddMigrationContext(upValorantSeparateTable, downValorantSeparateTable)
}

func upValorantSeparateTable(ctx context.Context, tx *sql.Tx) error {
	type valorantSeparateTableIntegrations20251130000000 struct {
		ID string
	}

	type valorantSeparateTableChannelsIntegrations20251130000000 struct {
		ID           string
		AccessToken  *string
		RefreshToken *string
		ChannelID    string
		Enabled      bool
		Data         []byte
	}

	type valorantSeparateTableData20251130000000 struct {
		UserName             *string `json:"username,omitempty"`
		ValorantActiveRegion *string `json:"valorantActiveRegion,omitempty"`
		ValorantPuuid        *string `json:"valorantPuuid,omitempty"`
	}

	// Get valorant integration ID
	valorantIntegration := valorantSeparateTableIntegrations20251130000000{}
	err := tx.QueryRowContext(
		ctx,
		"SELECT id FROM integrations WHERE service = 'VALORANT'",
	).Scan(&valorantIntegration.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("get valorant integration: %w", err)
	}

	// If no integration found, just create the table and exit
	if errors.Is(err, sql.ErrNoRows) {
		createTableQuery := `
CREATE TABLE IF NOT EXISTS channels_integrations_valorant (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	channel_id TEXT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
	enabled BOOLEAN NOT NULL DEFAULT false,
	access_token TEXT,
	refresh_token TEXT,
	username TEXT,
	valorant_active_region TEXT,
	valorant_puuid TEXT,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS channels_integrations_valorant_channel_id_idx ON channels_integrations_valorant(channel_id);
`
		if _, err := tx.ExecContext(ctx, createTableQuery); err != nil {
			return fmt.Errorf("create table: %w", err)
		}
		return nil
	}

	// Get all valorant integrations from old table
	rows, err := tx.QueryContext(
		ctx,
		`SELECT id, "accessToken", "refreshToken", enabled, "channelId", data FROM channels_integrations WHERE "integrationId" = $1`,
		valorantIntegration.ID,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("get channels integrations: %w", err)
	}

	channelsIntegrations := []valorantSeparateTableChannelsIntegrations20251130000000{}
	if rows != nil {
		defer rows.Close()

		for rows.Next() {
			channelIntegration := valorantSeparateTableChannelsIntegrations20251130000000{}
			err := rows.Scan(
				&channelIntegration.ID,
				&channelIntegration.AccessToken,
				&channelIntegration.RefreshToken,
				&channelIntegration.Enabled,
				&channelIntegration.ChannelID,
				&channelIntegration.Data,
			)
			if err != nil {
				return fmt.Errorf("scan row: %w", err)
			}

			if channelIntegration.Data == nil || channelIntegration.AccessToken == nil || channelIntegration.RefreshToken == nil {
				continue
			}

			channelsIntegrations = append(channelsIntegrations, channelIntegration)
		}

		if rows.Err() != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("rows error: %w", rows.Err())
		}
	}

	// Create new table
	createTableQuery := `
CREATE TABLE IF NOT EXISTS channels_integrations_valorant (
	id SERIAL PRIMARY KEY,
	channel_id TEXT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
	enabled BOOLEAN NOT NULL DEFAULT false,
	access_token TEXT,
	refresh_token TEXT,
	username TEXT,
	valorant_active_region TEXT,
	valorant_puuid TEXT,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS channels_integrations_valorant_channel_id_idx ON channels_integrations_valorant(channel_id);
`
	if _, err := tx.ExecContext(ctx, createTableQuery); err != nil {
		return fmt.Errorf("create table: %w", err)
	}

	// Migrate data
	for _, channelIntegration := range channelsIntegrations {
		var parsedData valorantSeparateTableData20251130000000
		if len(channelIntegration.Data) > 0 {
			err := json.Unmarshal(channelIntegration.Data, &parsedData)
			if err != nil {
				return fmt.Errorf("unmarshal data: %w", err)
			}
		}

		insertQuery := `
INSERT INTO channels_integrations_valorant (
	channel_id,
	enabled,
	access_token,
	refresh_token,
	username,
	valorant_active_region,
	valorant_puuid
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (channel_id) DO NOTHING
`

		_, err := tx.ExecContext(
			ctx,
			insertQuery,
			channelIntegration.ChannelID,
			channelIntegration.Enabled,
			channelIntegration.AccessToken,
			channelIntegration.RefreshToken,
			parsedData.UserName,
			parsedData.ValorantActiveRegion,
			parsedData.ValorantPuuid,
		)
		if err != nil {
			return fmt.Errorf("insert valorant integration: %w", err)
		}
	}

	return nil
}

func downValorantSeparateTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
