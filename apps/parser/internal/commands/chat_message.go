package commands

import (
	"github.com/twirapp/twir/libs/bus-core/generic"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	streamsmodel "github.com/twirapp/twir/libs/repositories/streams/model"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	usersstatsmodel "github.com/twirapp/twir/libs/repositories/users_stats/model"
)

type ChatMessageContext struct {
	generic.ChatMessage
	Channel       channelsmodel.Channel
	Stream        *streamsmodel.Stream
	User          usersmodel.User
	UserStats     *usersstatsmodel.UserStat
	CommandPrefix string
}
