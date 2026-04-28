package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/twirapp/twir/libs/entities/platform"
)

type User struct {
	ID                string
	Platform          platform.Platform
	PlatformID        string
	TokenID           null.String
	IsBotAdmin        bool
	ApiKey            string
	IsBanned          bool
	HideOnLandingPage bool
	CreatedAt         time.Time
	Login             string
	DisplayName       string
	Avatar            string

	isNil bool
}

func (c User) IsNil() bool {
	return c.isNil
}

var Nil = User{
	isNil: true,
}

type OnlineUser struct {
	ID        uuid.UUID
	ChannelID string
	UserID    string
	UserName  string

	isNil bool
}

func (c OnlineUser) IsNil() bool {
	return c.isNil
}

var NilOnlineUser = OnlineUser{
	isNil: true,
}
