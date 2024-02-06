package pubsub

import (
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/pubsub"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Config cfg.Config
}

func New(opts Opts) (*PubSub, error) {
	pb, err := pubsub.NewPubSub(opts.Config.RedisUrl)
	if err != nil {
		return nil, err
	}

	return &PubSub{
		Client: pb,
	}, nil
}

type PubSub struct {
	Client *pubsub.PubSub
}
