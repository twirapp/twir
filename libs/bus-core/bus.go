package buscore

import (
	"time"

	"github.com/nats-io/nats.go"
	cfg "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/libs/bus-core/parser"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/bus-core/websockets"
)

type Bus struct {
	ParserGetCommandResponse      Queue[twitch.TwitchChatMessage, parser.CommandParseResponse]
	ParserProcessMessageAsCommand Queue[twitch.TwitchChatMessage, struct{}]
	ParserParseVariablesInText    Queue[parser.ParseVariablesInTextRequest, parser.ParseVariablesInTextResponse]

	BotsMessages Queue[twitch.TwitchChatMessage, struct{}]

	WebsocketsDudesGrow        Queue[websockets.DudesGrowRequest, struct{}]
	WebsocketsDudesChangeColor Queue[websockets.DudesChangeColorRequest, struct{}]
}

const parserQueue = "parser"
const botsQueue = "bots"
const websocketsQueue = "websockets"

func NewNatsBus(nc *nats.Conn) *Bus {
	return &Bus{
		ParserGetCommandResponse: NewNatsQueue[twitch.TwitchChatMessage, parser.CommandParseResponse](
			nc,
			PARSER_COMMANDS_SUBJECT,
			parserQueue,
			30*time.Minute,
		),

		ParserParseVariablesInText: NewNatsQueue[parser.ParseVariablesInTextRequest, parser.ParseVariablesInTextResponse](
			nc,
			PARSER_TEXT_VARIABLES_SUBJECT,
			parserQueue,
			1*time.Minute,
		),

		ParserProcessMessageAsCommand: NewNatsQueue[twitch.TwitchChatMessage, struct{}](
			nc,
			PARSER_PROCESS_MESSAGE_AS_COMMAND_SUBJECT,
			parserQueue,
			30*time.Minute,
		),

		BotsMessages: NewNatsQueue[twitch.TwitchChatMessage, struct{}](
			nc,
			CHAT_MESSAGE_BOTS_SUBJECT,
			botsQueue,
			30*time.Minute,
		),

		WebsocketsDudesGrow: NewNatsQueue[websockets.DudesGrowRequest, struct{}](
			nc,
			WEBSOCKETS_DUDES_GROW_SUBJECT,
			websocketsQueue,
			1*time.Minute,
		),

		WebsocketsDudesChangeColor: NewNatsQueue[websockets.DudesChangeColorRequest, struct{}](
			nc,
			WEBSOCKETS_DUDES_CHANGE_COLOR_SUBJECT,
			websocketsQueue,
			1*time.Minute,
		),
	}
}

func NewNatsBusFx(serviceName string) func(config cfg.Config) (*Bus, error) {
	return func(config cfg.Config) (*Bus, error) {
		nc, err := nats.Connect(config.NatsUrl, nats.Name(serviceName))
		if err != nil {
			return nil, err
		}

		return NewNatsBus(nc), nil
	}
}