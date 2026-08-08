package mappers

import (
	"strings"

	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	entity "github.com/twirapp/twir/libs/entities/overlays_stream_stats"
)

func MapStreamStatsCounterOrderToEntity(
	order []gqlmodel.StreamStatsOverlayCounter,
) []entity.StreamStatsOverlayCounter {
	result := make([]entity.StreamStatsOverlayCounter, 0, len(order))
	for _, counter := range order {
		result = append(result, entity.StreamStatsOverlayCounter(counter))
	}
	return result
}

func MapStreamStatsEntityToGQL(e entity.StreamStatsOverlay) gqlmodel.StreamStatsOverlay {
	counterOrder := make([]gqlmodel.StreamStatsOverlayCounter, 0, len(e.CounterOrder))
	for _, counter := range e.CounterOrder {
		counterOrder = append(counterOrder, gqlmodel.StreamStatsOverlayCounter(counter))
	}

	return gqlmodel.StreamStatsOverlay{
		ID:                   e.ID,
		ChannelID:            e.ChannelID,
		Design:               gqlmodel.StreamStatsOverlayDesign(e.Design),
		Variant:              gqlmodel.StreamStatsOverlayVariant(e.Variant),
		ViewersEnabled:       e.ViewersEnabled,
		ViewersMode:          gqlmodel.StreamStatsOverlayViewersMode(e.ViewersMode),
		PlatformIconsEnabled: e.PlatformIconsEnabled,
		MessagesEnabled:      e.MessagesEnabled,
		UptimeEnabled:        e.UptimeEnabled,
		SubscribersEnabled:   e.SubscribersEnabled,
		FollowersEnabled:     e.FollowersEnabled,
		ViewersColor:         e.ViewersColor,
		MessagesColor:        e.MessagesColor,
		UptimeColor:          e.UptimeColor,
		SubscribersColor:     e.SubscribersColor,
		FollowersColor:       e.FollowersColor,
		CounterOrder:         counterOrder,
		CustomHTMLEnabled:    e.CustomHTMLEnabled,
		CustomHTML:           e.CustomHTML,
		CustomCSS:            e.CustomCSS,
		CreatedAt:            e.CreatedAt,
		UpdatedAt:            e.UpdatedAt,
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
