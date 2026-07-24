package resolvers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

func mapChannelByAPIKeyResult(channel channelentity.Channel) *gqlmodel.ChannelByAPIKeyResult {
	result := &gqlmodel.ChannelByAPIKeyResult{ID: channel.ID}
	if twitchBinding, found := channel.Binding(platformentity.PlatformTwitch); found {
		userID := twitchBinding.UserID
		result.TwitchUserID = &userID
	}
	if kickBinding, found := channel.Binding(platformentity.PlatformKick); found {
		userID := kickBinding.UserID
		result.KickUserID = &userID
	}

	return result
}
