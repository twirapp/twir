package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upKappagenEventsAnimations, downKappagenEventsAnimations)
}

func upKappagenEventsAnimations(ctx context.Context, tx *sql.Tx) error {
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
		var oldSettings KappagenOverlaySettings20231118120418old
		if err := json.Unmarshal(update.settingsBytes, &oldSettings); err != nil {
			return err
		}

		enabledEvents := make(
			[]KappagenOverlaySettingsEnabledEvents20231118120418new,
			len(oldSettings.EnabledEvents),
		)
		for i, event := range oldSettings.EnabledEvents {
			enabledEvents[i] = KappagenOverlaySettingsEnabledEvents20231118120418new{
				Event:   event,
				Enabled: true,
			}
		}

		newSettings := KappagenOverlaySettings20231118120418new{
			Emotes:        oldSettings.Emotes,
			Size:          oldSettings.Size,
			Cube:          oldSettings.Cube,
			Animation:     oldSettings.Animation,
			Animations:    oldSettings.Animations,
			EnableRave:    oldSettings.EnableRave,
			EnableSpawn:   oldSettings.EnableSpawn,
			EnabledEvents: enabledEvents,
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

func downKappagenEventsAnimations(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}

type KappagenOverlaySettingsEmotes20231118120418 struct {
	Time  int32 `json:"time,omitempty"`
	Max   int32 `json:"max,omitempty"`
	Queue int32 `json:"queue,omitempty"`
}

type KappagenOverlaySettingsSize20231118120418 struct {
	// from 7 to 20
	RatioNormal int32 `json:"ratioNormal,omitempty"`
	// from 14 to 40
	RatioSmall int32 `json:"ratioSmall,omitempty"`
	Min        int32 `json:"min,omitempty"`
	Max        int32 `json:"max,omitempty"`
}

type KappagenOverlaySettingsCube20231118120418 struct {
	Speed int32 `json:"speed,omitempty"`
}

type KappagenOverlaySettingsAnimation20231118120418 struct {
	FadeIn  bool `json:"fadeIn,omitempty"`
	FadeOut bool `json:"fadeOut,omitempty"`
	ZoomIn  bool `json:"zoomIn,omitempty"`
	ZoomOut bool `json:"zoomOut,omitempty"`
}

type KappagenOverlaySettingsAnimationSettingsPrefs20231118120418old struct {
	Size    *float32 `json:"size"`
	Center  *bool    `json:"center"`
	Speed   *int32   `json:"speed"`
	Faces   *bool    `json:"faces"`
	Time    *int32   `json:"time"`
	Message []string `json:"message"`
}

type KappagenOverlaySettingsAnimationSettings20231118120418 struct {
	Prefs   *KappagenOverlaySettingsAnimationSettingsPrefs20231118120418old `json:"prefs"`
	Count   *int32                                                          `json:"count"`
	Style   string                                                          `json:"style"`
	Enabled bool                                                            `json:"enabled"`
}

type KappagenOverlaySettings20231118120418old struct {
	Animations    []KappagenOverlaySettingsAnimationSettings20231118120418 `json:"animations,omitempty"`
	EnabledEvents []int32                                                  `json:"events,omitempty"`
	Size          KappagenOverlaySettingsSize20231118120418                `json:"size,omitempty"`
	Emotes        KappagenOverlaySettingsEmotes20231118120418              `json:"emotes,omitempty"`
	Cube          KappagenOverlaySettingsCube20231118120418                `json:"cube,omitempty"`
	Animation     KappagenOverlaySettingsAnimation20231118120418           `json:"animation,omitempty"`
	EnableRave    bool                                                     `json:"enableRave,omitempty"`
	EnableSpawn   bool                                                     `json:"enableSpawn,omitempty"`
}

//

type KappagenOverlaySettingsEnabledEvents20231118120418new struct {
	DisabledStyles []string `json:"disabledStyles,omitempty"`
	Event          int32    `json:"event"`
	Enabled        bool     `json:"enabled"`
}

type KappagenOverlaySettings20231118120418new struct {
	Animations    []KappagenOverlaySettingsAnimationSettings20231118120418 `json:"animations,omitempty"`
	EnabledEvents []KappagenOverlaySettingsEnabledEvents20231118120418new  `json:"events,omitempty"`
	Size          KappagenOverlaySettingsSize20231118120418                `json:"size,omitempty"`
	Emotes        KappagenOverlaySettingsEmotes20231118120418              `json:"emotes,omitempty"`
	Cube          KappagenOverlaySettingsCube20231118120418                `json:"cube,omitempty"`
	Animation     KappagenOverlaySettingsAnimation20231118120418           `json:"animation,omitempty"`
	EnableRave    bool                                                     `json:"enableRave,omitempty"`
	EnableSpawn   bool                                                     `json:"enableSpawn,omitempty"`
}
