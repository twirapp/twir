package vk

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	buscore "github.com/twirapp/twir/libs/bus-core"
	buscoretokens "github.com/twirapp/twir/libs/bus-core/tokens"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	vkintegrations "github.com/twirapp/twir/libs/integrations/vk"
)

func TestChatClientSendsTextWithGlobalVKVideoToken(t *testing.T) {
	// Given
	botUserID := uuid.New()
	ownerID := uuid.New()
	ownerRequester := &fakeUserTokenRequester{
		response: &buscore.QueueResponse[buscoretokens.TokenResponse]{
			Data: buscoretokens.TokenResponse{AccessToken: "vk-video-owner-token"},
		},
	}
	botRequester := &fakeBotTokenRequester{
		response: &buscore.QueueResponse[buscoretokens.TokenResponse]{
			Data: buscoretokens.TokenResponse{AccessToken: "vk-video-bot-token"},
		},
	}
	videoChat := &fakeVideoChatClient{}
	client := &ChatClient{requestOwnerToken: ownerRequester, requestBotToken: botRequester, videoChat: videoChat}
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
	require.Equal(t, 1, ownerRequester.calls)
	require.Equal(t, buscoretokens.GetUserTokenRequest{UserId: ownerID}, ownerRequester.request)
	require.Equal(t, 1, botRequester.calls)
	require.Equal(t, buscoretokens.GetBotTokenRequest{Platform: platformentity.PlatformVKVideoLive}, botRequester.request)
	require.Equal(t, []vkintegrations.SendTextMessageInput{{
		OwnerAccessToken: "vk-video-owner-token",
		BotAccessToken:   "vk-video-bot-token",
		Content:          "hello from bot",
	}}, videoChat.inputs)
}

func TestChatClientSkipsSendWhenBindingHasNoBotUser(t *testing.T) {
	// Given
	ownerRequester := &fakeUserTokenRequester{}
	botRequester := &fakeBotTokenRequester{}
	videoChat := &fakeVideoChatClient{}
	client := &ChatClient{requestOwnerToken: ownerRequester, requestBotToken: botRequester, videoChat: videoChat}
	binding := channelplatformentity.ChannelPlatform{
		Platform:          platformentity.PlatformVKVideoLive,
		PlatformChannelID: "provider-channel-42",
	}

	// When
	err := client.SendMessage(context.Background(), binding, "hello from bot")

	// Then
	require.NoError(t, err)
	require.Zero(t, ownerRequester.calls)
	require.Zero(t, botRequester.calls)
	require.Empty(t, videoChat.inputs)
}

type fakeBotTokenRequester struct {
	calls    int
	request  buscoretokens.GetBotTokenRequest
	response *buscore.QueueResponse[buscoretokens.TokenResponse]
	err      error
}

type fakeUserTokenRequester struct {
	calls    int
	request  buscoretokens.GetUserTokenRequest
	response *buscore.QueueResponse[buscoretokens.TokenResponse]
	err      error
}

func (f *fakeUserTokenRequester) Request(
	_ context.Context,
	request buscoretokens.GetUserTokenRequest,
) (*buscore.QueueResponse[buscoretokens.TokenResponse], error) {
	f.calls++
	f.request = request
	return f.response, f.err
}

func (f *fakeBotTokenRequester) Request(
	_ context.Context,
	request buscoretokens.GetBotTokenRequest,
) (*buscore.QueueResponse[buscoretokens.TokenResponse], error) {
	f.calls++
	f.request = request
	return f.response, f.err
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
