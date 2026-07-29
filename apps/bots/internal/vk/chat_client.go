package vk

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	vkintegrations "github.com/twirapp/twir/libs/integrations/vk"
	"github.com/twirapp/twir/libs/oauth"
)

type botTokenSource interface {
	Token(context.Context) (oauth.Credential, error)
}

type ownerTokenSource interface {
	Token(context.Context, uuid.UUID) (oauth.Credential, error)
}

type videoChatSender interface {
	SendTextMessage(context.Context, vkintegrations.SendTextMessageInput) error
}

type ChatClient struct {
	ownerTokens ownerTokenSource
	botTokens   botTokenSource
	videoChat   videoChatSender
}

func NewChatClient(ownerTokens ownerTokenSource, botTokens botTokenSource, videoChat videoChatSender) *ChatClient {
	return &ChatClient{
		ownerTokens: ownerTokens, botTokens: botTokens, videoChat: videoChat,
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

	ownerToken, err := c.ownerTokens.Token(ctx, binding.UserID)
	if err != nil {
		return fmt.Errorf("request VK Video Live binding owner token: %w", err)
	}

	botToken, err := c.botTokens.Token(ctx)
	if err != nil {
		return fmt.Errorf("request VK Video Live bot token: %w", err)
	}

	if err := c.videoChat.SendTextMessage(ctx, vkintegrations.SendTextMessageInput{
		OwnerAccessToken: ownerToken.AccessToken,
		BotAccessToken:   botToken.AccessToken,
		Content:          text,
	}); err != nil {
		return fmt.Errorf("send VK Video Live chat message: %w", err)
	}

	return nil
}
