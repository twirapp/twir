package nats

import (
	"github.com/nats-io/nats.go"
	cfg "github.com/satont/twir/libs/config"
)

type Opts struct {
	Config cfg.Config
}

func New(opts Opts) (*nats.Conn, error) {
	return nats.Connect(opts.Config.NatsUrl, nats.Name("parser"))
}
