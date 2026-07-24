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
	requester := &fakeBotTokenRequester{
		response: &buscore.QueueResponse[buscoretokens.TokenResponse]{
			Data: buscoretokens.TokenResponse{AccessToken: "vk-video-bot-token"},
		},
	}
	videoChat := &fakeVideoChatClient{}
	client := &ChatClient{requestBotToken: requester, videoChat: videoChat}
	binding := channelplatformentity.ChannelPlatform{
		Platform:          platformentity.PlatformVKVideoLive,
		PlatformChannelID: "provider-channel-42",
		BotUserID:         &botUserID,
	}

	// When
	err := client.SendMessage(context.Background(), binding, "hello from bot")

	// Then
	require.NoError(t, err)
	require.Equal(t, 1, requester.calls)
	require.Equal(t, buscoretokens.GetBotTokenRequest{Platform: platformentity.PlatformVKVideoLive}, requester.request)
	require.Equal(t, []vkintegrations.SendTextMessageInput{{
		AccessToken: "vk-video-bot-token",
		OwnerID:     "provider-channel-42",
		Content:     "hello from bot",
	}}, videoChat.inputs)
}

func TestChatClientSkipsSendWhenBindingHasNoBotUser(t *testing.T) {
	// Given
	requester := &fakeBotTokenRequester{}
	videoChat := &fakeVideoChatClient{}
	client := &ChatClient{requestBotToken: requester, videoChat: videoChat}
	binding := channelplatformentity.ChannelPlatform{
		Platform:          platformentity.PlatformVKVideoLive,
		PlatformChannelID: "provider-channel-42",
	}

	// When
	err := client.SendMessage(context.Background(), binding, "hello from bot")

	// Then
	require.NoError(t, err)
	require.Zero(t, requester.calls)
	require.Empty(t, videoChat.inputs)
}

type fakeBotTokenRequester struct {
	calls    int
	request  buscoretokens.GetBotTokenRequest
	response *buscore.QueueResponse[buscoretokens.TokenResponse]
	err      error
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
