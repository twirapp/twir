package messages_updater

import (
	"context"

	"github.com/satont/twir/apps/discord/internal/sended_messages_store"
	model "github.com/satont/twir/libs/gomodels"
)

func (c *MessagesUpdater) getChannelDiscordIntegration(
	ctx context.Context,
	channelId string,
) (*model.ChannelsIntegrations, error) {
	integration := &model.ChannelsIntegrations{}
	err := c.db.WithContext(ctx).Where(`"channelId" = ?`, channelId).First(integration).Error
	return integration, err
}

func (c *MessagesUpdater) sendOnlineMessage(
	ctx context.Context,
	stream model.ChannelsStreams,
) ([]sended_messages_store.Message, error) {
	// settings, err := c.getChannelDiscordIntegration(ctx, stream.UserId)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// if !settings.Enabled {
	// 	return nil, errors.New("discord integration is disabled")
	// }
	//
	// twitchUsersReq, err := c.twitchClient.GetUsers(&helix.UsersParams{IDs: []string{stream.UserId}})
	// if len(twitchUsersReq.Data.Users) == 0 {
	// 	return nil, errors.New("user not found")
	// }
	// twitchUser := twitchUsersReq.Data.Users[0]
	//
	// var sendedMessage []sended_messages_store.Message
	//
	// twitchUrl := fmt.Sprintf("https://twitch.tv/%s", stream.UserLogin)
	// width := 1920
	// height := 1080
	// thumbNailUrl := fmt.Sprintf("%s?t=%v", stream.ThumbnailUrl, time.Now().Unix())
	// thumbNailUrl = strings.Replace(thumbNailUrl, "{width}", fmt.Sprintf("%d", width), 1)
	// thumbNailUrl = strings.Replace(thumbNailUrl, "{height}", fmt.Sprintf("%d", height), 1)
	//
	// embed := discordgo.MessageEmbed{
	// 	URL:       twitchUrl,
	// 	Title:     stream.Title,
	// 	Timestamp: time.Now().Format(time.RFC3339),
	// 	Color:     0x6441a5,
	// 	Footer: &discordgo.MessageEmbedFooter{
	// 		Text:    "TwirApp",
	// 		IconURL: fmt.Sprintf("https://twir.app/favicon.svg?t=%v", time.Now().Unix()),
	// 	},
	// 	Image: &discordgo.MessageEmbedImage{
	// 		URL:    thumbNailUrl,
	// 		Width:  width,
	// 		Height: height,
	// 	},
	// 	Thumbnail: &discordgo.MessageEmbedThumbnail{
	// 		URL: fmt.Sprintf(
	// 			"%s?t=%v",
	// 			twitchUser.ProfileImageURL,
	// 			time.Now().Unix(),
	// 		),
	// 	},
	// 	Fields: []*discordgo.MessageEmbedField{
	// 		{
	// 			Name:   "Category",
	// 			Value:  stream.GameName,
	// 			Inline: false,
	// 		},
	// 		{
	// 			Name:   "Viewers",
	// 			Value:  fmt.Sprintf("%d", stream.ViewerCount),
	// 			Inline: false,
	// 		},
	// 	},
	// }
	//
	// for _, channelId := range settings.Data.DiscordChannels {
	// 	m, err := c.discord.ChannelMessageSendEmbed(channelId, &embed)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	sendedMessage = append(
	// 		sendedMessage,
	// 		sended_messages_store.Message{
	// 			GuildID:   m.GuildID,
	// 			MessageID: m.ID,
	// 			ChannelID: m.ChannelID,
	// 		},
	// 	)
	// }
	//
	// return sendedMessage, nil
	return nil, nil
}
