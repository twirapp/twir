package overlays_stream_stats

import (
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
)

type StreamStatsOverlayDesign string

const (
	StreamStatsOverlayDesignGlass    StreamStatsOverlayDesign = "GLASS"
	StreamStatsOverlayDesignCards    StreamStatsOverlayDesign = "CARDS"
	StreamStatsOverlayDesignNeon     StreamStatsOverlayDesign = "NEON"
	StreamStatsOverlayDesignSolid    StreamStatsOverlayDesign = "SOLID"
	StreamStatsOverlayDesignMinimal  StreamStatsOverlayDesign = "MINIMAL"
	StreamStatsOverlayDesignTerminal StreamStatsOverlayDesign = "TERMINAL"
	StreamStatsOverlayDesignOutline  StreamStatsOverlayDesign = "OUTLINE"
)

type StreamStatsOverlayVariant string

const (
	StreamStatsOverlayVariantHorizontal        StreamStatsOverlayVariant = "HORIZONTAL"
	StreamStatsOverlayVariantHorizontalCompact StreamStatsOverlayVariant = "HORIZONTAL_COMPACT"
	StreamStatsOverlayVariantVertical          StreamStatsOverlayVariant = "VERTICAL"
	StreamStatsOverlayVariantVerticalCompact   StreamStatsOverlayVariant = "VERTICAL_COMPACT"
	StreamStatsOverlayVariantLarge             StreamStatsOverlayVariant = "LARGE"
)

type StreamStatsOverlayViewersMode string

const (
	StreamStatsOverlayViewersModeCumulative StreamStatsOverlayViewersMode = "CUMULATIVE"
	StreamStatsOverlayViewersModeSeparate   StreamStatsOverlayViewersMode = "SEPARATE"
)

type StreamStatsOverlayCounter string

const (
	StreamStatsOverlayCounterViewers     StreamStatsOverlayCounter = "VIEWERS"
	StreamStatsOverlayCounterMessages    StreamStatsOverlayCounter = "MESSAGES"
	StreamStatsOverlayCounterUptime      StreamStatsOverlayCounter = "UPTIME"
	StreamStatsOverlayCounterSubscribers StreamStatsOverlayCounter = "SUBSCRIBERS"
	StreamStatsOverlayCounterFollowers   StreamStatsOverlayCounter = "FOLLOWERS"
)

func (c StreamStatsOverlayCounter) IsValid() bool {
	switch c {
	case StreamStatsOverlayCounterViewers,
		StreamStatsOverlayCounterMessages,
		StreamStatsOverlayCounterUptime,
		StreamStatsOverlayCounterSubscribers,
		StreamStatsOverlayCounterFollowers:
		return true
	}
	return false
}

var StreamStatsOverlayCountersDefaultOrder = []StreamStatsOverlayCounter{
	StreamStatsOverlayCounterViewers,
	StreamStatsOverlayCounterMessages,
	StreamStatsOverlayCounterUptime,
	StreamStatsOverlayCounterSubscribers,
	StreamStatsOverlayCounterFollowers,
}

// NormalizeCounterOrder returns every counter exactly once: input order first, then missing ones.
func NormalizeCounterOrder(order []StreamStatsOverlayCounter) []StreamStatsOverlayCounter {
	seen := make(map[StreamStatsOverlayCounter]bool, len(StreamStatsOverlayCountersDefaultOrder))
	result := make([]StreamStatsOverlayCounter, 0, len(StreamStatsOverlayCountersDefaultOrder))

	for _, counter := range order {
		if !counter.IsValid() || seen[counter] {
			continue
		}
		seen[counter] = true
		result = append(result, counter)
	}

	for _, counter := range StreamStatsOverlayCountersDefaultOrder {
		if !seen[counter] {
			result = append(result, counter)
		}
	}

	return result
}

type StreamStatsOverlay struct {
	ID                   uuid.UUID
	ChannelID            string
	Design               StreamStatsOverlayDesign
	Variant              StreamStatsOverlayVariant
	ViewersEnabled       bool
	ViewersMode          StreamStatsOverlayViewersMode
	PlatformIconsEnabled bool
	MessagesEnabled      bool
	UptimeEnabled        bool
	SubscribersEnabled   bool
	FollowersEnabled     bool
	ViewersColor         string
	MessagesColor        string
	UptimeColor          string
	SubscribersColor     string
	FollowersColor       string
	CounterOrder         []StreamStatsOverlayCounter
	CustomHTMLEnabled    bool
	CustomHTML           string
	CustomCSS            string
	CreatedAt            time.Time
	UpdatedAt            time.Time

	isNil bool
}

func (s StreamStatsOverlay) IsNil() bool {
	return s.isNil || s.ID == uuid.Nil
}

var Nil = StreamStatsOverlay{
	isNil: true,
}

type StreamStatsOverlayPlatformViewers struct {
	Platform platform.Platform
	Viewers  int
}

type StreamStatsOverlayCounters struct {
	Live            bool
	Viewers         int
	PlatformViewers []StreamStatsOverlayPlatformViewers
	Messages        int
	StartedAt       *time.Time
	Subscribers     *int
	Followers       *int
}
