package model

import (
	channelmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	usermodel "github.com/twirapp/twir/libs/repositories/users/model"
)

type UserWithChannel struct {
	User    usermodel.User
	Channel *channelmodel.Channel
}

var Nil = UserWithChannel{}
