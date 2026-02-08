package discordmessagesupdater

import (
	"context"
	"log/slog"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/twirapp/twir/apps/bots/internal/discord/sended_messages_store"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
	discordmodel "github.com/twirapp/twir/libs/repositories/channels_integrations_discord/model"
)

const (
	// Check interval for periodic updates
	checkInterval = 30 * time.Second
	// Message update interval - how often to update each message
	messageUpdateInterval = 5 * time.Minute
)

// StartPeriodicUpdater starts a goroutine that periodically checks for messages that need updating
func (c *MessagesUpdater) StartPeriodicUpdater(ctx context.Context) {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	c.logger.InfoContext(
		ctx,
		"Starting Discord messages periodic updater",
		slog.Duration("interval", checkInterval),
	)

	for {
		select {
		case <-ctx.Done():
			c.logger.InfoContext(ctx, "Stopping Discord messages periodic updater")
			return
		case <-ticker.C:
			if err := c.processPeriodicUpdates(ctx); err != nil {
				c.logger.ErrorContext(
					ctx,
					"Failed to process periodic updates",
					logger.Error(err),
				)
			}
		}
	}
}

// processPeriodicUpdates checks all live streams and updates messages that need updating
func (c *MessagesUpdater) processPeriodicUpdates(ctx context.Context) error {
	// Get all active streams
	var streams []model.ChannelsStreams
	if err := c.db.Find(&streams).Error; err != nil {
		return err
	}

	if len(streams) == 0 {
		return nil
	}

	c.logger.DebugContext(
		ctx,
		"Processing periodic updates",
		slog.Int("active_streams", len(streams)),
	)

	for _, stream := range streams {
		if err := c.updateMessagesForStream(ctx, stream); err != nil {
			c.logger.ErrorContext(
				ctx,
				"Failed to update messages for stream",
				logger.Error(err),
				slog.String("channel_id", stream.UserId),
			)
		}
	}

	return nil
}

// updateMessagesForStream updates Discord messages for a specific stream if interval has passed
func (c *MessagesUpdater) updateMessagesForStream(
	ctx context.Context,
	stream model.ChannelsStreams,
) error {
	integrations, err := c.getChannelDiscordIntegrations(ctx, stream.UserId)
	if err != nil {
		return err
	}

	if len(integrations) == 0 {
		return nil
	}

	for _, integration := range integrations {
		if !integration.LiveNotificationEnabled {
			continue
		}

		messages, err := c.store.GetByGuildId(ctx, integration.GuildID)
		if err != nil {
			c.logger.ErrorContext(
				ctx,
				"Failed to get messages for guild",
				logger.Error(err),
				slog.String("guild_id", integration.GuildID),
			)
			continue
		}

		for _, message := range messages {
			if message.TwitchChannelID != stream.UserId {
				continue
			}

			// Check if enough time has passed since last update
			if message.LastUpdatedAt != nil {
				timeSinceUpdate := time.Since(*message.LastUpdatedAt)

				if timeSinceUpdate < messageUpdateInterval {
					continue
				}
			}

			// Update the message
			if err := c.updateSingleMessage(ctx, &message, stream, integration); err != nil {
				c.logger.ErrorContext(
					ctx,
					"Failed to update single message",
					logger.Error(err),
					slog.String("message_id", message.MessageID),
					slog.String("guild_id", integration.GuildID),
				)
			}
		}
	}

	return nil
}

// updateSingleMessage updates a single Discord message with current stream data
func (c *MessagesUpdater) updateSingleMessage(
	ctx context.Context,
	message *sended_messages_store.Message,
	stream model.ChannelsStreams,
	integration discordmodel.ChannelIntegrationDiscord,
) error {
	twitchUser, err := c.getTwitchUser(stream.UserId)
	if err != nil {
		return err
	}

	embed := c.buildEmbed(twitchUser, stream, integration)

	err = retry.Do(
		func() error {
			content := c.replaceMessageVars(
				integration.LiveNotificationMessage, replaceMessageVarsOpts{
					UserName:     stream.UserLogin,
					DisplayName:  stream.UserName,
					CategoryName: stream.GameName,
					Title:        stream.Title,
				},
			)

			return c.discord.EditMessage(
				ctx,
				message.DiscordChannelID,
				message.MessageID,
				content,
				embed,
			)
		},
		retry.Attempts(3),
	)

	if err != nil {
		return err
	}

	// Update last_updated_at timestamp
	if err := c.store.UpdateLastUpdatedAt(ctx, message.MessageID); err != nil {
		c.logger.ErrorContext(
			ctx,
			"Failed to update last_updated_at timestamp",
			logger.Error(err),
			slog.String("message_id", message.MessageID),
		)
	}

	return nil
}
