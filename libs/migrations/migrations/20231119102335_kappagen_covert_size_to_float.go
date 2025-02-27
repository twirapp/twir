package migrations

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upKappagenCovertSizeToFloat, downKappagenCovertSizeToFloat)
}

func upKappagenCovertSizeToFloat(ctx context.Context, tx *sql.Tx) error {
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
		var oldSettings KappagenOverlaySettings20231119102335old
		if err := json.Unmarshal(update.settingsBytes, &oldSettings); err != nil {
			return err
		}

		formattedNormal := fmt.Sprintf("%.2f", 1/float64(oldSettings.Size.RatioNormal))
		formattedSmall := fmt.Sprintf("%.2f", 1/float64(oldSettings.Size.RatioSmall))

		normal, _ := strconv.ParseFloat(formattedNormal, 64)
		small, _ := strconv.ParseFloat(formattedSmall, 64)

		newSize := KappagenOverlaySettingsSize20231119102335new{
			RatioNormal: normal,
			RatioSmall:  small,
			Min:         oldSettings.Size.Min,
			Max:         oldSettings.Size.Max,
		}

		newSettings := KappagenOverlaySettings20231119102335new{
			Emotes:      oldSettings.Emotes,
			Size:        newSize,
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

func downKappagenCovertSizeToFloat(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}

type KappagenOverlaySettingsEmotes20231119102335old struct {
	Time  int32 `json:"time,omitempty"`
	Max   int32 `json:"max,omitempty"`
	Queue int32 `json:"queue,omitempty"`
}

type KappagenOverlaySettingsSize20231119102335old struct {
	// from 7 to 20
	RatioNormal int32 `json:"ratioNormal,omitempty"`
	// from 14 to 40
	RatioSmall int32 `json:"ratioSmall,omitempty"`
	Min        int32 `json:"min,omitempty"`
	Max        int32 `json:"max,omitempty"`
}

type KappagenOverlaySettingsSize20231119102335new struct {
	// from 7 to 20
	RatioNormal float64 `json:"ratioNormal,omitempty"`
	// from 14 to 40
	RatioSmall float64 `json:"ratioSmall,omitempty"`
	Min        int32   `json:"min,omitempty"`
	Max        int32   `json:"max,omitempty"`
}

type KappagenOverlaySettingsCube20231119102335old struct {
	Speed int32 `json:"speed,omitempty"`
}

type KappagenOverlaySettingsAnimation20231119102335old struct {
	FadeIn  bool `json:"fadeIn,omitempty"`
	FadeOut bool `json:"fadeOut,omitempty"`
	ZoomIn  bool `json:"zoomIn,omitempty"`
	ZoomOut bool `json:"zoomOut,omitempty"`
}

type KappagenOverlaySettingsAnimationSettingsPrefs20231119102335old struct {
	Size    *float64 `json:"size"`
	Center  *bool    `json:"center"`
	Speed   *int32   `json:"speed"`
	Faces   *bool    `json:"faces"`
	Time    *int32   `json:"time"`
	Message []string `json:"message"`
}

type KappagenOverlaySettingsAnimationSettings20231119102335old struct {
	Prefs   *KappagenOverlaySettingsAnimationSettingsPrefs20231119102335old `json:"prefs"`
	Count   *int32                                                          `json:"count"`
	Style   string                                                          `json:"style"`
	Enabled bool                                                            `json:"enabled"`
}

type KappagenOverlaySettingsEvent20231119102335old struct {
	DisabledStyles []string `json:"disabledStyles,omitempty"`
	Event          int32    `json:"event"`
	Enabled        bool     `json:"enabled,omitempty"`
}

type KappagenOverlaySettings20231119102335old struct {
	Animations  []KappagenOverlaySettingsAnimationSettings20231119102335old `json:"animations,omitempty"`
	Events      []KappagenOverlaySettingsEvent20231119102335old             `json:"events,omitempty"`
	Size        KappagenOverlaySettingsSize20231119102335old                `json:"size,omitempty"`
	Emotes      KappagenOverlaySettingsEmotes20231119102335old              `json:"emotes,omitempty"`
	Cube        KappagenOverlaySettingsCube20231119102335old                `json:"cube,omitempty"`
	Animation   KappagenOverlaySettingsAnimation20231119102335old           `json:"animation,omitempty"`
	EnableRave  bool                                                        `json:"enableRave,omitempty"`
	EnableSpawn bool                                                        `json:"enableSpawn,omitempty"`
}

type KappagenOverlaySettings20231119102335new struct {
	Animations  []KappagenOverlaySettingsAnimationSettings20231119102335old `json:"animations,omitempty"`
	Events      []KappagenOverlaySettingsEvent20231119102335old             `json:"events,omitempty"`
	Size        KappagenOverlaySettingsSize20231119102335new                `json:"size,omitempty"`
	Emotes      KappagenOverlaySettingsEmotes20231119102335old              `json:"emotes,omitempty"`
	Cube        KappagenOverlaySettingsCube20231119102335old                `json:"cube,omitempty"`
	Animation   KappagenOverlaySettingsAnimation20231119102335old           `json:"animation,omitempty"`
	EnableRave  bool                                                        `json:"enableRave,omitempty"`
	EnableSpawn bool                                                        `json:"enableSpawn,omitempty"`
}
