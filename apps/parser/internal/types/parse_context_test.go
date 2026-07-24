package types

import (
	"testing"

	"github.com/google/uuid"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	"github.com/twirapp/twir/libs/entities/platform"
)

func newContextTestChannel() (channelentity.Channel, channelplatformentity.ChannelPlatform, channelplatformentity.ChannelPlatform) {
	twitchBinding := channelplatformentity.ChannelPlatform{
		ID:                uuid.New(),
		Platform:          platform.PlatformTwitch,
		PlatformChannelID: "twitch-channel",
		UserID:            uuid.New(),
	}
	kickBinding := channelplatformentity.ChannelPlatform{
		ID:                uuid.New(),
		Platform:          platform.PlatformKick,
		PlatformChannelID: "kick-channel",
		UserID:            uuid.New(),
	}
	channel := channelentity.Channel{
		ID:       uuid.New(),
		Bindings: []channelplatformentity.ChannelPlatform{twitchBinding, kickBinding},
	}

	return channel, twitchBinding, kickBinding
}

func TestNewParseContextChannelEmptyBindingIDSelectsRequestedPlatform(t *testing.T) {
	t.Parallel()

	channel, _, kickBinding := newContextTestChannel()

	got, ok := NewParseContextChannel(channel, platform.PlatformKick, "kick-name", "")
	if !ok {
		t.Fatal("NewParseContextChannel() ok = false, want true")
	}
	if got.ID != kickBinding.PlatformChannelID {
		t.Fatalf("ID = %q, want %q", got.ID, kickBinding.PlatformChannelID)
	}
	if got.Name != "kick-name" {
		t.Fatalf("Name = %q, want kick-name", got.Name)
	}
	if got.DBChannelID != channel.ID.String() {
		t.Fatalf("DBChannelID = %q, want %q", got.DBChannelID, channel.ID.String())
	}
}

func TestNewParseContextChannelExplicitBindingIDWinsRegardlessOfPlatform(t *testing.T) {
	t.Parallel()

	channel, twitchBinding, _ := newContextTestChannel()

	got, ok := NewParseContextChannel(channel, platform.PlatformKick, "name", twitchBinding.ID.String())
	if !ok {
		t.Fatal("NewParseContextChannel() ok = false, want true")
	}
	if got.ID != twitchBinding.PlatformChannelID {
		t.Fatalf("ID = %q, want explicit binding channel %q", got.ID, twitchBinding.PlatformChannelID)
	}
}

func TestNewParseContextChannelRejectsUnknownBindingID(t *testing.T) {
	t.Parallel()

	channel, _, _ := newContextTestChannel()

	got, ok := NewParseContextChannel(channel, platform.PlatformKick, "name", uuid.New().String())
	if ok {
		t.Fatal("NewParseContextChannel() ok = true, want false for unknown binding ID")
	}
	if got != nil {
		t.Fatalf("NewParseContextChannel() = %#v, want nil", got)
	}
}

func TestNewParseContextChannelRejectsMalformedBindingIDWithoutFallback(t *testing.T) {
	t.Parallel()

	channel, _, _ := newContextTestChannel()

	got, ok := NewParseContextChannel(channel, platform.PlatformKick, "name", "not-a-uuid")
	if ok {
		t.Fatal("NewParseContextChannel() ok = true, want false for malformed binding ID")
	}
	if got != nil {
		t.Fatalf("NewParseContextChannel() = %#v, want nil (no platform fallback)", got)
	}
}

func TestNewParseContextChannelPopulatesTwitchUserIDFromFirstTwitchBinding(t *testing.T) {
	t.Parallel()

	channel, twitchBinding, _ := newContextTestChannel()

	got, ok := NewParseContextChannel(channel, platform.PlatformKick, "name", "")
	if !ok {
		t.Fatal("NewParseContextChannel() ok = false, want true")
	}
	if got.TwitchUserID != twitchBinding.UserID {
		t.Fatalf("TwitchUserID = %v, want %v", got.TwitchUserID, twitchBinding.UserID)
	}
}

func TestNewParseContextChannelReturnsFalseWhenPlatformAbsent(t *testing.T) {
	t.Parallel()

	channel, _, _ := newContextTestChannel()

	got, ok := NewParseContextChannel(channel, platform.PlatformVKVideoLive, "name", "")
	if ok {
		t.Fatal("NewParseContextChannel() ok = true, want false for absent platform")
	}
	if got != nil {
		t.Fatalf("NewParseContextChannel() = %#v, want nil", got)
	}
}
