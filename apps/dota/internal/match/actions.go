package match

import "github.com/google/uuid"

type ActionKind string

const (
	ActionCreate  ActionKind = "create"
	ActionResolve ActionKind = "resolve"
	ActionCancel  ActionKind = "cancel"
)

type LifecycleAction struct {
	Kind           ActionKind `json:"kind"`
	ChannelID      uuid.UUID  `json:"channelId"`
	MatchID        int64      `json:"matchId"`
	SteamAccountID string     `json:"steamAccountId,omitempty"`
	Win            bool       `json:"win,omitempty"`
	HeroName       string     `json:"heroName,omitempty"`
	TeamKnown      bool       `json:"teamKnown"`
}
