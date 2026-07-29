package vk

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	vkintegrations "github.com/twirapp/twir/libs/integrations/vk"
	"github.com/twirapp/twir/libs/oauth"
)

func TestChatClientSendsTextWithGlobalVKVideoToken(t *testing.T) {
	// Given
	botUserID := uuid.New()
	ownerID := uuid.New()
	ownerSource := &fakeUserTokenSource{credential: oauth.Credential{AccessToken: "vk-video-owner-token"}}
	botSource := &fakeBotTokenSource{credential: oauth.Credential{AccessToken: "vk-video-bot-token"}}
	videoChat := &fakeVideoChatClient{}
	client := &ChatClient{ownerTokens: ownerSource, botTokens: botSource, videoChat: videoChat}
	binding := channelplatformentity.ChannelPlatform{
		Platform:          platformentity.PlatformVKVideoLive,
		UserID:            ownerID,
		PlatformChannelID: "provider-channel-42",
		BotUserID:         &botUserID,
	}

	// When
	err := client.SendMessage(context.Background(), binding, "hello from bot")

	// Then
	require.NoError(t, err)
	require.Equal(t, 1, ownerSource.calls)
	require.Equal(t, ownerID, ownerSource.userID)
	require.Equal(t, 1, botSource.calls)
	require.Equal(t, []vkintegrations.SendTextMessageInput{{
		OwnerAccessToken: "vk-video-owner-token",
		BotAccessToken:   "vk-video-bot-token",
		Content:          "hello from bot",
	}}, videoChat.inputs)
}

func TestChatClientSkipsSendWhenBindingHasNoBotUser(t *testing.T) {
	// Given
	ownerSource := &fakeUserTokenSource{}
	botSource := &fakeBotTokenSource{}
	videoChat := &fakeVideoChatClient{}
	client := &ChatClient{ownerTokens: ownerSource, botTokens: botSource, videoChat: videoChat}
	binding := channelplatformentity.ChannelPlatform{
		Platform:          platformentity.PlatformVKVideoLive,
		PlatformChannelID: "provider-channel-42",
	}

	// When
	err := client.SendMessage(context.Background(), binding, "hello from bot")

	// Then
	require.NoError(t, err)
	require.Zero(t, ownerSource.calls)
	require.Zero(t, botSource.calls)
	require.Empty(t, videoChat.inputs)
}

type fakeBotTokenSource struct {
	calls      int
	credential oauth.Credential
	err        error
}

func (f *fakeBotTokenSource) Token(context.Context) (oauth.Credential, error) {
	f.calls++
	return f.credential, f.err
}

type fakeUserTokenSource struct {
	calls      int
	userID     uuid.UUID
	credential oauth.Credential
	err        error
}

func (f *fakeUserTokenSource) Token(_ context.Context, userID uuid.UUID) (oauth.Credential, error) {
	f.calls++
	f.userID = userID
	return f.credential, f.err
}

type fakeVideoChatClient struct {
	inputs []vkintegrations.SendTextMessageInput
	err    error
}

func (f *fakeVideoChatClient) SendTextMessage(
	_ context.Context,
	input vkintegrations.SendTextMessageInput,
) error {
	f.inputs = append(f.inputs, input)
	return f.err
}
