package nats

import (
	"tsuwari/parser/internal/config/cfg"

	"github.com/nats-io/nats.go"
)

var Nats *nats.Conn

func Connect() {
	nc, err := nats.Connect(cfg.Cfg.NatsUrl)

	if err != nil {
		panic(err)
	}

	Nats = nc
}