package messages_updater

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/avast/retry-go/v4"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/samber/lo"
	model "github.com/twirapp/twir/libs/gomodels"
)

func (c *MessagesUpdater) processOffline(
	ctx context.Context,
	twitchChannelId string,
) error {
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

		gUid, _ := strconv.ParseUint(guild.ID, 10, 64)
		shard, _ := c.discord.FromGuildID(discord.GuildID(gUid))
		if shard == nil {
			c.logger.Error("Shard not found", slog.Any("guild_id", guild.ID))
			continue
		}

		dChanUid, err := strconv.ParseUint(message.DiscordChannelID, 10, 64)
		if err != nil {
			return err
		}

		dMsgId, err := strconv.ParseUint(message.MessageID, 10, 64)
		if err != nil {
			return err
		}

		twitchUser, err := c.getTwitchUser(message.TwitchChannelID)
		if err != nil {
			return err
		}

		if guild.ShouldDeleteMessageOnOffline {
			err = retry.Do(
				func() error {
					return shard.(*state.State).DeleteMessage(
						discord.ChannelID(dChanUid),
						discord.MessageID(dMsgId),
						"stream offline",
					)
				},
				retry.Attempts(3),
			)

			if err != nil {
				c.logger.Error("Failed to delete message", slog.Any("err", err))
				continue
			}
		} else {
			content := guild.OfflineNotificationMessage
			if content == "" {
				content = "{userName} is offline now"
			}

			content = c.replaceMessageVars(
				content,
				replaceMessageVarsOpts{
					UserName:    twitchUser.Login,
					DisplayName: twitchUser.DisplayName,
				},
			)

			_, err = retry.DoWithData(
				func() (*discord.Message, error) {
					return shard.(*state.State).EditMessageComplex(
						discord.ChannelID(dChanUid),
						discord.MessageID(dMsgId),
						api.EditMessageData{
							Content: option.NewNullableString(content),
							Embeds:  &[]discord.Embed{},
						},
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
