package match

import (
	"context"

	"github.com/google/uuid"
)

type ActionKind string

const (
	ActionCreate  ActionKind = "create"
	ActionResolve ActionKind = "resolve"
	ActionCancel  ActionKind = "cancel"
)

type LifecycleAction struct {
	Kind      ActionKind `json:"kind"`
	ChannelID uuid.UUID  `json:"channelId"`
	MatchID   int64      `json:"matchId"`
	Revision  uint64     `json:"revision"`
	Win       bool       `json:"win,omitempty"`
	HeroName  string     `json:"heroName,omitempty"`
}

type StateStore interface {
	Load(context.Context, uuid.UUID) (Snapshot, error)
	CompareAndSwap(context.Context, Snapshot, Snapshot, []LifecycleAction) (bool, error)
	UpdateStats(context.Context, uuid.UUID, int, int, int) error
}
