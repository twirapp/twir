package model

import (
	"database/sql/driver"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

type PredictionSettings struct {
	Enabled       bool   `json:"enabled"`
	TitleTemplate string `json:"titleTemplate"`
	WindowSeconds int    `json:"windowSeconds"`
}

func (s PredictionSettings) Value() (driver.Value, error) {
	d, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return string(d), nil
}

func (s *PredictionSettings) Scan(src interface{}) error {
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

	return json.Unmarshal(source, s)
}

type ChatEventSettings struct {
	Enabled  bool   `json:"enabled"`
	Template string `json:"template"`
	Cooldown int    `json:"cooldown"`
}

type ChatEvents struct {
	MatchStarted ChatEventSettings `json:"matchStarted"`
	MatchEnded   ChatEventSettings `json:"matchEnded"`
	RoshanKilled ChatEventSettings `json:"roshanKilled"`
	AegisPickup  ChatEventSettings `json:"aegisPickup"`
}

func (s ChatEvents) Value() (driver.Value, error) {
	d, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return string(d), nil
}

func (s *ChatEvents) Scan(src interface{}) error {
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

	return json.Unmarshal(source, s)
}

type CommandsSettings struct {
	Mmr bool `json:"mmr"`
	Wl  bool `json:"wl"`
	Lg  bool `json:"lg"`
	Gm  bool `json:"gm"`
	Np  bool `json:"np"`
	Wp  bool `json:"wp"`
}

func (s CommandsSettings) Value() (driver.Value, error) {
	d, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return string(d), nil
}

func (s *CommandsSettings) Scan(src interface{}) error {
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

	return json.Unmarshal(source, s)
}

type ChannelDotaSettings struct {
	ID                 uuid.UUID
	ChannelID          uuid.UUID
	Enabled            bool
	SteamAccountID     *string
	GsiToken           string
	Mmr                int
	MmrDelta           int
	SessionWins        int
	SessionLosses      int
	PredictionSettings PredictionSettings
	ChatEvents         ChatEvents
	CommandsSettings   CommandsSettings
	CreatedAt          time.Time
	UpdatedAt          time.Time

	isNil bool
}

func (c ChannelDotaSettings) IsNil() bool {
	return c.isNil
}

var Nil = ChannelDotaSettings{isNil: true}
