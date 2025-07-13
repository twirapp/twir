package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

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
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	data jsonb
);

CREATE UNIQUE INDEX IF NOT EXISTS channels_overlays_kappagen_channel_id_unique ON channels_overlays_kappagen (channel_id);
`

	if _, err := tx.ExecContext(ctx, tablesCreateQuery); err != nil {
		return fmt.Errorf("create tables: %w", err)
	}

	type OldKappagenOverlaySettingsEmotes struct {
		Time           int32 `json:"time,omitempty"`
		Max            int32 `json:"max,omitempty"`
		Queue          int32 `json:"queue,omitempty"`
		FfzEnabled     bool  `json:"ffzEnabled,omitempty"`
		BttvEnabled    bool  `json:"bttvEnabled,omitempty"`
		SevenTvEnabled bool  `json:"sevenTvEnabled,omitempty"`
		EmojiStyle     int   `json:"emojiStyle,omitempty"`
	}

	type OldKappagenOverlaySettingsSize struct {
		// from 7 to 20
		RatioNormal float64 `json:"ratioNormal,omitempty"`
		// from 14 to 40
		RatioSmall float64 `json:"ratioSmall,omitempty"`
		Min        int32   `json:"min,omitempty"`
		Max        int32   `json:"max,omitempty"`
	}

	type OldKappagenOverlaySettingsCube struct {
		Speed int32 `json:"speed,omitempty"`
	}

	type OldKappagenOverlaySettingsAnimation struct {
		FadeIn  bool `json:"fadeIn,omitempty"`
		FadeOut bool `json:"fadeOut,omitempty"`
		ZoomIn  bool `json:"zoomIn,omitempty"`
		ZoomOut bool `json:"zoomOut,omitempty"`
	}

	type OldKappagenOverlaySettingsAnimationSettingsPrefs struct {
		Size    *float64 `json:"size"`
		Center  *bool    `json:"center"`
		Speed   *int64   `json:"speed"`
		Faces   *bool    `json:"faces"`
		Message []string `json:"message"`
		Time    *int64   `json:"time"`
	}

	type OldKappagenOverlaySettingsAnimationSettings struct {
		Style   string                                            `json:"style"`
		Prefs   *OldKappagenOverlaySettingsAnimationSettingsPrefs `json:"prefs"`
		Count   *int64                                            `json:"count"`
		Enabled bool                                              `json:"enabled"`
	}

	type OldKappagenOverlaySettingsEvent struct {
		Event          int32    `json:"event"`
		DisabledStyles []string `json:"disabledStyles,omitempty"`
		Enabled        bool     `json:"enabled,omitempty"`
	}

	type OldKappagenOverlaySettings struct {
		Emotes         OldKappagenOverlaySettingsEmotes              `json:"emotes,omitempty"`
		Size           OldKappagenOverlaySettingsSize                `json:"size,omitempty"`
		Cube           OldKappagenOverlaySettingsCube                `json:"cube,omitempty"`
		Animation      OldKappagenOverlaySettingsAnimation           `json:"animation,omitempty"`
		Animations     []OldKappagenOverlaySettingsAnimationSettings `json:"animations,omitempty"`
		EnableRave     bool                                          `json:"enableRave,omitempty"`
		Events         []OldKappagenOverlaySettingsEvent             `json:"events,omitempty"`
		EnableSpawn    bool                                          `json:"enableSpawn,omitempty"`
		ExcludedEmotes []string                                      `json:"excludedEmotes,omitempty"`
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
		return fmt.Errorf("rows err: %w", rows.Err())
	}

	eventToTwirEvent := map[int32]string{
		0:  "FOLLOW",
		1:  "SUBSCRIBE",
		2:  "RESUBSCRIBE",
		3:  "SUB_GIFT",
		4:  "REDEMPTION_CREATED",
		5:  "COMMAND_USED",
		6:  "FIRST_USER_MESSAGE",
		7:  "RAIDED",
		8:  "TITLE_OR_CATEGORY_CHANGED",
		9:  "STREAM_ONLINE",
		10: "STREAM_OFFLINE",
		11: "ON_CHAT_CLEAR",
		12: "DONATE",
		13: "KEYWORD_MATCHED",
		14: "GREETING_SENDED",
		15: "POLL_BEGIN",
		16: "POLL_PROGRESS",
		17: "POLL_END",
		18: "PREDICTION_BEGIN",
		19: "PREDICTION_PROGRESS",
		20: "PREDICTION_END",
		21: "PREDICTION_LOCK",
		22: "CHANNEL_BAN",
		23: "CHANNEL_UNBAN_REQUEST_CREATE",
		24: "CHANNEL_UNBAN_REQUEST_RESOLVE",
		25: "CHANNEL_MESSAGE_DELETE",
	}

	type KappagenOverlayEmotesSettings struct {
		Time           int  `json:"time,omitempty"`
		Max            int  `json:"max,omitempty"`
		Queue          int  `json:"queue,omitempty"`
		FfzEnabled     bool `json:"ffz_enabled,omitempty"`
		BttvEnabled    bool `json:"bttv_enabled,omitempty"`
		SevenTvEnabled bool `json:"seven_tv_enabled,omitempty"`
		EmojiStyle     int  `json:"emoji_style,omitempty"`
	}

	type KappagenOverlaySizeSettings struct {
		RatioNormal float64 `json:"ratio_normal,omitempty"`
		RatioSmall  float64 `json:"ratio_small,omitempty"`
		Min         int     `json:"min,omitempty"`
		Max         int     `json:"max,omitempty"`
	}

	type KappagenOverlayAnimationSettings struct {
		FadeIn  bool `json:"fade_in,omitempty"`
		FadeOut bool `json:"fade_out,omitempty"`
		ZoomIn  bool `json:"zoom_in,omitempty"`
		ZoomOut bool `json:"zoom_out,omitempty"`
	}

	type KappagenOverlayAnimationsPrefsSettings struct {
		Size    float64  `json:"size,omitempty"`
		Center  bool     `json:"center,omitempty"`
		Speed   int      `json:"speed,omitempty"`
		Faces   bool     `json:"faces,omitempty"`
		Message []string `json:"message,omitempty"`
		Time    int      `json:"time,omitempty"`
	}

	type KappagenOverlayAnimationsSettings struct {
		Style   string                                  `json:"style,omitempty"`
		Prefs   *KappagenOverlayAnimationsPrefsSettings `json:"prefs"`
		Count   *int                                    `json:"count,omitempty"`
		Enabled bool                                    `json:"enabled,omitempty"`
	}

	type KappagenOverlayEvent struct {
		Event              string   `json:"event,omitempty"`
		DisabledAnimations []string `json:"disabled_animations,omitempty"`
		Enabled            bool     `json:"enabled,omitempty"`
	}

	type KappagenOverlaySettings struct {
		EnableSpawn    bool                                `json:"enable_spawn,omitempty"`
		ExcludedEmotes []string                            `json:"excluded_emotes,omitempty"`
		EnableRave     bool                                `json:"enable_rave,omitempty"`
		Animation      KappagenOverlayAnimationSettings    `json:"animation"`
		Animations     []KappagenOverlayAnimationsSettings `json:"animations,omitempty"`
		Emotes         KappagenOverlayEmotesSettings       `json:"emotes"`
		Size           KappagenOverlaySizeSettings         `json:"size"`
		Events         []KappagenOverlayEvent              `json:"events,omitempty"`
	}

	for _, update := range forUpdate {
		var oldSettings OldKappagenOverlaySettings
		if err := json.Unmarshal(update.settingsBytes, &oldSettings); err != nil {
			return err
		}

		newSettings := KappagenOverlaySettings{
			EnableSpawn:    oldSettings.EnableSpawn,
			ExcludedEmotes: oldSettings.ExcludedEmotes,
			EnableRave:     oldSettings.EnableRave,
			Animation: KappagenOverlayAnimationSettings{
				FadeIn:  oldSettings.Animation.FadeIn,
				FadeOut: oldSettings.Animation.FadeOut,
				ZoomIn:  oldSettings.Animation.ZoomIn,
				ZoomOut: oldSettings.Animation.ZoomOut,
			},
			Animations: []KappagenOverlayAnimationsSettings{},
			Emotes: KappagenOverlayEmotesSettings{
				Time:           int(oldSettings.Emotes.Time),
				Max:            int(oldSettings.Emotes.Max),
				Queue:          int(oldSettings.Emotes.Queue),
				FfzEnabled:     oldSettings.Emotes.FfzEnabled,
				BttvEnabled:    oldSettings.Emotes.BttvEnabled,
				SevenTvEnabled: oldSettings.Emotes.SevenTvEnabled,
				EmojiStyle:     oldSettings.Emotes.EmojiStyle,
			},
			Size: KappagenOverlaySizeSettings{
				RatioNormal: oldSettings.Size.RatioNormal,
				RatioSmall:  oldSettings.Size.RatioSmall,
				Min:         int(oldSettings.Size.Min),
				Max:         int(oldSettings.Size.Max),
			},
		}
		for _, animation := range oldSettings.Animations {
			var (
				size      float64
				center    bool
				speed     int64
				faces     bool
				prefsTime int64
				count     *int64
				messages  []string
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
					prefsTime = *animation.Prefs.Time
				}
				if animation.Count != nil {
					count = animation.Count
				}
				messages = animation.Prefs.Message
			}

			s := KappagenOverlayAnimationsSettings{
				Style: animation.Style,
				Prefs: &KappagenOverlayAnimationsPrefsSettings{
					Size:    size,
					Center:  center,
					Speed:   int(speed),
					Faces:   faces,
					Message: messages,
					Time:    int(prefsTime),
				},
				Count:   nil,
				Enabled: animation.Enabled,
			}

			if count != nil {
				countInt := int(*count)
				s.Count = &countInt
			}

			newSettings.Animations = append(
				newSettings.Animations, s,
			)
		}
		for _, event := range oldSettings.Events {
			newSettings.Events = append(
				newSettings.Events,
				KappagenOverlayEvent{
					Event:              eventToTwirEvent[event.Event],
					DisabledAnimations: event.DisabledStyles,
					Enabled:            event.Enabled,
				},
			)
		}

		newSettingsBytes, err := json.Marshal(newSettings)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO channels_overlays_kappagen (id, channel_id, data) VALUES ($1, $2, $3)`,
			update.id,
			update.channelId,
			newSettingsBytes,
		)
		if err != nil {
			return fmt.Errorf("insert: %w", err)
		}
	}

	deleteQuery := `
DELETE FROM channels_modules_settings WHERE type = 'kappagen_overlay';
`
	_, err = tx.ExecContext(ctx, deleteQuery)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	// This code is executed when the migration is applied.
	return nil
}

func downKappagenSeparateTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
