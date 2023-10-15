package messages_updater

import (
	"context"
	"errors"
	"log/slog"

	"github.com/avast/retry-go/v4"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/apps/discord/internal/sended_messages_store"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/switchupcb/disgo"
)

func (c *MessagesUpdater) sendOnlineMessage(
	ctx context.Context,
	stream model.ChannelsStreams,
) ([]sended_messages_store.Message, error) {
	settings, err := c.getChannelDiscordIntegration(ctx, stream.UserId)
	if err != nil {
		return nil, err
	}

	if settings.Data.Discord == nil || len(settings.Data.Discord.Guilds) == 0 {
		return nil, nil
	}

	twitchUsersReq, err := c.twitchClient.GetUsers(&helix.UsersParams{IDs: []string{stream.UserId}})
	if len(twitchUsersReq.Data.Users) == 0 {
		return nil, errors.New("user not found")
	}
	twitchUser := twitchUsersReq.Data.Users[0]

	var sendedMessage []sended_messages_store.Message

	for _, guild := range settings.Data.Discord.Guilds {
		if !guild.LiveNotificationEnabled {
			continue
		}

		embed := c.buildEmbed(twitchUser, stream, guild)

		for _, channel := range guild.LiveNotificationChannelsIds {
			sendMsgReq := disgo.CreateMessage{
				ChannelID: channel,
				Embeds:    []*disgo.Embed{embed},
			}

			if guild.LiveNotificationMessage != "" {
				sendMsgReq.Content = &guild.LiveNotificationMessage
			}

			m, err := retry.DoWithData(
				func() (*disgo.Message, error) {
					return sendMsgReq.Send(c.discord.Client)
				},
				retry.Attempts(3),
			)

			if err != nil {
				c.logger.Error("Failed to send message", slog.Any("err", err))
				continue
			}
			sendedMessage = append(
				sendedMessage,
				sended_messages_store.Message{
					MessageID:        m.ID,
					TwitchChannelID:  stream.UserId,
					GuildID:          guild.ID,
					DiscordChannelID: m.ChannelID,
				},
			)
		}
	}

	return sendedMessage, nil
}
