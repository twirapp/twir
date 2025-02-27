package migrations

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upDuelToTable, downDuelToTable)
}

type DuelToColumns struct {
	ID        string
	ChannelId string `gorm:"column:channelId;type:text" json:"channelId"`
	Settings  []byte
}

type DuelToColumnsSettings struct {
	StartMessage    string `json:"start_message"`
	ResultMessage   string `json:"result_message"`
	BothDieMessage  string `json:"both_die_message"`
	UserCooldown    int32  `json:"user_cooldown"`
	GlobalCooldown  int32  `json:"global_cooldown"`
	TimeoutSeconds  int32  `json:"timeout_seconds"`
	SecondsToAccept int32  `json:"seconds_to_accept"`
	PointsPerWin    int32  `json:"points_per_win"`
	PointsPerLose   int32  `json:"points_per_lose"`
	BothDiePercent  int32  `json:"both_die_percent"`
	Enabled         bool   `json:"enabled"`
}

func upDuelToTable(ctx context.Context, tx *sql.Tx) error {
	var entities []DuelToColumns

	rows, err := tx.QueryContext(
		ctx,
		`SELECT id, settings, "channelId" FROM channels_modules_settings WHERE type = 'duel'`,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var entity DuelToColumns
		if err := rows.Scan(&entity.ID, &entity.Settings, &entity.ChannelId); err != nil {
			return err
		}
		entities = append(entities, entity)
	}

	_, err = tx.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS channels_games_duel (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			channel_id text references channels(id),
			enabled BOOLEAN NOT NULL default false,
			user_cooldown INT NOT NULL default 0,
			global_cooldown INT NOT NULL default 0,
			timeout_seconds INT NOT NULL default 0,
			start_message varchar(500) NOT NULL default '',
			result_message varchar(500) NOT NULL default '',
			seconds_to_accept INT NOT NULL default 0,
			points_per_win INT NOT NULL default 0,
			points_per_lose INT NOT NULL default 0,
			both_die_percent INT NOT NULL default 0,
			both_die_message varchar(500) NOT NULL default ''
		)`,
	)
	if err != nil {
		return err
	}

	for _, entity := range entities {
		var settings DuelToColumnsSettings
		if err := json.Unmarshal(entity.Settings, &settings); err != nil {
			return err
		}

		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO channels_games_duel (
					 id,
					 channel_id,
					 enabled,
					 user_cooldown,
					 global_cooldown,
					 timeout_seconds,
					 start_message,
					 result_message,
					 seconds_to_accept,
					 points_per_win,
					 points_per_lose,
					 both_die_percent,
					 both_die_message
					 ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
			entity.ID,
			entity.ChannelId,
			settings.Enabled,
			settings.UserCooldown,
			settings.GlobalCooldown,
			settings.TimeoutSeconds,
			settings.StartMessage,
			settings.ResultMessage,
			settings.SecondsToAccept,
			settings.PointsPerWin,
			settings.PointsPerLose,
			settings.BothDiePercent,
			settings.BothDieMessage,
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

	return nil
}

func downDuelToTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
