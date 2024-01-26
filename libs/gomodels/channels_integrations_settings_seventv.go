package model

import (
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/lib/pq"
)

type ChannelsIntegrationsSettingsSeventv struct {
	ID                         uuid.UUID      `gorm:"column:id;type:uuid"`
	RewardIdForAddEmote        null.String    `gorm:"column:reward_id_for_add_emote"`
	RewardIdForRemoveEmote     null.String    `gorm:"column:reward_id_for_remove_emote"`
	DeleteEmotesOnlyAddedByApp bool           `gorm:"column:delete_emotes_only_added_by_app"`
	AddedEmotes                pq.StringArray `gorm:"column:added_emotes;type:text[]"`

	ChannelID string    `gorm:"column:channel_id;type:text" json:"channel_id"`
	Channel   *Channels `gorm:"foreignkey:ChannelID;" json:"channel"`
}

func (c ChannelsIntegrationsSettingsSeventv) TableName() string {
	return "channels_integrations_seventv"
}
