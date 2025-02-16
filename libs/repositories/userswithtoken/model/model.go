package model

import (
	tokenmodel "github.com/twirapp/twir/libs/repositories/tokens/model"
	usermodel "github.com/twirapp/twir/libs/repositories/users/model"
)

type UserWithToken struct {
	User  usermodel.User
	Token *tokenmodel.Token
}
