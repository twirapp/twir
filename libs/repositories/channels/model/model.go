package model

import (
	"github.com/google/uuid"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

type Channel struct {
	ID             uuid.UUID
	TwitchUserID   *uuid.UUID
	KickUserID     *uuid.UUID
	IsEnabled      bool
	IsTwitchBanned bool
	IsBotMod       bool
	BotID          string

	isNil bool
}

func (c Channel) IsNil() bool {
	return c.isNil
}

func (c Channel) TwitchConnected() bool {
	return c.TwitchUserID != nil
}

func (c Channel) KickConnected() bool {
	return c.KickUserID != nil
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
