package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upDiscordIntegrationSeparateTable, downDiscordIntegrationSeparateTable)
}

func upDiscordIntegrationSeparateTable(ctx context.Context, tx *sql.Tx) error {
	type discordSeparateTableIntegrations20251201000000 struct {
		ID string
	}

	type discordSeparateTableChannelsIntegrations20251201000000 struct {
		ID        string
		ChannelID string
		Data      []byte
	}

	type discordGuild20251201000000 struct {
		ID                               string   `json:"id,omitempty"`
		LiveNotificationEnabled          bool     `json:"liveNotificationEnabled,omitempty"`
		LiveNotificationChannelsIds      []string `json:"liveNotificationChannelsIds,omitempty"`
		LiveNotificationShowTitle        bool     `json:"liveNotificationShowTitle,omitempty"`
		LiveNotificationShowCategory     bool     `json:"liveNotificationShowCategory,omitempty"`
		LiveNotificationShowViewers      bool     `json:"liveNotificationShowViewers,omitempty"`
		LiveNotificationMessage          string   `json:"liveNotificationMessage,omitempty"`
		LiveNotificationShowPreview      bool     `json:"liveNotificationShowPreview,omitempty"`
		LiveNotificationShowProfileImage bool     `json:"liveNotificationShowProfileImage,omitempty"`
		OfflineNotificationMessage       string   `json:"offlineNotificationMessage,omitempty"`
		ShouldDeleteMessageOnOffline     bool     `json:"shouldDeleteMessageOnOffline,omitempty"`
		AdditionalUsersIdsForLiveCheck   []string `json:"additionalUsersIdsForLiveCheck,omitempty"`
	}

	type discordData20251201000000 struct {
		Guilds []discordGuild20251201000000 `json:"guilds,omitempty"`
	}

	// Get discord integration ID
	discordIntegration := discordSeparateTableIntegrations20251201000000{}
	err := tx.QueryRowContext(
		ctx,
		"SELECT id FROM integrations WHERE service = 'DISCORD'",
	).Scan(&discordIntegration.ID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("get discord integration: %w", err)
	}

	// Create new table first
	createTableQuery := `
CREATE TABLE IF NOT EXISTS channels_integrations_discord (
	id SERIAL PRIMARY KEY,
	channel_id TEXT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
	guild_id TEXT NOT NULL,
	live_notification_enabled BOOLEAN NOT NULL DEFAULT false,
	live_notification_channels_ids TEXT[] NOT NULL DEFAULT '{}',
	live_notification_show_title BOOLEAN NOT NULL DEFAULT true,
	live_notification_show_category BOOLEAN NOT NULL DEFAULT true,
	live_notification_show_viewers BOOLEAN NOT NULL DEFAULT true,
	live_notification_message TEXT NOT NULL DEFAULT '',
	live_notification_show_preview BOOLEAN NOT NULL DEFAULT true,
	live_notification_show_profile_image BOOLEAN NOT NULL DEFAULT true,
	offline_notification_message TEXT NOT NULL DEFAULT '',
	should_delete_message_on_offline BOOLEAN NOT NULL DEFAULT false,
	additional_users_ids_for_live_check TEXT[] NOT NULL DEFAULT '{}',
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS channels_integrations_discord_channel_guild_idx ON channels_integrations_discord(channel_id, guild_id);
CREATE INDEX IF NOT EXISTS channels_integrations_discord_channel_id_idx ON channels_integrations_discord(channel_id);
`
	if _, err := tx.ExecContext(ctx, createTableQuery); err != nil {
		return fmt.Errorf("create table: %w", err)
	}

	// If no integration found, just exit
	if err == sql.ErrNoRows {
		return nil
	}

	// Get all discord integrations from old table
	rows, err := tx.QueryContext(
		ctx,
		`SELECT id, "channelId", data FROM channels_integrations WHERE "integrationId" = $1 AND data IS NOT NULL`,
		discordIntegration.ID,
	)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("get channels integrations: %w", err)
	}

	if rows == nil {
		return nil
	}
	defer rows.Close()

	channelsIntegrations := []discordSeparateTableChannelsIntegrations20251201000000{}
	for rows.Next() {
		channelIntegration := discordSeparateTableChannelsIntegrations20251201000000{}
		err := rows.Scan(
			&channelIntegration.ID,
			&channelIntegration.ChannelID,
			&channelIntegration.Data,
		)
		if err != nil {
			return fmt.Errorf("scan row: %w", err)
		}

		if channelIntegration.Data == nil {
			continue
		}

		channelsIntegrations = append(channelsIntegrations, channelIntegration)
	}

	if rows.Err() != nil {
		return fmt.Errorf("rows error: %w", rows.Err())
	}

	// Migrate data
	for _, channelIntegration := range channelsIntegrations {
		var parsedData discordData20251201000000
		if len(channelIntegration.Data) > 0 {
			err := json.Unmarshal(channelIntegration.Data, &parsedData)
			if err != nil {
				// Skip if we can't parse the data
				continue
			}
		}

		if parsedData.Guilds == nil {
			continue
		}

		for _, guild := range parsedData.Guilds {
			if guild.ID == "" {
				continue
			}

			insertQuery := `
INSERT INTO channels_integrations_discord (
	channel_id,
	guild_id,
	live_notification_enabled,
	live_notification_channels_ids,
	live_notification_show_title,
	live_notification_show_category,
	live_notification_show_viewers,
	live_notification_message,
	live_notification_show_preview,
	live_notification_show_profile_image,
	offline_notification_message,
	should_delete_message_on_offline,
	additional_users_ids_for_live_check
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
ON CONFLICT (channel_id, guild_id) DO NOTHING
`

			liveNotificationChannelsIds := guild.LiveNotificationChannelsIds
			if liveNotificationChannelsIds == nil {
				liveNotificationChannelsIds = []string{}
			}

			additionalUsersIds := guild.AdditionalUsersIdsForLiveCheck
			if additionalUsersIds == nil {
				additionalUsersIds = []string{}
			}

			_, err := tx.ExecContext(
				ctx,
				insertQuery,
				channelIntegration.ChannelID,
				guild.ID,
				guild.LiveNotificationEnabled,
				pq.Array(liveNotificationChannelsIds),
				guild.LiveNotificationShowTitle,
				guild.LiveNotificationShowCategory,
				guild.LiveNotificationShowViewers,
				guild.LiveNotificationMessage,
				guild.LiveNotificationShowPreview,
				guild.LiveNotificationShowProfileImage,
				guild.OfflineNotificationMessage,
				guild.ShouldDeleteMessageOnOffline,
				pq.Array(additionalUsersIds),
			)
			if err != nil {
				return fmt.Errorf("insert discord integration for channel %s guild %s: %w", channelIntegration.ChannelID, guild.ID, err)
			}
		}
	}

	return nil
}

func downDiscordIntegrationSeparateTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.ExecContext(ctx, "DROP TABLE IF EXISTS channels_integrations_discord")
	return err
}
