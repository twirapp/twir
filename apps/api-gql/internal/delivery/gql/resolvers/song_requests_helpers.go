package resolvers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/channelbinding"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

func mapChannelByAPIKeyResult(channel channelsmodel.Channel) *gqlmodel.ChannelByAPIKeyResult {
	result := &gqlmodel.ChannelByAPIKeyResult{ID: channel.ID}
	if twitchBinding, found := channelbinding.Find(channel, platformentity.PlatformTwitch); found {
		userID := twitchBinding.UserID
		result.TwitchUserID = &userID
	}
	if kickBinding, found := channelbinding.Find(channel, platformentity.PlatformKick); found {
		userID := kickBinding.UserID
		result.KickUserID = &userID
	}

	return result
}
