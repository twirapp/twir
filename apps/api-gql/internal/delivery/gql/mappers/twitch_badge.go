package mappers

import (
	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
)

func TwitchBadgeTo(badge helix.ChatBadge) gqlmodel.TwitchBadge {
	versions := make([]gqlmodel.TwitchBadgeVersion, 0, len(badge.Versions))
	for _, version := range badge.Versions {
		versions = append(
			versions, gqlmodel.TwitchBadgeVersion{
				ID:         version.ID,
				ImageURL1x: version.ImageUrl1x,
				ImageURL2x: version.ImageUrl2x,
				ImageURL4x: version.ImageUrl4x,
			},
		)
	}

	return gqlmodel.TwitchBadge{
		SetID:    badge.SetID,
		Versions: versions,
	}
}
