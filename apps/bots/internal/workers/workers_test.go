package workers

import (
	"testing"

	"github.com/stretchr/testify/require"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	"github.com/twirapp/twir/libs/entities/platform"
)

func TestCountEnabledTwitchChannelsUsesSelectedBinding(t *testing.T) {
	channels := []channelentity.Channel{
		{
			Bindings: []channelplatformentity.ChannelPlatform{
				{Platform: platform.PlatformKick, Enabled: true},
				{Platform: platform.PlatformTwitch, Enabled: false},
			},
		},
		{
			Bindings: []channelplatformentity.ChannelPlatform{
				{Platform: platform.PlatformKick, Enabled: false},
				{Platform: platform.PlatformTwitch, Enabled: true},
			},
		},
		{
			Bindings: []channelplatformentity.ChannelPlatform{
				{Platform: platform.PlatformKick, Enabled: true},
			},
		},
	}

	require.Equal(t, 1, countEnabledTwitchChannels(channels))
}
