package platform

import (
	"database/sql/driver"
	"fmt"
	"slices"

	"github.com/danielgtaylor/huma/v2"
)

type Platform string

type Capability string

type Capabilities []Capability

const (
	PlatformTwitch      Platform = "twitch"
	PlatformKick        Platform = "kick"
	PlatformVKVideoLive Platform = "vk_video_live"
)

const (
	CapabilityChatRead         Capability = "chat.read"
	CapabilityChatWrite        Capability = "chat.write"
	CapabilityChatReply        Capability = "chat.reply"
	CapabilityModerationDelete Capability = "moderation.delete"
	CapabilityStreamsRead      Capability = "streams.read"
	CapabilityEventsFollow     Capability = "events.follow"
	CapabilityEventsRaid       Capability = "events.raid"
	CapabilityEventsReward     Capability = "events.reward"
)

func (p Platform) IsValid() bool {
	switch p {
	case PlatformTwitch, PlatformKick, PlatformVKVideoLive:
		return true
	}
	return false
}

func (Platform) Schema(r huma.Registry) *huma.Schema {
	return &huma.Schema{
		Type: "string",
		Enum: []any{
			string(PlatformTwitch),
			string(PlatformKick),
			string(PlatformVKVideoLive),
		},
	}
}

func (p Platform) String() string { return string(p) }

func (c Capabilities) Supports(capability Capability) bool {
	return slices.Contains(c, capability)
}

type ErrUnsupportedCapability struct {
	Platform   Platform
	Capability Capability
}

func (e ErrUnsupportedCapability) Error() string {
	return fmt.Sprintf("platform %q does not support capability %q", e.Platform, e.Capability)
}

func (p *Platform) Scan(src any) error {
	switch v := src.(type) {
	case string:
		*p = Platform(v)
	case []byte:
		*p = Platform(v)
	case nil:
		*p = ""
	default:
		return fmt.Errorf("platform: cannot scan type %T into Platform", src)
	}
	return nil
}

func (p Platform) Value() (driver.Value, error) {
	return string(p), nil
}

func ShouldExecute(platforms []Platform, current Platform) bool {
	if len(platforms) == 0 {
		return true
	}

	return slices.Contains(platforms, current)
}

func All() []Platform {
	return []Platform{
		PlatformTwitch,
		PlatformKick,
		PlatformVKVideoLive,
	}
}
