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
}

var Nil = User{}

type OnlineUser struct {
	ID        uuid.UUID
	ChannelID string
	UserID    string
	UserName  string
}

var NilOnlineUser = OnlineUser{}
