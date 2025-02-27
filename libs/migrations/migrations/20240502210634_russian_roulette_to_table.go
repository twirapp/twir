package migrations

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upRussianRouletteToTable, downRussianRouletteToTable)
}

type RussianRouletteToTable struct {
	ID        string
	ChannelId string `gorm:"column:channelId;type:text" json:"channelId"`
	Settings  []byte
}

type RussianRouletteToTableSettings struct {
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

func upRussianRouletteToTable(ctx context.Context, tx *sql.Tx) error {
	var entities []RussianRouletteToTable

	rows, err := tx.QueryContext(
		ctx,
		`SELECT id, settings, "channelId" FROM channels_modules_settings WHERE type = 'russian_roulette'`,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var entity RussianRouletteToTable
		if err := rows.Scan(&entity.ID, &entity.Settings, &entity.ChannelId); err != nil {
			return err
		}
		entities = append(entities, entity)
	}

	_, err = tx.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS channels_games_russian_roulette (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			channel_id text references channels(id),
			enabled BOOLEAN NOT NULL default false,
			can_be_used_by_moderators BOOLEAN NOT NULL default false,
			timeout_seconds INT NOT NULL default 0,
			decision_seconds INT NOT NULL default 0,
			tumber_size INT NOT NULL default 0,
			charged_bullets INT NOT NULL default 0,
			init_message varchar(500) NOT NULL default '',
			survive_message varchar(500) NOT NULL default '',
			death_message varchar(500) NOT NULL default ''
		)`,
	)
	if err != nil {
		return err
	}

	for _, entity := range entities {
		var settings RussianRouletteToTableSettings
		if err := json.Unmarshal(entity.Settings, &settings); err != nil {
			return err
		}

		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO channels_games_russian_roulette (
				channel_id, enabled, can_be_used_by_moderators, timeout_seconds, decision_seconds, tumber_size, charged_bullets, init_message, survive_message, death_message
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			entity.ChannelId,
			settings.Enabled,
			settings.CanBeUsedByModerators,
			settings.TimeoutSeconds,
			settings.DecisionSeconds,
			settings.TumberSize,
			settings.ChargedBullets,
			settings.InitMessage,
			settings.SurviveMessage,
			settings.DeathMessage,
		)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(
			ctx,
			`DELETE FROM channels_modules_settings WHERE id = $1`,
			entity.ID,
		)
	}

	// This code is executed when the migration is applied.
	return nil
}

func downRussianRouletteToTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
