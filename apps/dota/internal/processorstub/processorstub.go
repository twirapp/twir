package processorstub

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/dota/internal/gsi"
)

// TODO(task-6): replace with the real match processor state machine.

type Noop struct{}

var _ gsi.MatchProcessor = (*Noop)(nil)

func New() *Noop {
	return &Noop{}
}

func (n *Noop) Process(_ context.Context, _ uuid.UUID, _ gsi.Payload) error {
	return nil
}
