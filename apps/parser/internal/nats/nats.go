package nats

import (
	"github.com/nats-io/nats.go"
	cfg "github.com/twirapp/twir/libs/config"
)

func New(config cfg.Config) (*nats.Conn, error) {
	return nats.Connect(config.NatsUrl, nats.Name("parser"))
}
