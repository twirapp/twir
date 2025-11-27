package messages_updater

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/avast/retry-go/v4"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
)

func (c *MessagesUpdater) updateDiscordMessages(
	ctx context.Context,
	stream model.ChannelsStreams,
) error {
	settings, err := c.getChannelDiscordIntegration(ctx, stream.UserId)
	if err != nil {
		return err
	}

	if settings.Data.Discord == nil || len(settings.Data.Discord.Guilds) == 0 {
		return nil
	}

	for _, guild := range settings.Data.Discord.Guilds {
		if !guild.LiveNotificationEnabled {
			continue
		}

		messages, err := c.store.GetByGuildId(ctx, guild.ID)
		if err != nil {
			return err
		}

		for _, message := range messages {
			if message.TwitchChannelID != stream.UserId {
				continue
			}

			twitchUser, err := c.getTwitchUser(stream.UserId)
			if err != nil {
				c.logger.Error("Failed to get twitch user", logger.Error(err))
				continue
			}

			embed := c.buildEmbed(twitchUser, stream, guild)

			gUid, _ := strconv.ParseUint(guild.ID, 10, 64)
			shard, _ := c.discord.FromGuildID(discord.GuildID(gUid))
			if shard == nil {
				c.logger.Error("Shard not found", slog.Any("guild_id", guild.ID))
				continue
			}

			_, err = retry.DoWithData(
				func() (*discord.Message, error) {
					dsChannelUid, err := strconv.ParseUint(message.DiscordChannelID, 10, 64)
					if err != nil {
						return nil, err
					}

					dMsgId, err := strconv.ParseUint(message.MessageID, 10, 64)
					if err != nil {
						return nil, err
					}

					content := c.replaceMessageVars(
						guild.LiveNotificationMessage, replaceMessageVarsOpts{
							UserName:     stream.UserLogin,
							DisplayName:  stream.UserName,
							CategoryName: stream.GameName,
							Title:        stream.Title,
						},
					)

					return shard.(*state.State).EditMessage(
						discord.ChannelID(dsChannelUid),
						discord.MessageID(dMsgId),
						content,
						embed,
					)
				},
				retry.Attempts(3),
			)

			if err != nil {
				c.logger.Error("Failed to edit message", logger.Error(err))
				continue
			}
		}
	}

	return nil
}
