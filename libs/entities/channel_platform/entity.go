package channel_platform

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
)

type ChannelPlatform struct {
	ID                uuid.UUID
	ChannelID         uuid.UUID
	Platform          platform.Platform
	UserID            uuid.UUID
	PlatformChannelID string
	Enabled           bool
	BotUserID         *uuid.UUID
	BotConfig         json.RawMessage
	CreatedAt         time.Time
	UpdatedAt         time.Time

	isNil bool
}

func (c ChannelPlatform) IsNil() bool {
	return c.isNil
}

var Nil = ChannelPlatform{isNil: true}

type TwitchBotConfig struct {
	BotID          string `json:"bot_id"`
	IsBotMod       bool   `json:"is_bot_mod"`
	IsTwitchBanned bool   `json:"is_twitch_banned"`
}

func (c ChannelPlatform) ParseTwitchBotConfig() (TwitchBotConfig, error) {
	if len(c.BotConfig) == 0 {
		return TwitchBotConfig{}, nil
	}

	var config TwitchBotConfig
	if err := json.Unmarshal(c.BotConfig, &config); err != nil {
		return TwitchBotConfig{}, err
	}

	return config, nil
}
