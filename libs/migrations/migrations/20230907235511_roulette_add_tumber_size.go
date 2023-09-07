package migrations

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upRouletteAddTumberSize, downRouletteAddTumberSize)
}

type rouletteAddTumberSize20230907235511Old struct {
	Enabled               bool `json:"enabled"`
	CanBeUsedByModerators bool `json:"canBeUsedByModerator"`
	TimeoutSeconds        int  `json:"timeoutTime"`
	DecisionSeconds       int  `json:"decisionTime"`
	ChargedBullets        int  `json:"chargedBullets"`

	InitMessage    string `json:"initMessage"`
	SurviveMessage string `json:"surviveMessage"`
	DeathMessage   string `json:"deathMessage"`
}

type rouletteAddTumberSize20230907235511New struct {
	Enabled               bool `json:"enabled"`
	CanBeUsedByModerators bool `json:"canBeUsedByModerator"`
	TimeoutSeconds        int  `json:"timeoutTime"`
	DecisionSeconds       int  `json:"decisionTime"`
	TumberSize            int  `json:"tumberSize"`
	ChargedBullets        int  `json:"chargedBullets"`

	InitMessage    string `json:"initMessage"`
	SurviveMessage string `json:"surviveMessage"`
	DeathMessage   string `json:"deathMessage"`
}

func upRouletteAddTumberSize(ctx context.Context, tx *sql.Tx) error {
	findQuery := `
SELECT id, settings FROM channels_modules_settings WHERE type = 'russian_roulette'
`
	rows, err := tx.QueryContext(ctx, findQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	// iterate over all modules settings
	for rows.Next() {
		var id string
		var settingsBytes []byte
		if err := rows.Scan(&id, &settingsBytes); err != nil {
			return err
		}

		// unmarshal settings
		var settings rouletteAddTumberSize20230907235511Old
		if err := json.Unmarshal(settingsBytes, &settings); err != nil {
			return err
		}

		// update settings
		newSettingsBytes, err := json.Marshal(
			&rouletteAddTumberSize20230907235511New{
				Enabled:               settings.Enabled,
				CanBeUsedByModerators: settings.CanBeUsedByModerators,
				TimeoutSeconds:        settings.TimeoutSeconds,
				DecisionSeconds:       settings.DecisionSeconds,
				TumberSize:            6,
				ChargedBullets:        settings.ChargedBullets,
				InitMessage:           settings.InitMessage,
				SurviveMessage:        settings.SurviveMessage,
				DeathMessage:          settings.DeathMessage,
			},
		)
		if err != nil {
			return err
		}

		// update settings
		_, err = tx.ExecContext(
			ctx,
			"UPDATE channels_modules_settings SET settings = ? WHERE id = ?",
			newSettingsBytes,
			id,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func downRouletteAddTumberSize(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
