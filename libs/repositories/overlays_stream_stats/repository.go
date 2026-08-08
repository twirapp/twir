package overlays_stream_stats

import (
	"context"
	"errors"

	"github.com/google/uuid"
	entity "github.com/twirapp/twir/libs/entities/overlays_stream_stats"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID uuid.UUID) (entity.StreamStatsOverlay, error)
	Create(ctx context.Context, input CreateInput) (entity.StreamStatsOverlay, error)
	Update(ctx context.Context, channelID uuid.UUID, input UpdateInput) (entity.StreamStatsOverlay, error)
}

type CreateInput struct {
	ChannelID            uuid.UUID
	Design               string
	Variant              string
	ViewersEnabled       bool
	ViewersMode          string
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
	CounterOrder         []string
	CustomHTMLEnabled    bool
	CustomHTML           string
	CustomCSS            string
}

type UpdateInput struct {
	Design               string
	Variant              string
	ViewersEnabled       bool
	ViewersMode          string
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
	CounterOrder         []string
	CustomHTMLEnabled    bool
	CustomHTML           string
	CustomCSS            string
}

var ErrNotFound = errors.New("not found")
