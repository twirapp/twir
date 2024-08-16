package buscore

import (
	"time"

	"github.com/nats-io/nats.go"
	cfg "github.com/satont/twir/libs/config"
	botsservice "github.com/twirapp/twir/libs/bus-core/bots"
	chat_messages_store "github.com/twirapp/twir/libs/bus-core/chat-messages-store"
	emotes_cacher "github.com/twirapp/twir/libs/bus-core/emotes-cacher"
	"github.com/twirapp/twir/libs/bus-core/eval"
	"github.com/twirapp/twir/libs/bus-core/eventsub"
	"github.com/twirapp/twir/libs/bus-core/parser"
	"github.com/twirapp/twir/libs/bus-core/scheduler"
	"github.com/twirapp/twir/libs/bus-core/timers"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/bus-core/websockets"
)

type Bus struct {
	Parser            *parserBus
	Websocket         *websocketBus
	Channel           *channelBus
	Bots              *botsBus
	EmotesCacher      *emotesCacherBus
	Timers            *timersBus
	Eval              *evalBus
	EventSub          *eventSubBus
	Scheduler         *schedulerBus
	ChatMessages      Queue[twitch.TwitchChatMessage, struct{}]
	ChatMessagesStore *chatMessagesStoreBus
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

			GetBuiltInVariables: NewNatsQueue[struct{}, []parser.BuiltInVariable](
				nc,
				parser.GetBuiltInVariablesSubject,
				5*time.Second,
				nats.GOB_ENCODER,
			),
		},

		Bots: &botsBus{
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
			BanUser: NewNatsQueue[botsservice.BanRequest, struct{}](
				nc,
				botsservice.BanSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
		},

		Websocket: &websocketBus{
			DudesGrow: NewNatsQueue[websockets.DudesGrowRequest, struct{}](
				nc,
				websockets.DudesGrowSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),

			DudesUserSettings: NewNatsQueue[websockets.DudesChangeUserSettingsRequest, struct{}](
				nc,
				websockets.DudesUserSettingsSubjsect,
				1*time.Minute,
				nats.GOB_ENCODER,
			),

			DudesLeave: NewNatsQueue[websockets.DudesLeaveRequest, struct{}](
				nc,
				websockets.DudesLeaveSubject,
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
			SubscribeToAllEvents: NewNatsQueue[eventsub.EventsubSubscribeToAllEventsRequest, struct{}](
				nc,
				eventsub.EventsubSubscribeAllSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			Subscribe: NewNatsQueue[eventsub.EventsubSubscribeRequest, struct{}](
				nc,
				eventsub.EventsubSubscribeSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
		},

		Scheduler: &schedulerBus{
			CreateDefaultCommands: NewNatsQueue[scheduler.CreateDefaultCommandsRequest, struct{}](
				nc,
				scheduler.CreateDefaultCommandsSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			CreateDefaultRoles: NewNatsQueue[scheduler.CreateDefaultRolesRequest, struct{}](
				nc,
				scheduler.CreateDefaultRolesSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
		},

		ChatMessages: NewNatsQueue[twitch.TwitchChatMessage, struct{}](
			nc,
			CHAT_MESSAGES_SUBJECT,
			30*time.Minute,
			nats.JSON_ENCODER,
		),
		ChatMessagesStore: &chatMessagesStoreBus{
			GetChatMessagesByTextForDelete: NewNatsQueue[chat_messages_store.GetChatMessagesByTextRequest, chat_messages_store.GetChatMessagesByTextResponse](
				nc,
				CHAT_MESSAGES_STORE_GET_BY_TEXT_FOR_DELETE_SUBJECT,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			RemoveMessages: NewNatsQueue[chat_messages_store.RemoveMessagesRequest, struct{}](
				nc,
				CHAT_MESSAGES_STRORE_REMOVE_MESSAGES_SUBJECT,
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
