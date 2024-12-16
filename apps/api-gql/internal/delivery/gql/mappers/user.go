package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func UserWithChannelToAdminUser(e entity.UserWithChannel) gqlmodel.TwirAdminUser {
	user := gqlmodel.TwirAdminUser{
		ID:         e.User.ID,
		IsBotAdmin: e.User.IsBotAdmin,
		IsBanned:   e.User.IsBanned,
		APIKey:     e.User.ApiKey,
	}

	if e.Channel != nil {
		user.IsBotEnabled = e.Channel.IsEnabled
		user.IsBotModerator = e.Channel.IsBotMod
	}

	return user
}
