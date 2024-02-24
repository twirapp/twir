package nats

import (
	"github.com/nats-io/nats.go"
	cfg "github.com/satont/twir/libs/config"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Config cfg.Config
}

func New(opts Opts) (*nats.Conn, error) {
	return nats.Connect(opts.Config.NatsUrl, nats.Name("timers"))
}
