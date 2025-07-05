package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upEightBallToTable, downEightBallToTable)
}

type EightBallToColumns struct {
	ID        string
	ChannelId string `gorm:"column:channelId;type:text" json:"channelId"`
	Settings  []byte
}

type EightBallToColumnsSettings struct {
	Answers []string `json:"answers"`
	Enabled bool     `json:"enabled"`
}

func upEightBallToTable(ctx context.Context, tx *sql.Tx) error {
	var entities []EightBallToColumns

	rows, err := tx.QueryContext(
		ctx,
		`SELECT id, settings, "channelId" FROM channels_modules_settings WHERE type = '8ball'`,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var entity EightBallToColumns
		if err := rows.Scan(&entity.ID, &entity.Settings, &entity.ChannelId); err != nil {
			return err
		}
		entities = append(entities, entity)
	}

	_, err = tx.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS channels_games_8ball (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			channel_id text references channels(id),
			enabled BOOLEAN NOT NULL default false,
			answers TEXT[] NOT NULL default '{}'
		)`,
	)
	if err != nil {
		return err
	}

	for _, entity := range entities {
		var settings EightBallToColumnsSettings
		if err := json.Unmarshal(entity.Settings, &settings); err != nil {
			return err
		}

		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO channels_games_8ball (id, channel_id, enabled, answers) VALUES ($1, $2, $3, $4)`,
			entity.ID,
			entity.ChannelId,
			settings.Enabled,
			append(pq.StringArray{}, settings.Answers...),
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

func downEightBallToTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
