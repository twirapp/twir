package workers

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

func TestCountEnabledTwitchChannelsUsesSelectedBinding(t *testing.T) {
	channels := []channelsmodel.Channel{
		{
			Bindings: []channelplatformsmodel.ChannelPlatform{
				{Platform: platform.PlatformKick, Enabled: true},
				{Platform: platform.PlatformTwitch, Enabled: false},
			},
		},
		{
			Bindings: []channelplatformsmodel.ChannelPlatform{
				{Platform: platform.PlatformKick, Enabled: false},
				{Platform: platform.PlatformTwitch, Enabled: true},
			},
		},
		{
			Bindings: []channelplatformsmodel.ChannelPlatform{
				{Platform: platform.PlatformKick, Enabled: true},
			},
		},
	}

	require.Equal(t, 1, countEnabledTwitchChannels(channels))
}
