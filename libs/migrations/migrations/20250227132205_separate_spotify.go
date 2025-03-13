package migrations

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upSeparateSpotify, downSeparateSpotify)
}

type integrations20250227132205 struct {
	ID string
}

type channelsIntegrationsSpotify20250227132205 struct {
	ID            string
	AccessToken   string
	RefreshToken  string
	IntegrationID string
	ChannelID     string
	Data          []byte
	Enabled       bool
}

type channelsIntegrationsSpotifyData20250227132205 struct {
	Avatar   string `json:"avatar"`
	Username string `json:"username"`
}

func upSeparateSpotify(ctx context.Context, tx *sql.Tx) error {
	spotifyIntegration := integrations20250227132205{}
	err := tx.QueryRowContext(
		ctx,
		"SELECT id FROM integrations WHERE service = 'SPOTIFY'",
	).Scan(&spotifyIntegration.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	rows, err := tx.QueryContext(
		ctx,
		`SELECT id, "accessToken", "refreshToken", enabled, "channelId", data FROM channels_integrations WHERE "integrationId" = $1`,
		spotifyIntegration.ID,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	channelsIntegrations := []channelsIntegrationsSpotify20250227132205{}
	defer rows.Close()

	for rows.Next() {
		channelIntegration := channelsIntegrationsSpotify20250227132205{}
		err := rows.Scan(
			&channelIntegration.ID,
			&channelIntegration.AccessToken,
			&channelIntegration.RefreshToken,
			&channelIntegration.Enabled,
			&channelIntegration.ChannelID,
			&channelIntegration.Data,
		)
		if err != nil {
			return err
		}
		channelsIntegrations = append(channelsIntegrations, channelIntegration)
	}
	if rows.Err() != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	newTableQuery := `
CREATE TABLE IF NOT EXISTS channels_integrations_spotify (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	access_token TEXT NOT NULL,
	refresh_token TEXT NOT NULL,
	enabled BOOLEAN NOT NULL DEFAULT true,
	scopes TEXT[] NOT NULL DEFAULT '{}',
	channel_id text NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
	avatar_uri TEXT,
	username TEXT,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX IF NOT EXISTS channels_integrations_spotify_channel_id_idx ON channels_integrations_spotify(channel_id);
`
	_, err = tx.ExecContext(ctx, newTableQuery)
	if err != nil {
		return err
	}

	for _, channelIntegration := range channelsIntegrations {
		var parsedData channelsIntegrationsSpotifyData20250227132205
		err := json.Unmarshal(channelIntegration.Data, &parsedData)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO channels_integrations_spotify (id, access_token, refresh_token, enabled, channel_id, avatar_uri, username, scopes) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			channelIntegration.ID,
			channelIntegration.AccessToken,
			channelIntegration.RefreshToken,
			channelIntegration.Enabled,
			channelIntegration.ChannelID,
			parsedData.Avatar,
			parsedData.Username,
			pq.StringArray{"user-read-currently-playing"},
		)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(
			ctx,
			`DELETE FROM channels_integrations WHERE id = $1`,
			channelIntegration.ID,
		)
	}

	// This code is executed when the migration is applied.
	return nil
}

func downSeparateSpotify(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
