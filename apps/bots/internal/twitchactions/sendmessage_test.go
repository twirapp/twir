package twitchactions

import (
	"testing"

	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	"github.com/twirapp/twir/libs/entities/platform"
)

func TestActiveTwitchBindingUsesNormalizedState(t *testing.T) {
	t.Parallel()

	config, active, err := activeTwitchBinding(channelplatformentity.ChannelPlatform{
		Platform:          platform.PlatformTwitch,
		PlatformChannelID: "broadcaster-id",
		Enabled:           true,
		BotConfig:         []byte(`{"bot_id":"bot-id","is_bot_mod":true,"is_twitch_banned":false}`),
	})
	if err != nil {
		t.Fatalf("activeTwitchBinding() error = %v", err)
	}
	if !active {
		t.Fatal("expected enabled moderated binding to be active")
	}
	if config.BotID != "bot-id" {
		t.Errorf("expected bot ID %q, got %q", "bot-id", config.BotID)
	}
}

func TestActiveTwitchBindingSkipsDisabledBinding(t *testing.T) {
	t.Parallel()

	_, active, err := activeTwitchBinding(channelplatformentity.ChannelPlatform{
		Platform:          platform.PlatformTwitch,
		PlatformChannelID: "broadcaster-id",
		Enabled:           false,
		BotConfig:         []byte(`{"bot_id":"bot-id","is_bot_mod":true}`),
	})
	if err != nil {
		t.Fatalf("activeTwitchBinding() error = %v", err)
	}
	if active {
		t.Error("expected disabled binding to be skipped")
	}
}
