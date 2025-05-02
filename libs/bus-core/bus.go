package buscore

import (
	"time"

	"github.com/nats-io/nats.go"
	cfg "github.com/satont/twir/libs/config"
	auditlog "github.com/twirapp/twir/libs/bus-core/audit-logs"
	botsservice "github.com/twirapp/twir/libs/bus-core/bots"
	emotes_cacher "github.com/twirapp/twir/libs/bus-core/emotes-cacher"
	"github.com/twirapp/twir/libs/bus-core/eval"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/bus-core/eventsub"
	"github.com/twirapp/twir/libs/bus-core/giveaways"
	"github.com/twirapp/twir/libs/bus-core/parser"
	"github.com/twirapp/twir/libs/bus-core/scheduler"
	"github.com/twirapp/twir/libs/bus-core/timers"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/bus-core/websockets"
)

type Bus struct {
	AuditLogs     *auditLogsBus
	Parser        *parserBus
	Websocket     *websocketBus
	Channel       *channelBus
	Bots          *botsBus
	EmotesCacher  *emotesCacherBus
	Timers        *timersBus
	Eval          *evalBus
	EventSub      *eventSubBus
	Scheduler     *schedulerBus
	Giveaways     *giveawaysBus
	ChatMessages  Queue[twitch.TwitchChatMessage, struct{}]
	RedemptionAdd Queue[twitch.ActivatedRedemption, struct{}]
	Events        *eventsBus
}

func NewNatsBus(nc *nats.Conn) *Bus {
	return &Bus{
		Giveaways: &giveawaysBus{
			TryAddParticipant: NewNatsQueue[giveaways.TryAddParticipantRequest, struct{}](
				nc,
				giveaways.TryAddParticipantSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			ChooseWinner: NewNatsQueue[giveaways.ChooseWinnerRequest, giveaways.ChooseWinnerResponse](
				nc,
				giveaways.ChooseWinnerSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			NewParticipants: NewNatsQueue[giveaways.NewParticipant, struct{}](
				nc,
				giveaways.NewParticipantsSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
		},
		AuditLogs: &auditLogsBus{
			Logs: NewNatsQueue[auditlog.NewAuditLogMessage, struct{}](
				nc,
				AUDIT_LOGS_SUBJECT,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
		},

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
			BanUsers: NewNatsQueue[[]botsservice.BanRequest, struct{}](
				nc,
				botsservice.BanMultipleSubject,
				5*time.Minute,
				nats.GOB_ENCODER,
			),
			ShoutOut: NewNatsQueue[botsservice.SentShoutOutRequest, struct{}](
				nc,
				botsservice.ShoutOutSubject,
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
			InitChannels: NewNatsQueue[struct{}, struct{}](
				nc,
				eventsub.EventsubInitChannelsSubject,
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

		RedemptionAdd: NewNatsQueue[twitch.ActivatedRedemption, struct{}](
			nc,
			twitch.RedemptionAddSubject,
			30*time.Minute,
			nats.GOB_ENCODER,
		),

		Events: &eventsBus{
			Follow: NewNatsQueue[events.FollowMessage, struct{}](
				nc,
				events.FollowSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			Subscribe: NewNatsQueue[events.SubscribeMessage, struct{}](
				nc,
				events.SubscribeSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			SubGift: NewNatsQueue[events.SubGiftMessage, struct{}](
				nc,
				events.SubGiftSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			ReSubscribe: NewNatsQueue[events.ReSubscribeMessage, struct{}](
				nc,
				events.ReSubscribeSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			RedemptionCreated: NewNatsQueue[events.RedemptionCreatedMessage, struct{}](
				nc,
				events.RedemptionCreatedSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			CommandUsed: NewNatsQueue[events.CommandUsedMessage, struct{}](
				nc,
				events.CommandUsedSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			FirstUserMessage: NewNatsQueue[events.FirstUserMessageMessage, struct{}](
				nc,
				events.FirstUserMessageSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			Raided: NewNatsQueue[events.RaidedMessage, struct{}](
				nc,
				events.RaidedSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			TitleOrCategoryChanged: NewNatsQueue[events.TitleOrCategoryChangedMessage, struct{}](
				nc,
				events.TitleOrCategoryChangedSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			ChatClear: NewNatsQueue[events.ChatClearMessage, struct{}](
				nc,
				events.ChatClearSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			Donate: NewNatsQueue[events.DonateMessage, struct{}](
				nc,
				events.DonateSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			KeywordMatched: NewNatsQueue[events.KeywordMatchedMessage, struct{}](
				nc,
				events.KeywordMatchedSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			GreetingSended: NewNatsQueue[events.GreetingSendedMessage, struct{}](
				nc,
				events.GreetingSendedSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			PollBegin: NewNatsQueue[events.PollBeginMessage, struct{}](
				nc,
				events.PollBeginSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			PollProgress: NewNatsQueue[events.PollProgressMessage, struct{}](
				nc,
				events.PollProgressSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			PollEnd: NewNatsQueue[events.PollEndMessage, struct{}](
				nc,
				events.PollEndSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			PredictionBegin: NewNatsQueue[events.PredictionBeginMessage, struct{}](
				nc,
				events.PredictionBeginSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			PredictionProgress: NewNatsQueue[events.PredictionProgressMessage, struct{}](
				nc,
				events.PredictionProgressSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			PredictionLock: NewNatsQueue[events.PredictionLockMessage, struct{}](
				nc,
				events.PredictionLockSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			PredictionEnd: NewNatsQueue[events.PredictionEndMessage, struct{}](
				nc,
				events.PredictionEndSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			StreamFirstUserJoin: NewNatsQueue[events.StreamFirstUserJoinMessage, struct{}](
				nc,
				events.StreamFirstUserJoinSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			ChannelBan: NewNatsQueue[events.ChannelBanMessage, struct{}](
				nc,
				events.ChannelBanSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			ChannelUnbanRequestCreate: NewNatsQueue[events.ChannelUnbanRequestCreateMessage, struct{}](
				nc,
				events.ChannelUnbanRequestCreateSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			ChannelUnbanRequestResolve: NewNatsQueue[events.ChannelUnbanRequestResolveMessage, struct{}](
				nc,
				events.ChannelUnbanRequestResolveSubject,
				1*time.Minute,
				nats.GOB_ENCODER,
			),
			ChannelMessageDelete: NewNatsQueue[events.ChannelMessageDeleteMessage, struct{}](
				nc,
				events.ChannelMessageDeleteSubject,
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
