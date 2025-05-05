package migrations

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upKappagenSeparateTable, downKappagenSeparateTable)
}

func upKappagenSeparateTable(ctx context.Context, tx *sql.Tx) error {
	tablesCreateQuery := `
CREATE TABLE channels_overlays_kappagen (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	channel_id TEXT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
	enable_spawn BOOLEAN NOT NULL,
	excluded_emotes TEXT[] NOT NULL,
	enable_rave BOOLEAN NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE channels_overlays_kappagen_animations (
	id ulid PRIMARY KEY DEFAULT gen_ulid(),
	overlay_id UUID NOT NULL REFERENCES channels_overlays_kappagen(id) ON DELETE CASCADE,
	style TEXT NOT NULL,
	count INT NOT NULL,
	enabled BOOLEAN NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE channels_overlays_kappagen_animations_prefs (
	id ulid PRIMARY KEY DEFAULT gen_ulid(),
	animation_id ulid NOT NULL REFERENCES channels_overlays_kappagen_animations(id) ON DELETE CASCADE,
	size float4 NOT NULL,
	center BOOLEAN NOT NULL,
	speed INT NOT NULL,
	faces BOOLEAN NOT NULL,
	message TEXT[] NOT NULL,
	time INT NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE UNIQUE INDEX channels_overlays_kappagen.channel_id_unique ON channels_overlays_kappagen (channel_id);
`

	if _, err := tx.ExecContext(ctx, tablesCreateQuery); err != nil {
		return err
	}

	type KappagenOverlaySettingsEmotes struct {
		Time           int32 `json:"time,omitempty"`
		Max            int32 `json:"max,omitempty"`
		Queue          int32 `json:"queue,omitempty"`
		FfzEnabled     bool  `json:"ffzEnabled,omitempty"`
		BttvEnabled    bool  `json:"bttvEnabled,omitempty"`
		SevenTvEnabled bool  `json:"sevenTvEnabled,omitempty"`
		EmojiStyle     int   `json:"emojiStyle,omitempty"`
	}

	type KappagenOverlaySettingsSize struct {
		// from 7 to 20
		RatioNormal float64 `json:"ratioNormal,omitempty"`
		// from 14 to 40
		RatioSmall float64 `json:"ratioSmall,omitempty"`
		Min        int32   `json:"min,omitempty"`
		Max        int32   `json:"max,omitempty"`
	}

	type KappagenOverlaySettingsCube struct {
		Speed int32 `json:"speed,omitempty"`
	}

	type KappagenOverlaySettingsAnimation struct {
		FadeIn  bool `json:"fadeIn,omitempty"`
		FadeOut bool `json:"fadeOut,omitempty"`
		ZoomIn  bool `json:"zoomIn,omitempty"`
		ZoomOut bool `json:"zoomOut,omitempty"`
	}

	type KappagenOverlaySettingsAnimationSettingsPrefs struct {
		Size    *float64 `json:"size"`
		Center  *bool    `json:"center"`
		Speed   *int64   `json:"speed"`
		Faces   *bool    `json:"faces"`
		Message []string `json:"message"`
		Time    *int64   `json:"time"`
	}

	type KappagenOverlaySettingsAnimationSettings struct {
		Style   string                                         `json:"style"`
		Prefs   *KappagenOverlaySettingsAnimationSettingsPrefs `json:"prefs"`
		Count   *int64                                         `json:"count"`
		Enabled bool                                           `json:"enabled"`
	}

	type KappagenOverlaySettingsEvent struct {
		Event          int32    `json:"event"`
		DisabledStyles []string `json:"disabledStyles,omitempty"`
		Enabled        bool     `json:"enabled,omitempty"`
	}

	type KappagenOverlaySettings struct {
		Emotes         KappagenOverlaySettingsEmotes              `json:"emotes,omitempty"`
		Size           KappagenOverlaySettingsSize                `json:"size,omitempty"`
		Cube           KappagenOverlaySettingsCube                `json:"cube,omitempty"`
		Animation      KappagenOverlaySettingsAnimation           `json:"animation,omitempty"`
		Animations     []KappagenOverlaySettingsAnimationSettings `json:"animations,omitempty"`
		EnableRave     bool                                       `json:"enableRave,omitempty"`
		Events         []KappagenOverlaySettingsEvent             `json:"events,omitempty"`
		EnableSpawn    bool                                       `json:"enableSpawn,omitempty"`
		ExcludedEmotes []string                                   `json:"excludedEmotes,omitempty"`
	}

	findQuery := `
SELECT id, settings, "channelId" FROM channels_modules_settings WHERE type = 'kappagen_overlay'
`

	rows, err := tx.QueryContext(ctx, findQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	var forUpdate []struct {
		id            string
		channelId     string
		settingsBytes []byte
	}

	for rows.Next() {
		var id string
		var channelId string
		var settingsBytes []byte
		if err := rows.Scan(&id, &settingsBytes, &channelId); err != nil {
			return err
		}

		forUpdate = append(
			forUpdate, struct {
				id            string
				channelId     string
				settingsBytes []byte
			}{
				id:            id,
				settingsBytes: settingsBytes,
				channelId:     channelId,
			},
		)
	}

	if rows.Err() != nil {
		return rows.Err()
	}

	for _, update := range forUpdate {
		var oldSettings KappagenOverlaySettings
		if err := json.Unmarshal(update.settingsBytes, &oldSettings); err != nil {
			return err
		}

		var overlayId string
		overlayRow := tx.QueryRowContext(
			ctx,
			`INSERT INTO channels_overlays_kappagen (
			 	id,
				channel_id,
				enable_spawn,
				excluded_emotes,
				enable_rave
			) VALUES (
				$1,
				$2,
				$3,
				$4,
				$5
			) RETURNING id`,
			update.id,
			update.channelId,
			oldSettings.EnableSpawn,
			append(pq.StringArray{}, oldSettings.ExcludedEmotes...),
			oldSettings.EnableRave,
		)
		if err := overlayRow.Scan(&overlayId); err != nil {
			return err
		}

		for _, animation := range oldSettings.Animations {
			var animationId string
			var count int64
			if animation.Count != nil {
				count = *animation.Count
			}

			animationRow := tx.QueryRowContext(
				ctx,
				`INSERT INTO channels_overlays_kappagen_animations (
						overlay_id,
						style,
						count,
						enabled
				) VALUES (
						$1,
						$2,
						$3,
						$4
				) RETURNING id`,
				overlayId,
				animation.Style,
				count,
				animation.Enabled,
			)
			if err := animationRow.Scan(&animationId); err != nil {
				return err
			}

			var (
				size    float64
				center  bool
				speed   int64
				faces   bool
				message []string
				time    int64
			)

			if animation.Prefs != nil {
				if animation.Prefs.Size != nil {
					size = *animation.Prefs.Size
				}

				if animation.Prefs.Center != nil {
					center = *animation.Prefs.Center
				}

				if animation.Prefs.Speed != nil {
					speed = *animation.Prefs.Speed
				}

				if animation.Prefs.Faces != nil {
					faces = *animation.Prefs.Faces
				}

				if animation.Prefs.Time != nil {
					time = *animation.Prefs.Time
				}

				message = animation.Prefs.Message
			}

			_, err = tx.ExecContext(
				ctx,
				`INSERT INTO channels_overlays_kappagen_animations_prefs (
						animation_id,
						size,
						center,
						speed,
						faces,
						message,
						time
				) VALUES (
						$1,
						$2,
						$3,
						$4,
						$5,
						$6,
						$7
				)`,
				animationId,
				size,
				center,
				speed,
				faces,
				append(pq.StringArray{}, message...),
				time,
			)
			if err != nil {
				return err
			}
		}
	}

	deleteQuery := `
DELETE FROM channels_modules_settings WHERE type = 'kappagen_overlay'
`
	_, err = tx.ExecContext(ctx, deleteQuery)
	if err != nil {
		return err
	}

	// This code is executed when the migration is applied.
	return nil
}

func downKappagenSeparateTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
