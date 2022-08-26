package nats

import (
	"tsuwari/parser/internal/config/cfg"

	"github.com/nats-io/nats.go"
)

var Nats *nats.Conn

func Connect() {
	cfg := cfg.Cfg

	if cfg.NatsUrl == nil {
		panic("Nats url not setuped.")
	}

	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		panic(err)
	}

	Nats = nc
}