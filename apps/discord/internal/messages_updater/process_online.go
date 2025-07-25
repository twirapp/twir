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

	settings, err := c.getChannelDiscordIntegration(ctx, stream.UserId)
	if err != nil {
		return err
	}

	if settings.Data.Discord == nil || len(settings.Data.Discord.Guilds) == 0 {
		return nil
	}

	twitchUser, err := c.getTwitchUser(stream.UserId)
	if err != nil {
		return err
	}

	var sendedMessage []sended_messages_store.Message

	for _, guild := range settings.Data.Discord.Guilds {
		if !guild.LiveNotificationEnabled {
			continue
		}

		embed := c.buildEmbed(twitchUser, stream, guild)

		gUid, _ := strconv.ParseUint(guild.ID, 10, 64)
		shard, _ := c.discord.FromGuildID(discord.GuildID(gUid))
		if shard == nil {
			c.logger.Error("Shard not found", slog.Any("guild_id", guild.ID))
			continue
		}

		for _, channel := range guild.LiveNotificationChannelsIds {
			dChanUid, err := strconv.ParseUint(channel, 10, 64)
			if err != nil {
				c.logger.Error("Failed to parse channel id", slog.Any("err", err))
				continue
			}

			message := c.replaceMessageVars(
				guild.LiveNotificationMessage, replaceMessageVarsOpts{
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
				c.logger.Error("Failed to send message", slog.Any("err", err))
				continue
			}
			sendedMessage = append(
				sendedMessage,
				sended_messages_store.Message{
					MessageID:        m.ID.String(),
					TwitchChannelID:  stream.UserId,
					GuildID:          guild.ID,
					DiscordChannelID: m.ChannelID.String(),
				},
			)
		}
	}

	return nil
}
