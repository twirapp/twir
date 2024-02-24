package nats

import (
	"context"

	"github.com/nats-io/nats.go"
	cfg "github.com/satont/twir/libs/config"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Config cfg.Config
}

func New(opts Opts) (*nats.Conn, error) {
	nc, err := nats.Connect(opts.Config.NatsUrl, nats.Name("eventsub"))
	if err != nil {
		return nil, err
	}

	opts.LC.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				nc.Close()
				return nil
			},
		},
	)

	return nc, nil
}
