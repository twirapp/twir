package recorder

import (
	"github.com/twirapp/twir/libs/audit"
	"go.uber.org/fx"
)

type FxFanoutOptions struct {
	fx.In

	Recorders []audit.Recorder `group:"audit-recorders"`
}

func NewFxFanout(options FxFanoutOptions) Fanout {
	return NewFanout(options.Recorders...)
}
