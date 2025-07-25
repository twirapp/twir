package messages_updater

import (
	"fmt"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/nicklaw5/helix/v2"
	model "github.com/twirapp/twir/libs/gomodels"
)

func (c *MessagesUpdater) buildEmbed(
	twitchUser helix.User,
	stream model.ChannelsStreams,
	guild model.ChannelIntegrationDataDiscordGuild,
) discord.Embed {
	twitchUrl := fmt.Sprintf("https://twitch.tv/%s", stream.UserLogin)
	var width uint = 1920
	var height uint = 1080
	thumbNailUrl := fmt.Sprintf("%s?t=%v", stream.ThumbnailUrl, time.Now().Unix())
	thumbNailUrl = strings.Replace(thumbNailUrl, "{width}", fmt.Sprintf("%d", width), 1)
	thumbNailUrl = strings.Replace(thumbNailUrl, "{height}", fmt.Sprintf("%d", height), 1)

	embed := discord.Embed{
		URL:       twitchUrl,
		Title:     stream.Title,
		Timestamp: discord.NewTimestamp(time.Now()),
		Color:     0x6441a5,
		Footer: &discord.EmbedFooter{
			Text: "TwirApp",
			Icon: fmt.Sprintf("https://twir.app/favicon.svg?t=%v", time.Now().Unix()),
		},
		Fields: []discord.EmbedField{},
	}

	if guild.LiveNotificationShowPreview {
		embed.Image = &discord.EmbedImage{
			URL:    thumbNailUrl,
			Width:  width,
			Height: height,
		}
	}

	if guild.LiveNotificationShowProfileImage {
		embed.Thumbnail = &discord.EmbedThumbnail{
			URL: fmt.Sprintf(
				"%s?t=%v",
				twitchUser.ProfileImageURL,
				time.Now().Unix(),
			),
		}
	}

	if guild.LiveNotificationShowTitle {
		embed.Fields = append(
			embed.Fields,
			discord.EmbedField{
				Name:   "Title",
				Value:  stream.Title,
				Inline: true,
			},
		)
	}

	if guild.LiveNotificationShowViewers {
		embed.Fields = append(
			embed.Fields,
			discord.EmbedField{
				Name:   "Viewers",
				Value:  fmt.Sprintf("%d", stream.ViewerCount),
				Inline: true,
			},
		)
	}

	if guild.LiveNotificationShowCategory {
		embed.Fields = append(
			embed.Fields,
			discord.EmbedField{
				Name:   "Category",
				Value:  stream.GameName,
				Inline: false,
			},
		)
	}

	return embed
}
