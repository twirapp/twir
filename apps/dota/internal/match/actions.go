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
	Kind       ActionKind `json:"kind"`
	ChannelID  uuid.UUID  `json:"channelId"`
	MatchID    int64      `json:"matchId"`
	Revision   uint64     `json:"revision"`
	MutationID string     `json:"mutationId"`
	Win        bool       `json:"win,omitempty"`
	HeroName   string     `json:"heroName,omitempty"`
}

func ActionMatchesSnapshot(action LifecycleAction, snapshot Snapshot) bool {
	return action.ChannelID != uuid.Nil &&
		action.ChannelID == snapshot.ChannelID &&
		action.Revision == snapshot.Revision &&
		action.MutationID != "" &&
		action.MutationID == snapshot.MutationID
}

type StateStore interface {
	Load(context.Context, uuid.UUID) (Snapshot, error)
	CompareAndSwap(context.Context, Snapshot, Snapshot, []LifecycleAction) (bool, error)
	UpdateStats(context.Context, uuid.UUID, int, int, int) error
}
