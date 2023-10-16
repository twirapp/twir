package messages_updater

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/avast/retry-go/v4"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/switchupcb/disgo"
)

func (c *MessagesUpdater) processOffline(
	ctx context.Context,
	twitchChannelId string,
) error {
	// settings, err := c.getChannelDiscordIntegration(ctx, twitchChannelId)
	// if err != nil {
	// 	return err
	// }
	//
	// if settings.Data.Discord == nil || len(settings.Data.Discord.Guilds) == 0 {
	// 	return nil
	// }

	messages, err := c.store.GetByChannelId(ctx, twitchChannelId)
	if err != nil {
		return fmt.Errorf("failed to get messages: %w", err)
	}

	for _, message := range messages {
		settings, err := c.getChannelDiscordIntegration(ctx, twitchChannelId)
		if err != nil {
			return err
		}

		if settings.Data.Discord == nil || len(settings.Data.Discord.Guilds) == 0 {
			return nil
		}

		guild, ok := lo.Find(
			settings.Data.Discord.Guilds,
			func(guild model.ChannelIntegrationDataDiscordGuild) bool {
				return guild.ID == message.GuildID
			},
		)

		if !ok {
			continue
		}

		if guild.ShouldDeleteMessageOnOffline {
			delMsg := disgo.DeleteMessage{
				ChannelID: message.DiscordChannelID,
				MessageID: message.MessageID,
			}
			err = retry.Do(
				func() error {
					return delMsg.Send(c.discord.Client)
				},
				retry.Attempts(3),
			)

			if err != nil {
				c.logger.Error("Failed to delete message", slog.Any("err", err))
				continue
			}
		} else if guild.OfflineNotificationMessage != "" {
			content := &guild.OfflineNotificationMessage
			if *content == "" {
				content = lo.ToPtr("Stream is offline")
			}

			editMsg := disgo.EditMessage{
				ChannelID: message.DiscordChannelID,
				MessageID: message.MessageID,
				Content:   lo.ToPtr(content),
				Embeds:    &[]*disgo.Embed{},
			}

			_, err = retry.DoWithData(
				func() (*disgo.Message, error) {
					return editMsg.Send(c.discord.Client)
				},
				retry.Attempts(3),
			)

			if err != nil {
				c.logger.Error("Failed to edit message", slog.Any("err", err))
				continue
			}
		}
	}

	return nil
}
