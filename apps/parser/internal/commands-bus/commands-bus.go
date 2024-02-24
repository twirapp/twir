package commands_bus

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/satont/twir/libs/types/types/services"
	service_parser "github.com/satont/twir/libs/types/types/services/parser"
	"github.com/satont/twir/libs/types/types/services/twitch"
)

type CommandsBus struct {
	bus services.Queue[twitch.TwitchChatMessage, service_parser.CommandParseResponse]
}

func New(nc *nats.Conn) *CommandsBus {
	b := &CommandsBus{
		bus: services.NewNatsQueue[twitch.TwitchChatMessage, service_parser.CommandParseResponse](
			nc,
			service_parser.PARSER_COMMANDS_QUEUE,
			30*time.Minute,
		),
	}

	return b
}

func (c *CommandsBus) Subscribe() error {
	return c.bus.Subscribe(
		func(ctx context.Context, data twitch.TwitchChatMessage) service_parser.CommandParseResponse {
			return service_parser.CommandParseResponse{}
		},
	)
}

func (c *CommandsBus) Unsubscribe() {
	c.bus.Unsubscribe()
}
