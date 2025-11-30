package messages_updater

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/avast/retry-go/v4"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/twirapp/twir/apps/discord/internal/sended_messages_store"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
)

func (c *MessagesUpdater) processOnline(
	ctx context.Context,
	channelId string,
) error {
	stream := model.ChannelsStreams{}
	err := c.db.
		Where(`"userId" = ?`, channelId).
		First(&stream).
		Error
	if err != nil {
		return err
	}

	integrations, err := c.getChannelDiscordIntegrations(ctx, stream.UserId)
	if err != nil {
		return err
	}

	if len(integrations) == 0 {
		return nil
	}

	twitchUser, err := c.getTwitchUser(stream.UserId)
	if err != nil {
		return err
	}

	var sendedMessage []sended_messages_store.Message

	for _, integration := range integrations {
		if !integration.LiveNotificationEnabled {
			continue
		}

		embed := c.buildEmbed(twitchUser, stream, integration)

		gUid, _ := strconv.ParseUint(integration.GuildID, 10, 64)
		shard, _ := c.discord.FromGuildID(discord.GuildID(gUid))
		if shard == nil {
			c.logger.Error("Shard not found", slog.Any("guild_id", integration.GuildID))
			continue
		}

		for _, channel := range integration.LiveNotificationChannelsIds {
			dChanUid, err := strconv.ParseUint(channel, 10, 64)
			if err != nil {
				c.logger.Error("Failed to parse channel id", logger.Error(err))
				continue
			}

			message := c.replaceMessageVars(
				integration.LiveNotificationMessage, replaceMessageVarsOpts{
					UserName:     stream.UserLogin,
					DisplayName:  stream.UserName,
					CategoryName: stream.GameName,
					Title:        stream.Title,
				},
			)

			m, err := retry.DoWithData(
				func() (*discord.Message, error) {
					return shard.(*state.State).SendMessage(
						discord.ChannelID(dChanUid),
						message,
						embed,
					)
				},
				retry.Attempts(3),
			)

			if err != nil {
				c.logger.Error("Failed to send message", logger.Error(err))
				continue
			}
			sendedMessage = append(
				sendedMessage,
				sended_messages_store.Message{
					MessageID:        m.ID.String(),
					TwitchChannelID:  stream.UserId,
					GuildID:          integration.GuildID,
					DiscordChannelID: m.ChannelID.String(),
				},
			)
		}
	}

	return nil
}
