package channelbinding

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

func TestFindTwitchSelectsTwitchBindingRegardlessOfOrder(t *testing.T) {
	twitchUserID := uuid.New()
	channel := channelsmodel.Channel{
		Bindings: []channelplatformsmodel.ChannelPlatform{
			{
				Platform:          platform.PlatformKick,
				PlatformChannelID: "kick-channel",
				UserID:            uuid.New(),
				Enabled:           true,
				BotConfig:         json.RawMessage(`{"bot_id":"kick-bot"}`),
			},
			{
				Platform:          platform.PlatformTwitch,
				PlatformChannelID: "twitch-channel",
				UserID:            twitchUserID,
				Enabled:           true,
				BotConfig: json.RawMessage(
					`{"bot_id":"twitch-bot","is_bot_mod":true,"is_twitch_banned":false}`,
				),
			},
		},
	}

	binding, config, found, err := FindTwitch(channel)

	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, platform.PlatformTwitch, binding.Platform)
	require.Equal(t, "twitch-channel", binding.PlatformChannelID)
	require.Equal(t, twitchUserID, binding.UserID)
	require.Equal(t, "twitch-bot", config.BotID)
	require.True(t, config.IsBotMod)
	require.False(t, config.IsTwitchBanned)
}

func TestFindTwitchReturnsConfigError(t *testing.T) {
	channel := channelsmodel.Channel{
		Bindings: []channelplatformsmodel.ChannelPlatform{
			{
				Platform:  platform.PlatformTwitch,
				BotConfig: json.RawMessage(`{`),
			},
		},
	}

	_, _, found, err := FindTwitch(channel)

	require.Error(t, err)
	require.False(t, found)
}

func TestFindReturnsRequestedPlatformBinding(t *testing.T) {
	kickBinding := channelplatformsmodel.ChannelPlatform{
		Platform:          platform.PlatformKick,
		PlatformChannelID: "kick-channel",
	}
	twitchBinding := channelplatformsmodel.ChannelPlatform{
		Platform:          platform.PlatformTwitch,
		PlatformChannelID: "twitch-channel",
	}
	channel := channelsmodel.Channel{
		Bindings: []channelplatformsmodel.ChannelPlatform{kickBinding, twitchBinding},
	}

	binding, found := Find(channel, platform.PlatformTwitch)

	require.True(t, found)
	require.Equal(t, twitchBinding, binding)
}
