package messages_updater

import (
	"fmt"
	"strings"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/switchupcb/disgo"
)

func (c *MessagesUpdater) buildEmbed(
	twitchUser helix.User,
	stream model.ChannelsStreams,
	guild model.ChannelIntegrationDataDiscordGuild,
) *disgo.Embed {
	twitchUrl := fmt.Sprintf("https://twitch.tv/%s", stream.UserLogin)
	width := 1920
	height := 1080
	thumbNailUrl := fmt.Sprintf("%s?t=%v", stream.ThumbnailUrl, time.Now().Unix())
	thumbNailUrl = strings.Replace(thumbNailUrl, "{width}", fmt.Sprintf("%d", width), 1)
	thumbNailUrl = strings.Replace(thumbNailUrl, "{height}", fmt.Sprintf("%d", height), 1)

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

	return &embed
}
