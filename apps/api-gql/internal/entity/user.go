package entity

import (
	"time"

	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

type User struct {
	ID                string
	Platform          platformentity.Platform
	PlatformID        string
	TokenID           *string
	IsBotAdmin        bool
	ApiKey            string
	IsBanned          bool
	HideOnLandingPage bool
	Login             string
	DisplayName       string
	Avatar            string
}

var UserNil = User{}

type TwitchUser struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"displayName"`
	ProfileImageURL string `json:"profileImageUrl"`
	Description     string `json:"description"`
	NotFound        bool   `json:"notFound"`
}

type ChannelUserInfo struct {
	ID                string
	Messages          int
	Watched           int
	UsedEmotes        int
	UsedChannelPoints int
	IsMod             bool
	IsVip             bool
	IsSubscriber      bool
	FollowerSince     *time.Time
}
