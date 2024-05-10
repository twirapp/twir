package migrations

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upChannelRolesSettingsToColumns, downChannelRolesSettingsToColumns)
}

type ChannelRoleChannelRolesSettingsToColumns struct {
	ID       string
	Settings []byte
}

type ChannelRoleSettingsChannelRolesSettingsToColumns struct {
	RequiredWatchTime         int64 `json:"requiredWatchTime"`
	RequiredMessages          int32 `json:"requiredMessages"`
	RequiredUsedChannelPoints int64 `json:"requiredUsedChannelPoints"`
}

func upChannelRolesSettingsToColumns(ctx context.Context, tx *sql.Tx) error {
	var roles []ChannelRoleChannelRolesSettingsToColumns

	rows, err := tx.QueryContext(ctx, "SELECT id, settings FROM channels_roles")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var role ChannelRoleChannelRolesSettingsToColumns
		if err := rows.Scan(&role.ID, &role.Settings); err != nil {
			return err
		}
		roles = append(roles, role)
	}

	_, err = tx.QueryContext(ctx, "ALTER TABLE channels_roles DROP COLUMN settings")
	if err != nil {
		return err
	}

	_, err = tx.QueryContext(ctx, "ALTER TABLE channels_roles ADD COLUMN required_watch_time bigint")
	if err != nil {
		return err
	}

	_, err = tx.QueryContext(ctx, "ALTER TABLE channels_roles ADD COLUMN required_messages integer")
	if err != nil {
		return err
	}

	_, err = tx.QueryContext(
		ctx,
		"ALTER TABLE channels_roles ADD COLUMN required_used_channel_points bigint",
	)
	if err != nil {
		return err
	}

	for _, role := range roles {
		var settings ChannelRoleSettingsChannelRolesSettingsToColumns
		if err := json.Unmarshal(role.Settings, &settings); err != nil {
			return err
		}

		_, err = tx.ExecContext(
			ctx,
			"UPDATE channels_roles SET required_watch_time = $1, required_messages = $2, required_used_channel_points = $3 WHERE id = $4",
			settings.RequiredWatchTime,
			settings.RequiredMessages,
			settings.RequiredUsedChannelPoints,
			role.ID,
		)
		if err != nil {
			return err
		}
	}

	// This code is executed when the migration is applied.
	return nil
}

func downChannelRolesSettingsToColumns(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
