package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
)

type StreamStatsOverlayDesign string

const (
	StreamStatsOverlayDesignBar     StreamStatsOverlayDesign = "BAR"
	StreamStatsOverlayDesignCards   StreamStatsOverlayDesign = "CARDS"
	StreamStatsOverlayDesignMinimal StreamStatsOverlayDesign = "MINIMAL"
)

type StreamStatsOverlayViewersMode string

const (
	StreamStatsOverlayViewersModeCumulative StreamStatsOverlayViewersMode = "CUMULATIVE"
	StreamStatsOverlayViewersModeSeparate   StreamStatsOverlayViewersMode = "SEPARATE"
)

type StreamStatsOverlay struct {
	ID                 uuid.UUID
	ChannelID          string
	Design             StreamStatsOverlayDesign
	ViewersEnabled     bool
	ViewersMode        StreamStatsOverlayViewersMode
	MessagesEnabled    bool
	UptimeEnabled      bool
	SubscribersEnabled bool
	FollowersEnabled   bool
	CustomHTMLEnabled  bool
	CustomHTML         string
	CustomCSS          string
	CreatedAt          time.Time
	UpdatedAt          time.Time
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
