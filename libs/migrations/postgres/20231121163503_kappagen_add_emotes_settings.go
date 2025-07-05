package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upKappagenAddEmotesSettings, downKappagenAddEmotesSettings)
}

func upKappagenAddEmotesSettings(ctx context.Context, tx *sql.Tx) error {
	findQuery := `
SELECT id, settings FROM channels_modules_settings WHERE type = 'kappagen_overlay'
`

	rows, err := tx.QueryContext(ctx, findQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	var forUpdate []struct {
		id            string
		settingsBytes []byte
	}

	for rows.Next() {
		var id string
		var settingsBytes []byte
		if err := rows.Scan(&id, &settingsBytes); err != nil {
			return err
		}

		forUpdate = append(
			forUpdate, struct {
				id            string
				settingsBytes []byte
			}{
				id:            id,
				settingsBytes: settingsBytes,
			},
		)
	}

	for _, update := range forUpdate {
		var oldSettings KappagenOverlaySettings20231121163503old
		if err := json.Unmarshal(update.settingsBytes, &oldSettings); err != nil {
			return err
		}

		newEmotes := KappagenOverlaySettingsEmotes20231121163503new{
			Time:           oldSettings.Emotes.Time,
			Max:            oldSettings.Emotes.Max,
			Queue:          oldSettings.Emotes.Queue,
			FfzEnabled:     true,
			BttvEnabled:    true,
			SevenTvEnabled: true,
			EmojiStyle:     0,
		}

		newSettings := KappagenOverlaySettings20231121163503new{
			Emotes:      newEmotes,
			Size:        oldSettings.Size,
			Cube:        oldSettings.Cube,
			Animation:   oldSettings.Animation,
			Animations:  oldSettings.Animations,
			EnableRave:  oldSettings.EnableRave,
			Events:      oldSettings.Events,
			EnableSpawn: oldSettings.EnableSpawn,
		}

		newSettingsBytes, err := json.Marshal(newSettings)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(
			ctx,
			`UPDATE channels_modules_settings SET settings = $1 WHERE id = $2`,
			newSettingsBytes,
			update.id,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func downKappagenAddEmotesSettings(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}

type KappagenOverlaySettingsEmotes20231121163503old struct {
	Time  int32 `json:"time,omitempty"`
	Max   int32 `json:"max,omitempty"`
	Queue int32 `json:"queue,omitempty"`
}

type KappagenOverlaySettingsEmotes20231121163503new struct {
	Time           int32 `json:"time,omitempty"`
	Max            int32 `json:"max,omitempty"`
	Queue          int32 `json:"queue,omitempty"`
	FfzEnabled     bool  `json:"ffzEnabled,omitempty"`
	BttvEnabled    bool  `json:"bttvEnabled,omitempty"`
	SevenTvEnabled bool  `json:"sevenTvEnabled,omitempty"`
	EmojiStyle     int32 `json:"emojiStyle,omitempty"`
}

type KappagenOverlaySettingsSize20231121163503old struct {
	// from 7 to 20
	RatioNormal float64 `json:"ratioNormal,omitempty"`
	// from 14 to 40
	RatioSmall float64 `json:"ratioSmall,omitempty"`
	Min        int32   `json:"min,omitempty"`
	Max        int32   `json:"max,omitempty"`
}

type KappagenOverlaySettingsCube20231121163503old struct {
	Speed int32 `json:"speed,omitempty"`
}

type KappagenOverlaySettingsAnimation20231121163503old struct {
	FadeIn  bool `json:"fadeIn,omitempty"`
	FadeOut bool `json:"fadeOut,omitempty"`
	ZoomIn  bool `json:"zoomIn,omitempty"`
	ZoomOut bool `json:"zoomOut,omitempty"`
}

type KappagenOverlaySettingsAnimationSettingsPrefs20231121163503old struct {
	Size    *float64 `json:"size"`
	Center  *bool    `json:"center"`
	Speed   *int32   `json:"speed"`
	Faces   *bool    `json:"faces"`
	Time    *int32   `json:"time"`
	Message []string `json:"message"`
}

type KappagenOverlaySettingsAnimationSettings20231121163503old struct {
	Prefs   *KappagenOverlaySettingsAnimationSettingsPrefs20231121163503old `json:"prefs"`
	Count   *int32                                                          `json:"count"`
	Style   string                                                          `json:"style"`
	Enabled bool                                                            `json:"enabled"`
}

type KappagenOverlaySettingsEvent20231121163503old struct {
	DisabledStyles []string `json:"disabledStyles,omitempty"`
	Event          int32    `json:"event"`
	Enabled        bool     `json:"enabled,omitempty"`
}

type KappagenOverlaySettings20231121163503old struct {
	Animations  []KappagenOverlaySettingsAnimationSettings20231121163503old `json:"animations,omitempty"`
	Events      []KappagenOverlaySettingsEvent20231121163503old             `json:"events,omitempty"`
	Size        KappagenOverlaySettingsSize20231121163503old                `json:"size,omitempty"`
	Emotes      KappagenOverlaySettingsEmotes20231121163503old              `json:"emotes,omitempty"`
	Cube        KappagenOverlaySettingsCube20231121163503old                `json:"cube,omitempty"`
	Animation   KappagenOverlaySettingsAnimation20231121163503old           `json:"animation,omitempty"`
	EnableRave  bool                                                        `json:"enableRave,omitempty"`
	EnableSpawn bool                                                        `json:"enableSpawn,omitempty"`
}

type KappagenOverlaySettings20231121163503new struct {
	Animations  []KappagenOverlaySettingsAnimationSettings20231121163503old `json:"animations,omitempty"`
	Events      []KappagenOverlaySettingsEvent20231121163503old             `json:"events,omitempty"`
	Size        KappagenOverlaySettingsSize20231121163503old                `json:"size,omitempty"`
	Emotes      KappagenOverlaySettingsEmotes20231121163503new              `json:"emotes,omitempty"`
	Cube        KappagenOverlaySettingsCube20231121163503old                `json:"cube,omitempty"`
	Animation   KappagenOverlaySettingsAnimation20231121163503old           `json:"animation,omitempty"`
	EnableRave  bool                                                        `json:"enableRave,omitempty"`
	EnableSpawn bool                                                        `json:"enableSpawn,omitempty"`
}
