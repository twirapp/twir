package dota

import (
	"time"

	"github.com/google/uuid"
)

type PredictionSettings struct {
	Enabled       bool   `json:"enabled"`
	TitleTemplate string `json:"titleTemplate"`
	WindowSeconds int    `json:"windowSeconds"`
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

type CommandsSettings struct {
	Mmr bool `json:"mmr"`
	Wl  bool `json:"wl"`
	Lg  bool `json:"lg"`
	Gm  bool `json:"gm"`
	Np  bool `json:"np"`
	Wp  bool `json:"wp"`
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

func (c ChannelDotaSettings) IsNil() bool { return c.isNil }

var Nil = ChannelDotaSettings{isNil: true}

func (c ChannelDotaSettings) Winrate() float64 {
	total := c.SessionWins + c.SessionLosses
	if total == 0 {
		return 0
	}
	return float64(c.SessionWins) / float64(total) * 100
}
