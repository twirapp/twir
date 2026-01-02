package discordmessagesupdater

import (
	"context"

	"github.com/avast/retry-go/v4"
	"github.com/twirapp/twir/apps/bots/internal/discord/discord_go"
	"github.com/twirapp/twir/apps/bots/internal/discord/sended_messages_store"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
)

func (c *MessagesUpdater) ProcessOnline(
	ctx context.Context,
	twitchChannelId string,
) error {
	stream := model.ChannelsStreams{}
	err := c.db.
		Where(`"userId" = ?`, twitchChannelId).
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
		for _, channel := range integration.LiveNotificationChannelsIds {
			message := c.replaceMessageVars(
				integration.LiveNotificationMessage, replaceMessageVarsOpts{
					UserName:     stream.UserLogin,
					DisplayName:  stream.UserName,
					CategoryName: stream.GameName,
					Title:        stream.Title,
				},
			)

			m, err := retry.DoWithData(
				func() (discord_go.SendMessageResponse, error) {
					return c.discord.SendMessage(
						ctx,
						channel,
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
					MessageID:        m.MessageID,
					TwitchChannelID:  stream.UserId,
					GuildID:          integration.GuildID,
					DiscordChannelID: channel,
				},
			)
		}
	}

	if err := c.store.Add(ctx, sendedMessage...); err != nil {
		c.logger.Error("Failed to store sent messages", logger.Error(err))
	}

	return nil
}
