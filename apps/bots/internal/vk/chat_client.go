package vk

import (
	"context"
	"fmt"

	buscore "github.com/twirapp/twir/libs/bus-core"
	buscoretokens "github.com/twirapp/twir/libs/bus-core/tokens"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	vkintegrations "github.com/twirapp/twir/libs/integrations/vk"
)

type botTokenRequester interface {
	Request(context.Context, buscoretokens.GetBotTokenRequest) (*buscore.QueueResponse[buscoretokens.TokenResponse], error)
}

type videoChatSender interface {
	SendTextMessage(context.Context, vkintegrations.SendTextMessageInput) error
}

type ChatClient struct {
	requestBotToken botTokenRequester
	videoChat       videoChatSender
}

func NewChatClient(twirBus *buscore.Bus, videoChat *vkintegrations.VideoChatClient) *ChatClient {
	return &ChatClient{
		requestBotToken: twirBus.Tokens.RequestBotToken,
		videoChat:       videoChat,
	}
}

func (c *ChatClient) SendMessage(
	ctx context.Context,
	binding channelplatformentity.ChannelPlatform,
	text string,
) error {
	if binding.BotUserID == nil {
		return nil
	}

	tokenResponse, err := c.requestBotToken.Request(
		ctx,
		buscoretokens.GetBotTokenRequest{Platform: platformentity.PlatformVKVideoLive},
	)
	if err != nil {
		return fmt.Errorf("request VK Video Live bot token: %w", err)
	}

	if err := c.videoChat.SendTextMessage(ctx, vkintegrations.SendTextMessageInput{
		AccessToken: tokenResponse.Data.AccessToken,
		OwnerID:     binding.PlatformChannelID,
		Content:     text,
	}); err != nil {
		return fmt.Errorf("send VK Video Live chat message: %w", err)
	}

	return nil
}
