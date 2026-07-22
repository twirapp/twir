package timers

import (
	"testing"

	"github.com/google/uuid"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

func TestFindTwitchBindingSelectsTwitchBindingWhenKickComesFirst(t *testing.T) {
	t.Parallel()

	twitchOwnerID := uuid.New()
	channel := channelsmodel.Channel{
		Bindings: []channelplatformsmodel.ChannelPlatform{
			{
				Platform:          platformentity.PlatformKick,
				PlatformChannelID: "kick-channel-id",
				UserID:            uuid.New(),
				BotConfig:         []byte(`{"bot_id":"kick-bot"}`),
			},
			{
				Platform:          platformentity.PlatformTwitch,
				PlatformChannelID: "twitch-channel-id",
				UserID:            twitchOwnerID,
				BotConfig:         []byte(`{"bot_id":"twitch-bot"}`),
			},
		},
	}

	binding, ok := findTwitchBinding(channel)
	if !ok {
		t.Fatal("expected Twitch binding")
	}
	if binding.PlatformChannelID != "twitch-channel-id" {
		t.Fatalf("platform channel ID = %q, want Twitch binding ID", binding.PlatformChannelID)
	}
	if binding.UserID != twitchOwnerID {
		t.Fatalf("owner user ID = %s, want %s", binding.UserID, twitchOwnerID)
	}
	if string(binding.BotConfig) != `{"bot_id":"twitch-bot"}` {
		t.Fatalf("bot config = %s, want selected Twitch binding config", binding.BotConfig)
	}
}
