package messages_updater

import (
	"context"
	"errors"
	"log/slog"
	"strconv"

	"github.com/avast/retry-go/v4"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/nicklaw5/helix/v2"
	model "github.com/satont/twir/libs/gomodels"
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

	twitchUsersReq, err := c.twitchClient.GetUsers(&helix.UsersParams{IDs: []string{stream.UserId}})
	if len(twitchUsersReq.Data.Users) == 0 {
		return errors.New("user not found")
	}
	twitchUser := twitchUsersReq.Data.Users[0]

	for _, guild := range settings.Data.Discord.Guilds {
		if !guild.LiveNotificationEnabled {
			continue
		}

		messages, err := c.store.GetByGuildId(ctx, guild.ID)
		if err != nil {
			return err
		}

		for _, message := range messages {
			embed := c.buildEmbed(twitchUser, stream, guild)

			gUid, _ := strconv.ParseUint(guild.ID, 10, 64)
			shard, _ := c.discord.FromGuildID(discord.GuildID(gUid))
			if shard == nil {
				c.logger.Error("Shard not found", slog.Any("guild_id", guild.ID))
				continue
			}

			_, err := retry.DoWithData(
				func() (*discord.Message, error) {
					dsChannelUid, err := strconv.ParseUint(message.DiscordChannelID, 10, 64)
					if err != nil {
						return nil, err
					}

					dMsgId, err := strconv.ParseUint(message.MessageID, 10, 64)
					if err != nil {
						return nil, err
					}

					return shard.(*state.State).EditMessage(
						discord.ChannelID(dsChannelUid),
						discord.MessageID(dMsgId),
						"",
						embed,
					)
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
