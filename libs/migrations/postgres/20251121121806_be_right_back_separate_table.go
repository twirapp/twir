package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upBeRightBackSeparateTable, downBeRightBackSeparateTable)
}

func upBeRightBackSeparateTable(ctx context.Context, tx *sql.Tx) error {
	tablesCreateQuery := `
CREATE TABLE channels_overlays_be_right_back (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	channel_id TEXT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	text TEXT NOT NULL DEFAULT 'BRB',
	late_enabled BOOLEAN NOT NULL DEFAULT false,
	late_text TEXT NOT NULL DEFAULT 'Streamer is late',
	late_display_brb_time BOOLEAN NOT NULL DEFAULT false,
	background_color TEXT NOT NULL DEFAULT '#000000',
	font_size INTEGER NOT NULL DEFAULT 48,
	font_color TEXT NOT NULL DEFAULT '#FFFFFF',
	font_family TEXT NOT NULL DEFAULT 'Roboto'
);

CREATE UNIQUE INDEX IF NOT EXISTS channels_overlays_be_right_back_channel_id_unique ON channels_overlays_be_right_back (channel_id);
`

	if _, err := tx.ExecContext(ctx, tablesCreateQuery); err != nil {
		return fmt.Errorf("create tables: %w", err)
	}

	type OldBeRightBackOverlaySettingsLate struct {
		Enabled        bool   `json:"enabled"`
		Text           string `json:"text"`
		DisplayBrbTime bool   `json:"displayBrbTime"`
	}

	type OldBeRightBackOverlaySettings struct {
		Text            string                            `json:"text"`
		Late            OldBeRightBackOverlaySettingsLate `json:"late"`
		BackgroundColor string                            `json:"backgroundColor"`
		FontSize        int32                             `json:"fontSize"`
		FontColor       string                            `json:"fontColor"`
		FontFamily      string                            `json:"fontFamily"`
	}

	findQuery := `
SELECT id, settings, "channelId" FROM channels_modules_settings WHERE type = 'be_right_back_overlay'
`

	rows, err := tx.QueryContext(ctx, findQuery)
	if err != nil {
		return fmt.Errorf("find query: %w", err)
	}
	defer rows.Close()

	type migrationData struct {
		id        string
		channelID string
		settings  OldBeRightBackOverlaySettings
	}

	data := make([]migrationData, 0)
	for rows.Next() {
		var id string
		var channelID string
		var settingsRaw []byte

		if err := rows.Scan(&id, &settingsRaw, &channelID); err != nil {
			return fmt.Errorf("scan: %w", err)
		}

		var settings OldBeRightBackOverlaySettings
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
INSERT INTO channels_overlays_be_right_back (
	channel_id,
	text,
	late_enabled,
	late_text,
	late_display_brb_time,
	background_color,
	font_size,
	font_color,
	font_family
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
		`

		if _, err := tx.ExecContext(
			ctx,
			insertQuery,
			d.channelID,
			d.settings.Text,
			d.settings.Late.Enabled,
			d.settings.Late.Text,
			d.settings.Late.DisplayBrbTime,
			d.settings.BackgroundColor,
			d.settings.FontSize,
			d.settings.FontColor,
			d.settings.FontFamily,
		); err != nil {
			return fmt.Errorf("insert: %w", err)
		}
	}

	deleteQuery := `
DELETE FROM channels_modules_settings WHERE type = 'be_right_back_overlay'
	`

	if _, err := tx.ExecContext(ctx, deleteQuery); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func downBeRightBackSeparateTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
