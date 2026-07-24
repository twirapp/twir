package channel

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/google/uuid"
	channelplatform "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

func TestChannelBinding(t *testing.T) {
	t.Parallel()

	kickBinding := channelplatform.ChannelPlatform{
		ID:                uuid.New(),
		Platform:          platformentity.PlatformKick,
		PlatformChannelID: "kick-channel",
	}
	twitchBinding := channelplatform.ChannelPlatform{
		ID:                uuid.New(),
		Platform:          platformentity.PlatformTwitch,
		PlatformChannelID: "twitch-channel",
	}

	t.Run("returns first matching platform binding", func(t *testing.T) {
		t.Parallel()

		channel := Channel{Bindings: []channelplatform.ChannelPlatform{kickBinding, twitchBinding}}

		got, found := channel.Binding(platformentity.PlatformTwitch)
		if !found {
			t.Fatal("Binding() found = false, want true")
		}
		if !reflect.DeepEqual(got, twitchBinding) {
			t.Fatalf("Binding() = %#v, want %#v", got, twitchBinding)
		}
	})

	t.Run("returns zero and false when absent", func(t *testing.T) {
		t.Parallel()

		channel := Channel{Bindings: []channelplatform.ChannelPlatform{kickBinding}}

		got, found := channel.Binding(platformentity.PlatformVKVideoLive)
		if found {
			t.Fatal("Binding() found = true, want false")
		}
		if !reflect.DeepEqual(got, channelplatform.ChannelPlatform{}) {
			t.Fatalf("Binding() = %#v, want zero value", got)
		}
	})

	t.Run("preserves first match with duplicated platforms", func(t *testing.T) {
		t.Parallel()

		second := channelplatform.ChannelPlatform{
			ID:                uuid.New(),
			Platform:          platformentity.PlatformTwitch,
			PlatformChannelID: "twitch-channel-2",
		}
		channel := Channel{Bindings: []channelplatform.ChannelPlatform{twitchBinding, second}}

		got, found := channel.Binding(platformentity.PlatformTwitch)
		if !found {
			t.Fatal("Binding() found = false, want true")
		}
		if !reflect.DeepEqual(got, twitchBinding) {
			t.Fatalf("Binding() = %#v, want first match %#v", got, twitchBinding)
		}
	})
}

func TestChannelBindingByID(t *testing.T) {
	t.Parallel()

	target := channelplatform.ChannelPlatform{
		ID:                uuid.New(),
		Platform:          platformentity.PlatformKick,
		PlatformChannelID: "kick-channel",
	}
	channel := Channel{Bindings: []channelplatform.ChannelPlatform{
		{ID: uuid.New(), Platform: platformentity.PlatformTwitch},
		target,
	}}

	t.Run("returns binding with matching id", func(t *testing.T) {
		t.Parallel()

		got, found := channel.BindingByID(target.ID)
		if !found {
			t.Fatal("BindingByID() found = false, want true")
		}
		if !reflect.DeepEqual(got, target) {
			t.Fatalf("BindingByID() = %#v, want %#v", got, target)
		}
	})

	t.Run("returns zero and false when absent", func(t *testing.T) {
		t.Parallel()

		got, found := channel.BindingByID(uuid.New())
		if found {
			t.Fatal("BindingByID() found = true, want false")
		}
		if !reflect.DeepEqual(got, channelplatform.ChannelPlatform{}) {
			t.Fatalf("BindingByID() = %#v, want zero value", got)
		}
	})
}

func TestChannelTwitchBinding(t *testing.T) {
	t.Parallel()

	twitchUserID := uuid.New()
	fullConfig := json.RawMessage(`{"bot_id":"twitch-bot","is_bot_mod":true,"is_twitch_banned":false}`)

	newTwitchBinding := func(botConfig json.RawMessage) channelplatform.ChannelPlatform {
		return channelplatform.ChannelPlatform{
			ID:                uuid.New(),
			Platform:          platformentity.PlatformTwitch,
			PlatformChannelID: "twitch-channel",
			UserID:            twitchUserID,
			Enabled:           true,
			BotConfig:         botConfig,
		}
	}

	t.Run("returns absent semantics when no twitch binding", func(t *testing.T) {
		t.Parallel()

		channel := Channel{Bindings: []channelplatform.ChannelPlatform{
			{ID: uuid.New(), Platform: platformentity.PlatformKick},
		}}

		binding, config, found, err := channel.TwitchBinding()
		if err != nil {
			t.Fatalf("TwitchBinding() error = %v, want nil", err)
		}
		if found {
			t.Fatal("TwitchBinding() found = true, want false")
		}
		if !reflect.DeepEqual(binding, channelplatform.ChannelPlatform{}) || !reflect.DeepEqual(config, channelplatform.TwitchBotConfig{}) {
			t.Fatalf("TwitchBinding() = (%#v, %#v), want zero values", binding, config)
		}
	})

	t.Run("selects twitch binding regardless of order and parses config", func(t *testing.T) {
		t.Parallel()

		channel := Channel{Bindings: []channelplatform.ChannelPlatform{
			{
				ID:                uuid.New(),
				Platform:          platformentity.PlatformKick,
				PlatformChannelID: "kick-channel",
				BotConfig:         json.RawMessage(`{"bot_id":"kick-bot"}`),
			},
			newTwitchBinding(fullConfig),
		}}

		binding, config, found, err := channel.TwitchBinding()
		if err != nil {
			t.Fatalf("TwitchBinding() error = %v, want nil", err)
		}
		if !found {
			t.Fatal("TwitchBinding() found = false, want true")
		}
		if binding.Platform != platformentity.PlatformTwitch || binding.PlatformChannelID != "twitch-channel" || binding.UserID != twitchUserID {
			t.Fatalf("TwitchBinding() binding = %#v", binding)
		}
		if config.BotID != "twitch-bot" || !config.IsBotMod || config.IsTwitchBanned {
			t.Fatalf("TwitchBinding() config = %#v", config)
		}
	})

	t.Run("accepts empty object and null configs", func(t *testing.T) {
		t.Parallel()

		for _, raw := range []json.RawMessage{nil, {}, json.RawMessage(`{}`), json.RawMessage(`null`)} {
			channel := Channel{Bindings: []channelplatform.ChannelPlatform{newTwitchBinding(raw)}}

			_, config, found, err := channel.TwitchBinding()
			if err != nil {
				t.Fatalf("TwitchBinding(%s) error = %v, want nil", raw, err)
			}
			if !found {
				t.Fatalf("TwitchBinding(%s) found = false, want true", raw)
			}
			if !reflect.DeepEqual(config, channelplatform.TwitchBotConfig{}) {
				t.Fatalf("TwitchBinding(%s) config = %#v, want zero value", raw, config)
			}
		}
	})

	t.Run("returns all zero values and error for malformed config", func(t *testing.T) {
		t.Parallel()

		channel := Channel{Bindings: []channelplatform.ChannelPlatform{
			newTwitchBinding(json.RawMessage(`{`)),
		}}

		binding, config, found, err := channel.TwitchBinding()
		if err == nil {
			t.Fatal("TwitchBinding() error = nil, want non-nil")
		}
		if found {
			t.Fatal("TwitchBinding() found = true, want false")
		}
		if !reflect.DeepEqual(binding, channelplatform.ChannelPlatform{}) || !reflect.DeepEqual(config, channelplatform.TwitchBotConfig{}) {
			t.Fatalf("TwitchBinding() = (%#v, %#v), want zero values", binding, config)
		}
	})
}

func TestChannelPlatforms(t *testing.T) {
	t.Parallel()

	channel := Channel{Bindings: []channelplatform.ChannelPlatform{
		{ID: uuid.New(), Platform: platformentity.PlatformKick},
		{ID: uuid.New(), Platform: platformentity.PlatformTwitch},
	}}

	got := channel.Platforms()
	want := []platformentity.Platform{platformentity.PlatformKick, platformentity.PlatformTwitch}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Platforms() = %#v, want %#v", got, want)
	}
}
