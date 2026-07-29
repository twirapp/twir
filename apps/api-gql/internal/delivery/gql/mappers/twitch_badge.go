package mappers

import (
	"github.com/kvizyx/twitchy/helix"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
)

func TwitchBadgeTo(badge helix.ChatBadgeSet) gqlmodel.TwitchBadge {
	versions := make([]gqlmodel.TwitchBadgeVersion, 0, len(badge.Versions))
	for _, version := range badge.Versions {
		versions = append(
			versions, gqlmodel.TwitchBadgeVersion{
				ID:         version.ID,
				ImageURL1x: version.ImageURL1x,
				ImageURL2x: version.ImageURL2x,
				ImageURL4x: version.ImageURL4x,
			},
		)
	}

	return gqlmodel.TwitchBadge{
		SetID:    badge.SetID,
		Versions: versions,
	}
}
