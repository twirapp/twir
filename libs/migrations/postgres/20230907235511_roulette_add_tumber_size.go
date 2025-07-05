package postgres

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
	InitMessage     string `json:"initMessage"`
	SurviveMessage  string `json:"surviveMessage"`
	DeathMessage    string `json:"deathMessage"`
	TimeoutSeconds  int    `json:"timeoutTime"`
	DecisionSeconds int    `json:"decisionTime"`
	ChargedBullets  int    `json:"chargedBullets"`

	Enabled               bool `json:"enabled"`
	CanBeUsedByModerators bool `json:"canBeUsedByModerator"`
}

type rouletteAddTumberSize20230907235511New struct {
	InitMessage     string `json:"initMessage"`
	SurviveMessage  string `json:"surviveMessage"`
	DeathMessage    string `json:"deathMessage"`
	TimeoutSeconds  int    `json:"timeoutTime"`
	DecisionSeconds int    `json:"decisionTime"`
	TumberSize      int    `json:"tumberSize"`
	ChargedBullets  int    `json:"chargedBullets"`

	Enabled               bool `json:"enabled"`
	CanBeUsedByModerators bool `json:"canBeUsedByModerator"`
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

	var forUpdate []struct {
		id            string
		settingsBytes []byte
	}

	// iterate over all modules settings
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
		var oldSettings rouletteAddTumberSize20230907235511Old
		if err := json.Unmarshal(update.settingsBytes, &oldSettings); err != nil {
			return err
		}

		newSettings := rouletteAddTumberSize20230907235511New{
			Enabled:               oldSettings.Enabled,
			CanBeUsedByModerators: oldSettings.CanBeUsedByModerators,
			TimeoutSeconds:        oldSettings.TimeoutSeconds,
			DecisionSeconds:       oldSettings.DecisionSeconds,
			TumberSize:            6,
			ChargedBullets:        oldSettings.ChargedBullets,
			InitMessage:           oldSettings.InitMessage,
			SurviveMessage:        oldSettings.SurviveMessage,
			DeathMessage:          oldSettings.DeathMessage,
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
	}

	return nil
}

func downRouletteAddTumberSize(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
