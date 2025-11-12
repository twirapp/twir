package model

import (
	usermodel "github.com/twirapp/twir/libs/repositories/users/model"
	userstatsmodel "github.com/twirapp/twir/libs/repositories/users_stats/model"
)

type UserWithStats struct {
	User  usermodel.User
	Stats *userstatsmodel.UserStat
}
