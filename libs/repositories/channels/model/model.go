package model

import (
	"github.com/google/uuid"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

type Channel struct {
	ID               uuid.UUID
	TwitchUserID     *uuid.UUID
	TwitchPlatformID *string `db:"twitch_platform_id"`
	TwitchBotEnabled bool    `db:"twitch_bot_enabled"`
	KickUserID       *uuid.UUID
	KickPlatformID   *string `db:"kick_platform_id"`
	KickBotEnabled   bool    `db:"kick_bot_enabled"`
	IsEnabled        bool
	IsTwitchBanned   bool
	IsBotMod         bool
	BotID            string
	KickBotID        *uuid.UUID
	ApiKey           *string `db:"api_key"`

	TwitchUser *usersmodel.User `json:"twitch_user,omitempty"`
	KickUser   *usersmodel.User `json:"kick_user,omitempty"`

	isNil bool
}

func (c Channel) IsNil() bool {
	return c.isNil
}

func (c Channel) TwitchConnected() bool {
	return c.TwitchUserID != nil && c.TwitchPlatformID != nil
}

func (c Channel) TwitchBotJoined() bool {
	return c.TwitchConnected() && c.TwitchBotEnabled
}

func (c Channel) KickConnected() bool {
	return c.KickUserID != nil && c.KickPlatformID != nil
}

func (c Channel) KickBotJoined() bool {
	return c.KickConnected() && c.KickBotEnabled
}

func (c Channel) AnyBotJoined() bool {
	return c.TwitchBotJoined() || c.KickBotJoined()
}

func (c Channel) Platforms() []platformentity.Platform {
	platforms := make([]platformentity.Platform, 0, 2)

	if c.TwitchConnected() {
		platforms = append(platforms, platformentity.PlatformTwitch)
	}

	if c.KickConnected() {
		platforms = append(platforms, platformentity.PlatformKick)
	}

	return platforms
}

var Nil = Channel{
	isNil: true,
}
