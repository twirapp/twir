package buscore

import (
	"time"

	"github.com/nats-io/nats.go"
	cfg "github.com/satont/twir/libs/config"
	botsservice "github.com/twirapp/twir/libs/bus-core/bots"
	emotes_cacher "github.com/twirapp/twir/libs/bus-core/emotes-cacher"
	"github.com/twirapp/twir/libs/bus-core/eval"
	"github.com/twirapp/twir/libs/bus-core/eventsub"
	"github.com/twirapp/twir/libs/bus-core/parser"
	"github.com/twirapp/twir/libs/bus-core/timers"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/bus-core/websockets"
)

type Bus struct {
	Parser       *parserBus
	Websocket    *websocketBus
	Channel      *channelBus
	Bots         *botsBus
	EmotesCacher *emotesCacherBus
	Timers       *timersBus
	Eval         *evalBus
	EventSub     *eventSubBus
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

		Bots: &botsBus{
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

		EmotesCacher: &emotesCacherBus{
			CacheGlobalEmotes: NewNatsQueue[struct{}, struct{}](
				nc,
				emotes_cacher.EMOTES_CACHER_GLOBAL_EMOTES_SUBJECT,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			CacheChannelEmotes: NewNatsQueue[emotes_cacher.EmotesCacheRequest, struct{}](
				nc,
				emotes_cacher.EMOTES_CACHER_CHANNEL_EMOTES_SUBJECT,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
		},

		Timers: &timersBus{
			AddTimer: NewNatsQueue[timers.AddOrRemoveTimerRequest, struct{}](
				nc,
				timers.AddTimerSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			RemoveTimer: NewNatsQueue[timers.AddOrRemoveTimerRequest, struct{}](
				nc,
				timers.RemoveTimerSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
		},

		Eval: &evalBus{
			Evaluate: NewNatsQueue[eval.EvalRequest, eval.EvalResponse](
				nc,
				eval.EvalEvaluateSubject,
				1*time.Minute,
				nats.JSON_ENCODER,
			),
		},

		EventSub: &eventSubBus{
			Subscribe: NewNatsQueue[eventsub.EventsubSubscribeRequest, struct{}](
				nc,
				eventsub.EventsubSubscribeSubject,
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
