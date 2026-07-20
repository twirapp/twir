package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
)

type ChannelsDotaSettingsPrediction struct {
	Enabled       bool   `json:"enabled"`
	TitleTemplate string `json:"titleTemplate"`
	WindowSeconds int    `json:"windowSeconds"`
}

func (c ChannelsDotaSettingsPrediction) Value() (driver.Value, error) {
	d, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	return string(d), nil
}

func (c *ChannelsDotaSettingsPrediction) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	var source []byte
	switch v := src.(type) {
	case []byte:
		source = v
	case string:
		source = []byte(v)
	default:
		return nil
	}

	return json.Unmarshal(source, c)
}

type ChannelsDotaSettingsChatEvent struct {
	Enabled  bool   `json:"enabled"`
	Template string `json:"template"`
	Cooldown int    `json:"cooldown"`
}

type ChannelsDotaSettingsChatEvents struct {
	MatchStarted ChannelsDotaSettingsChatEvent `json:"matchStarted"`
	MatchEnded   ChannelsDotaSettingsChatEvent `json:"matchEnded"`
	RoshanKilled ChannelsDotaSettingsChatEvent `json:"roshanKilled"`
	AegisPickup  ChannelsDotaSettingsChatEvent `json:"aegisPickup"`
}

func (c ChannelsDotaSettingsChatEvents) Value() (driver.Value, error) {
	d, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	return string(d), nil
}

func (c *ChannelsDotaSettingsChatEvents) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	var source []byte
	switch v := src.(type) {
	case []byte:
		source = v
	case string:
		source = []byte(v)
	default:
		return nil
	}

	return json.Unmarshal(source, c)
}

type ChannelsDotaSettingsCommands struct {
	Mmr bool `json:"mmr"`
	Wl  bool `json:"wl"`
	Lg  bool `json:"lg"`
	Gm  bool `json:"gm"`
	Np  bool `json:"np"`
	Wp  bool `json:"wp"`
}

func (c ChannelsDotaSettingsCommands) Value() (driver.Value, error) {
	d, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	return string(d), nil
}

func (c *ChannelsDotaSettingsCommands) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	var source []byte
	switch v := src.(type) {
	case []byte:
		source = v
	case string:
		source = []byte(v)
	default:
		return nil
	}

	return json.Unmarshal(source, c)
}

type ChannelsDotaSettings struct {
	ID                 uuid.UUID                      `gorm:"column:id;type:uuid"                  json:"id"`
	ChannelID          uuid.UUID                      `gorm:"column:channel_id;type:uuid"          json:"channelId"`
	Enabled            bool                           `gorm:"column:enabled;type:bool"             json:"enabled"`
	SteamAccountID     null.String                    `gorm:"column:steam_account_id;type:text"    json:"steamAccountId"`
	GsiToken           string                         `gorm:"column:gsi_token;type:text"           json:"gsiToken"`
	Mmr                int                            `gorm:"column:mmr;type:int"                  json:"mmr"`
	MmrDelta           int                            `gorm:"column:mmr_delta;type:int"            json:"mmrDelta"`
	SessionWins        int                            `gorm:"column:session_wins;type:int"         json:"sessionWins"`
	SessionLosses      int                            `gorm:"column:session_losses;type:int"       json:"sessionLosses"`
	PredictionSettings ChannelsDotaSettingsPrediction `gorm:"column:prediction_settings;type:jsonb" json:"predictionSettings"`
	ChatEvents         ChannelsDotaSettingsChatEvents `gorm:"column:chat_events;type:jsonb"        json:"chatEvents"`
	CommandsSettings   ChannelsDotaSettingsCommands   `gorm:"column:commands_settings;type:jsonb"  json:"commandsSettings"`
	CreatedAt          time.Time                      `gorm:"column:created_at"                    json:"createdAt"`
	UpdatedAt          time.Time                      `gorm:"column:updated_at"                    json:"updatedAt"`
}

func (ChannelsDotaSettings) TableName() string {
	return "channels_dota_settings"
}
