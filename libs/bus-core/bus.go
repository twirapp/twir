package buscore

import (
	"time"

	"github.com/nats-io/nats.go"
	cfg "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/libs/bus-core/parser"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/bus-core/websockets"
)

type parserBus struct {
	GetCommandResponse      Queue[twitch.TwitchChatMessage, parser.CommandParseResponse]
	ProcessMessageAsCommand Queue[twitch.TwitchChatMessage, struct{}]
	ParseVariablesInText    Queue[parser.ParseVariablesInTextRequest, parser.ParseVariablesInTextResponse]
}

type websocketBus struct {
	DudesGrow         Queue[websockets.DudesGrowRequest, struct{}]
	DudesUserSettings Queue[websockets.DudesChangeUserSettingsRequest, struct{}]
}

type channelBus struct {
	StreamOnline  Queue[twitch.StreamOnlineMessage, struct{}]
	StreamOffline Queue[twitch.StreamOfflineMessage, struct{}]
}

type Bus struct {
	Parser    *parserBus
	Websocket *websocketBus
	Channel   *channelBus

	BotsMessages Queue[twitch.TwitchChatMessage, struct{}]
}

func NewNatsBus(nc *nats.Conn) *Bus {
	return &Bus{
		Parser: &parserBus{
			GetCommandResponse: NewNatsQueue[twitch.TwitchChatMessage, parser.CommandParseResponse](
				nc,
				PARSER_COMMANDS_SUBJECT,
				30*time.Minute,
			),

			ParseVariablesInText: NewNatsQueue[parser.ParseVariablesInTextRequest, parser.ParseVariablesInTextResponse](
				nc,
				PARSER_TEXT_VARIABLES_SUBJECT,
				1*time.Minute,
			),

			ProcessMessageAsCommand: NewNatsQueue[twitch.TwitchChatMessage, struct{}](
				nc,
				PARSER_PROCESS_MESSAGE_AS_COMMAND_SUBJECT,
				30*time.Minute,
			),
		},

		BotsMessages: NewNatsQueue[twitch.TwitchChatMessage, struct{}](
			nc,
			CHAT_MESSAGE_BOTS_SUBJECT,
			30*time.Minute,
		),

		Websocket: &websocketBus{
			DudesGrow: NewNatsQueue[websockets.DudesGrowRequest, struct{}](
				nc,
				WEBSOCKETS_DUDES_GROW_SUBJECT,
				1*time.Minute,
			),

			DudesUserSettings: NewNatsQueue[websockets.DudesChangeUserSettingsRequest, struct{}](
				nc,
				WEBSOCKETS_DUDES_CHANGE_COLOR_SUBJECT,
				1*time.Minute,
			),
		},

		Channel: &channelBus{
			StreamOnline: NewNatsQueue[twitch.StreamOnlineMessage, struct{}](
				nc,
				STREAM_ONLINE_SUBJECT,
				1*time.Minute,
			),
			StreamOffline: NewNatsQueue[twitch.StreamOfflineMessage, struct{}](
				nc,
				STREAM_OFFLINE_SUBJECT,
				1*time.Minute,
			),
		},
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
