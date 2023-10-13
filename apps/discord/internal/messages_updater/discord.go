package messages_updater

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/discord/internal/sended_messages_store"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/switchupcb/disgo"
)

func (c *MessagesUpdater) getChannelDiscordIntegration(
	ctx context.Context,
	channelId string,
) (*model.ChannelsIntegrations, error) {
	discordIntegration := model.Integrations{}
	err := c.db.WithContext(ctx).Where(
		`service = ?`,
		model.IntegrationServiceDiscord,
	).First(&discordIntegration).Error
	if err != nil {
		return nil, err
	}

	integration := &model.ChannelsIntegrations{}
	err = c.db.WithContext(ctx).Where(
		`"channelId" = ? AND "integrationId" = ?`,
		channelId,
		discordIntegration.ID,
	).First(integration).Error
	return integration, err
}

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

	twitchUrl := fmt.Sprintf("https://twitch.tv/%s", stream.UserLogin)
	width := 1920
	height := 1080
	thumbNailUrl := fmt.Sprintf("%s?t=%v", stream.ThumbnailUrl, time.Now().Unix())
	thumbNailUrl = strings.Replace(thumbNailUrl, "{width}", fmt.Sprintf("%d", width), 1)
	thumbNailUrl = strings.Replace(thumbNailUrl, "{height}", fmt.Sprintf("%d", height), 1)

	for _, guild := range settings.Data.Discord.Guilds {
		if !guild.LiveNotificationEnabled {
			continue
		}

		embed := disgo.Embed{
			URL:       &twitchUrl,
			Title:     &stream.Title,
			Timestamp: lo.ToPtr(time.Now()),
			Color:     lo.ToPtr(0x6441a5),
			Footer: &disgo.EmbedFooter{
				Text:    "TwirApp",
				IconURL: lo.ToPtr(fmt.Sprintf("https://twir.app/favicon.svg?t=%v", time.Now().Unix())),
			},
			Image: &disgo.EmbedImage{
				URL:    thumbNailUrl,
				Width:  &width,
				Height: &height,
			},
			Thumbnail: &disgo.EmbedThumbnail{
				URL: fmt.Sprintf(
					"%s?t=%v",
					twitchUser.ProfileImageURL,
					time.Now().Unix(),
				),
			},
			Fields: []*disgo.EmbedField{},
		}

		if guild.LiveNotificationShowTitle {
			embed.Fields = append(
				embed.Fields, &disgo.EmbedField{
					Name:   "Title",
					Value:  stream.Title,
					Inline: lo.ToPtr(true),
				},
			)
		}

		if guild.LiveNotificationShowViewers {
			embed.Fields = append(
				embed.Fields, &disgo.EmbedField{
					Name:   "Viewers",
					Value:  fmt.Sprintf("%d", stream.ViewerCount),
					Inline: lo.ToPtr(true),
				},
			)
		}

		if guild.LiveNotificationShowCategory {
			embed.Fields = append(
				embed.Fields, &disgo.EmbedField{
					Name:   "Category",
					Value:  stream.GameName,
					Inline: lo.ToPtr(false),
				},
			)
		}

		for _, channel := range guild.LiveNotificationChannelsIds {
			sendMsgReq := disgo.CreateMessage{
				ChannelID: channel,
				Embeds:    []*disgo.Embed{&embed},
			}

			if guild.LiveNotificationMessage != "" {
				sendMsgReq.Content = &guild.LiveNotificationMessage
			}

			m, err := sendMsgReq.Send(c.discord.Client)
			if err != nil {
				c.logger.Error("Failed to send message", slog.Any("err", err))
				continue
			}
			sendedMessage = append(
				sendedMessage, sended_messages_store.Message{
					MessageID: m.ID,
					ChannelID: stream.UserId,
					GuildID:   guild.ID,
				},
			)
		}
	}

	return sendedMessage, nil
}
