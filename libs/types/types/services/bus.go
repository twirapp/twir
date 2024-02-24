package services

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/satont/twir/libs/types/types/services/parser"
	"github.com/satont/twir/libs/types/types/services/twitch"
)

type Bus struct {
	ParserCommands Queue[twitch.TwitchChatMessage, parser.CommandParseResponse]
	BotsMessages   Queue[twitch.TwitchChatMessage, struct{}]
}

func NewNatsBus(nc *nats.Conn) *Bus {
	return &Bus{
		ParserCommands: NewNatsQueue[twitch.TwitchChatMessage, parser.CommandParseResponse](
			nc,
			parser.PARSER_COMMANDS_QUEUE,
			30*time.Minute,
		),
		BotsMessages: NewNatsQueue[twitch.TwitchChatMessage, struct{}](
			nc,
			twitch.CHAT_MESSAGE_BOTS_QUEUE,
			30*time.Minute,
		),
	}
}
