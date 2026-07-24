package model

import (
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	usermodel "github.com/twirapp/twir/libs/repositories/users/model"
)

type UserWithChannel struct {
	User    usermodel.User
	Channel *channelentity.Channel
}

var Nil = UserWithChannel{}
