package mappers

import (
	"strings"

	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func MapStreamStatsEntityToGQL(e entity.StreamStatsOverlay) gqlmodel.StreamStatsOverlay {
	return gqlmodel.StreamStatsOverlay{
		ID:                 e.ID,
		ChannelID:          e.ChannelID,
		Design:             gqlmodel.StreamStatsOverlayDesign(e.Design),
		ViewersEnabled:     e.ViewersEnabled,
		ViewersMode:        gqlmodel.StreamStatsOverlayViewersMode(e.ViewersMode),
		MessagesEnabled:    e.MessagesEnabled,
		UptimeEnabled:      e.UptimeEnabled,
		SubscribersEnabled: e.SubscribersEnabled,
		FollowersEnabled:   e.FollowersEnabled,
		CustomHTMLEnabled:  e.CustomHTMLEnabled,
		CustomHTML:         e.CustomHTML,
		CustomCSS:          e.CustomCSS,
		CreatedAt:          e.CreatedAt,
		UpdatedAt:          e.UpdatedAt,
	}
}

func MapStreamStatsCountersEntityToGQL(e entity.StreamStatsOverlayCounters) gqlmodel.StreamStatsOverlayCounters {
	platformViewers := make([]gqlmodel.StreamStatsOverlayPlatformViewers, 0, len(e.PlatformViewers))
	for _, platformViewer := range e.PlatformViewers {
		platformViewers = append(platformViewers, gqlmodel.StreamStatsOverlayPlatformViewers{
			Platform: gqlmodel.Platform(strings.ToUpper(string(platformViewer.Platform))),
			Viewers:  platformViewer.Viewers,
		})
	}

	return gqlmodel.StreamStatsOverlayCounters{
		Live:            e.Live,
		Viewers:         e.Viewers,
		PlatformViewers: platformViewers,
		Messages:        e.Messages,
		StartedAt:       e.StartedAt,
		Subscribers:     e.Subscribers,
		Followers:       e.Followers,
	}
}
