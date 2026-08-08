package entity

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
	CustomHTMLEnabled    bool
	CustomHTML           string
	CustomCSS            string
	CreatedAt            time.Time
	UpdatedAt            time.Time
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
