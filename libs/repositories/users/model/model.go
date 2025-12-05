package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
)

type User struct {
	ID                string
	TokenID           null.String
	IsBotAdmin        bool
	ApiKey            string
	IsBanned          bool
	HideOnLandingPage bool
	CreatedAt         time.Time

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
