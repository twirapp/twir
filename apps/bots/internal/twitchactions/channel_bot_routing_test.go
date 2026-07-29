package twitchactions

import (
	"context"
	"errors"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/stretchr/testify/require"
	"github.com/twirapp/twir/libs/twitch"
)

func TestTwitchActions_createChannelBotClient_uses_distinct_channel_routes(t *testing.T) {
	// Given
	var routes [][2]string
	client := &helix.Client{}
	actions := &TwitchActions{
		newChannelBotClient: func(_ context.Context, botID string, channelID string) (*helix.Client, error) {
			routes = append(routes, [2]string{botID, channelID})
			return client, nil
		},
	}

	// When
	firstClient, firstErr := actions.createChannelBotClient(context.Background(), "bot-one", "channel-one")
	secondClient, secondErr := actions.createChannelBotClient(context.Background(), "bot-two", "channel-two")

	// Then
	require.NoError(t, firstErr)
	require.NoError(t, secondErr)
	require.Same(t, client, firstClient)
	require.Same(t, client, secondClient)
	require.Equal(t, [][2]string{{"bot-one", "channel-one"}, {"bot-two", "channel-two"}}, routes)
}

func TestTwitchActions_createChannelBotClient_preserves_missing_channel_bot_error(t *testing.T) {
	// Given
	actions := &TwitchActions{
		newChannelBotClient: func(context.Context, string, string) (*helix.Client, error) {
			return nil, twitch.ErrChannelBotNotRegistered
		},
	}

	// When
	_, err := actions.createChannelBotClient(context.Background(), "bot", "channel")

	// Then
	require.True(t, errors.Is(err, twitch.ErrChannelBotNotRegistered))
}
