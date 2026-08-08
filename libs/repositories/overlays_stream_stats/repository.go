package overlays_stream_stats

import (
	"context"
	"errors"

	"github.com/twirapp/twir/libs/repositories/overlays_stream_stats/model"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (model.StreamStatsOverlay, error)
	Create(ctx context.Context, input CreateInput) (model.StreamStatsOverlay, error)
	Update(ctx context.Context, channelID string, input UpdateInput) (model.StreamStatsOverlay, error)
}

type CreateInput struct {
	ChannelID          string
	Design             string
	ViewersEnabled     bool
	ViewersMode        string
	MessagesEnabled    bool
	UptimeEnabled      bool
	SubscribersEnabled bool
	FollowersEnabled   bool
	CustomHTMLEnabled  bool
	CustomHTML         string
	CustomCSS          string
}

type UpdateInput struct {
	Design             string
	ViewersEnabled     bool
	ViewersMode        string
	MessagesEnabled    bool
	UptimeEnabled      bool
	SubscribersEnabled bool
	FollowersEnabled   bool
	CustomHTMLEnabled  bool
	CustomHTML         string
	CustomCSS          string
}

var ErrNotFound = errors.New("not found")
