package buscore

import (
	"time"

	"github.com/nats-io/nats.go"
	cfg "github.com/satont/twir/libs/config"
	botsservice "github.com/twirapp/twir/libs/bus-core/bots"
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

type bots struct {
	ProcessMessage Queue[twitch.TwitchChatMessage, struct{}]
	SendMessage    Queue[botsservice.SendMessageRequest, struct{}]
	DeleteMessage  Queue[botsservice.DeleteMessageRequest, struct{}]
}

type Bus struct {
	Parser    *parserBus
	Websocket *websocketBus
	Channel   *channelBus
	Bots      *bots
}

func NewNatsBus(nc *nats.Conn) *Bus {
	return &Bus{
		Parser: &parserBus{
			GetCommandResponse: NewNatsQueue[twitch.TwitchChatMessage, parser.CommandParseResponse](
				nc,
				PARSER_COMMANDS_SUBJECT,
				30*time.Minute,
				nats.GOB_ENCODER,
			),

			ParseVariablesInText: NewNatsQueue[parser.ParseVariablesInTextRequest, parser.ParseVariablesInTextResponse](
				nc,
				PARSER_TEXT_VARIABLES_SUBJECT,
				1*time.Minute,
				nats.GOB_ENCODER,
			),

			ProcessMessageAsCommand: NewNatsQueue[twitch.TwitchChatMessage, struct{}](
				nc,
				PARSER_PROCESS_MESSAGE_AS_COMMAND_SUBJECT,
				30*time.Minute,
				nats.GOB_ENCODER,
			),
		},

		Bots: &bots{
			ProcessMessage: NewNatsQueue[twitch.TwitchChatMessage, struct{}](
				nc,
				CHAT_MESSAGE_BOTS_SUBJECT,
				30*time.Minute,
				nats.GOB_ENCODER,
			),
			SendMessage: NewNatsQueue[botsservice.SendMessageRequest, struct{}](
				nc,
				botsservice.SendMessageSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			DeleteMessage: NewNatsQueue[botsservice.DeleteMessageRequest, struct{}](
				nc,
				botsservice.DeleteMessageSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
		},

		Websocket: &websocketBus{
			DudesGrow: NewNatsQueue[websockets.DudesGrowRequest, struct{}](
				nc,
				WEBSOCKETS_DUDES_GROW_SUBJECT,
				1*time.Minute,
				nats.GOB_ENCODER,
			),

			DudesUserSettings: NewNatsQueue[websockets.DudesChangeUserSettingsRequest, struct{}](
				nc,
				WEBSOCKETS_DUDES_CHANGE_COLOR_SUBJECT,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
		},

		Channel: &channelBus{
			StreamOnline: NewNatsQueue[twitch.StreamOnlineMessage, struct{}](
				nc,
				STREAM_ONLINE_SUBJECT,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			StreamOffline: NewNatsQueue[twitch.StreamOfflineMessage, struct{}](
				nc,
				STREAM_OFFLINE_SUBJECT,
				1*time.Minute,
				nats.GOB_ENCODER,
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
