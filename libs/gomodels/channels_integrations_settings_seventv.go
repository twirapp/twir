package model

import (
	"github.com/google/uuid"
	"github.com/guregu/null"
)

type ChannelsIntegrationsSettingsSeventv struct {
	ID                     uuid.UUID   `gorm:"column:id;type:uuid"`
	RewardIdForAddEmote    null.String `gorm:"column:reward_id_for_add_emote"`
	RewardIdForRemoveEmote null.String `gorm:"column:reward_id_for_remove_emote"`

	ChannelID string    `gorm:"column:channel_id;type:text" json:"channel_id"`
	Channel   *Channels `gorm:"foreignkey:ChannelID;" json:"channel"`
}

func (c ChannelsIntegrationsSettingsSeventv) TableName() string {
	return "channels_integrations_seventv"
}
