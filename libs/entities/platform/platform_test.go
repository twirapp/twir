package platform

import (
	"slices"
	"testing"
)

func TestPlatformVKVideoLive(t *testing.T) {
	if got, want := PlatformVKVideoLive.String(), "vk_video_live"; got != want {
		t.Errorf("PlatformVKVideoLive.String() = %q, want %q", got, want)
	}
}

func TestAll(t *testing.T) {
	want := []Platform{
		PlatformTwitch,
		PlatformKick,
		PlatformVKVideoLive,
	}

	if got := All(); !slices.Equal(got, want) {
		t.Errorf("All() = %v, want %v", got, want)
	}
}

func TestPlatformIsValid(t *testing.T) {
	tests := []struct {
		name     string
		platform Platform
		want     bool
	}{
		{name: "twitch", platform: PlatformTwitch, want: true},
		{name: "kick", platform: PlatformKick, want: true},
		{name: "VK Video Live", platform: PlatformVKVideoLive, want: true},
		{name: "unknown", platform: "unknown", want: false},
		{name: "empty", platform: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.platform.IsValid(); got != tt.want {
				t.Errorf("Platform(%q).IsValid() = %t, want %t", tt.platform, got, tt.want)
			}
		})
	}
}

func TestPlatformSchema(t *testing.T) {
	schema := Platform("").Schema(nil)
	want := []any{
		string(PlatformTwitch),
		string(PlatformKick),
		string(PlatformVKVideoLive),
	}

	if !slices.Equal(schema.Enum, want) {
		t.Errorf("Platform.Schema().Enum = %v, want %v", schema.Enum, want)
	}
}

func TestCapabilityValues(t *testing.T) {
	tests := []struct {
		capability Capability
		want       string
	}{
		{capability: CapabilityChatRead, want: "chat.read"},
		{capability: CapabilityChatWrite, want: "chat.write"},
		{capability: CapabilityChatReply, want: "chat.reply"},
		{capability: CapabilityModerationDelete, want: "moderation.delete"},
		{capability: CapabilityStreamsRead, want: "streams.read"},
		{capability: CapabilityEventsFollow, want: "events.follow"},
		{capability: CapabilityEventsRaid, want: "events.raid"},
		{capability: CapabilityEventsReward, want: "events.reward"},
	}

	for _, tt := range tests {
		if got := string(tt.capability); got != tt.want {
			t.Errorf("Capability = %q, want %q", got, tt.want)
		}
	}
}

func TestCapabilitiesSupports(t *testing.T) {
	capabilities := Capabilities{
		CapabilityChatRead,
		CapabilityChatReply,
	}

	tests := []struct {
		capability Capability
		want       bool
	}{
		{capability: CapabilityChatRead, want: true},
		{capability: CapabilityChatReply, want: true},
		{capability: CapabilityChatWrite, want: false},
		{capability: CapabilityStreamsRead, want: false},
	}

	for _, tt := range tests {
		if got := capabilities.Supports(tt.capability); got != tt.want {
			t.Errorf("Capabilities.Supports(%q) = %t, want %t", tt.capability, got, tt.want)
		}
	}
}

func TestPlatformCapabilities(t *testing.T) {
	tests := []struct {
		platform Platform
		want     Capabilities
	}{
		{
			platform: PlatformTwitch,
			want: Capabilities{
				CapabilityChatRead,
				CapabilityChatWrite,
				CapabilityChatReply,
				CapabilityModerationDelete,
				CapabilityStreamsRead,
				CapabilityEventsFollow,
				CapabilityEventsRaid,
				CapabilityEventsReward,
			},
		},
		{
			platform: PlatformKick,
			want: Capabilities{
				CapabilityChatWrite,
				CapabilityChatReply,
				CapabilityStreamsRead,
				CapabilityEventsFollow,
				CapabilityEventsReward,
			},
		},
		{platform: PlatformVKVideoLive, want: Capabilities{}},
		{platform: "unknown", want: Capabilities{}},
	}

	for _, tt := range tests {
		if got := tt.platform.Capabilities(); !slices.Equal(got, tt.want) {
			t.Errorf("Platform(%q).Capabilities() = %v, want %v", tt.platform, got, tt.want)
		}
	}
}

func TestErrUnsupportedCapability(t *testing.T) {
	var err error = ErrUnsupportedCapability{
		Platform:   PlatformVKVideoLive,
		Capability: CapabilityEventsReward,
	}

	if got, want := err.Error(), `platform "vk_video_live" does not support capability "events.reward"`; got != want {
		t.Errorf("ErrUnsupportedCapability.Error() = %q, want %q", got, want)
	}
}
