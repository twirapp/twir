package platforms

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	"github.com/twirapp/twir/libs/entities/platform"
)

func TestVKVideoLiveChatAdapterSendsPlainText(t *testing.T) {
	// Given
	sender := &fakeVKVideoLiveSender{}
	adapter := NewVKVideoLiveChatAdapter(sender)
	binding := channelplatformentity.ChannelPlatform{Platform: platform.PlatformVKVideoLive}

	// When
	err := adapter.SendMessage(context.Background(), binding, "hello", "reply-id", ChatOptions{IsAnnounce: true})

	// Then
	require.NoError(t, err)
	require.Equal(t, binding, sender.binding)
	require.Equal(t, "hello", sender.message)
	require.Equal(t, platform.Capabilities{platform.CapabilityChatWrite}, adapter.Capabilities())
}

func TestDispatchSkipsDisabledVKVideoLiveBinding(t *testing.T) {
	// Given
	sender := &fakeVKVideoLiveSender{}
	adapter := NewVKVideoLiveChatAdapter(sender)
	binding := channelplatformentity.ChannelPlatform{
		Platform: platform.PlatformVKVideoLive,
		Enabled:  false,
	}

	// When
	err := Dispatch(
		context.Background(),
		newRegistry(adapter),
		[]channelplatformentity.ChannelPlatform{binding},
		[]platform.Platform{platform.PlatformVKVideoLive},
		"hello",
		"",
		ChatOptions{},
	)

	// Then
	require.NoError(t, err)
	require.Empty(t, sender.message)
}

func TestDispatchReportsVKVideoLiveSendError(t *testing.T) {
	// Given
	sendError := errors.New("VK Video Live send failed")
	sender := &fakeVKVideoLiveSender{err: sendError}
	adapter := NewVKVideoLiveChatAdapter(sender)
	binding := channelplatformentity.ChannelPlatform{
		Platform:          platform.PlatformVKVideoLive,
		PlatformChannelID: "vk-channel",
		Enabled:           true,
	}

	// When
	err := Dispatch(
		context.Background(),
		newRegistry(adapter),
		[]channelplatformentity.ChannelPlatform{binding},
		nil,
		"hello",
		"",
		ChatOptions{},
	)

	// Then
	require.ErrorIs(t, err, sendError)
	require.ErrorContains(t, err, `send chat message to "vk_video_live" binding "vk-channel"`)
}

type fakeVKVideoLiveSender struct {
	binding channelplatformentity.ChannelPlatform
	message string
	err     error
}

func (s *fakeVKVideoLiveSender) SendMessage(
	_ context.Context,
	binding channelplatformentity.ChannelPlatform,
	message string,
) error {
	s.binding = binding
	s.message = message
	return s.err
}
