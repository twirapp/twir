package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func UserWithChannelToAdminUser(e entity.UserWithChannel) gqlmodel.TwirAdminUser {
	platform, err := EntityPlatformToGraphQL(e.User.Platform)
	if err != nil {
		platform = gqlmodel.PlatformTwitch
	}

	user := gqlmodel.TwirAdminUser{
		ID:          e.User.ID,
		Platform:    platform,
		PlatformID:  e.User.PlatformID,
		Login:       e.User.Login,
		DisplayName: e.User.DisplayName,
		Avatar:      &e.User.Avatar,
		IsBotAdmin:  e.User.IsBotAdmin,
		IsBanned:    e.User.IsBanned,
		APIKey:      e.User.ApiKey,
	}

	if e.User.Avatar == "" {
		user.Avatar = nil
	}

	if e.Channel != nil {
		user.IsBotEnabled = e.Channel.IsEnabled
		user.IsBotModerator = e.Channel.IsBotMod
	}

	return user
}
