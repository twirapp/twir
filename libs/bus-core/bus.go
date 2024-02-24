package buscore

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/twirapp/twir/libs/bus-core/parser"
	"github.com/twirapp/twir/libs/bus-core/twitch"
)

type Bus struct {
	ParserGetCommandResponse      Queue[twitch.TwitchChatMessage, parser.CommandParseResponse]
	ParserProcessMessageAsCommand Queue[twitch.TwitchChatMessage, struct{}]
	ParserParseVariablesInText    Queue[parser.ParseVariablesInTextRequest, parser.ParseVariablesInTextResponse]
	BotsMessages                  Queue[twitch.TwitchChatMessage, struct{}]
}

const parserQueue = "parser"
const botsQueue = "bots"

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
	}
}
