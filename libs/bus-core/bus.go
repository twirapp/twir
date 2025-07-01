package buscore

import (
	"time"

	"github.com/nats-io/nats.go"
	cfg "github.com/satont/twir/libs/config"
	auditlog "github.com/twirapp/twir/libs/bus-core/audit-logs"
	botsservice "github.com/twirapp/twir/libs/bus-core/bots"
	emotes_cacher "github.com/twirapp/twir/libs/bus-core/emotes-cacher"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/bus-core/eventsub"
	"github.com/twirapp/twir/libs/bus-core/giveaways"
	"github.com/twirapp/twir/libs/bus-core/parser"
	"github.com/twirapp/twir/libs/bus-core/scheduler"
	"github.com/twirapp/twir/libs/bus-core/timers"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/bus-core/websockets"
	"github.com/twirapp/twir/libs/bus-core/ytsr"
)

type Bus struct {
	AuditLogs     *auditLogsBus
	Parser        *parserBus
	Websocket     *websocketBus
	Channel       *channelBus
	Bots          *botsBus
	EmotesCacher  *emotesCacherBus
	Timers        *timersBus
	EventSub      *eventSubBus
	Scheduler     *schedulerBus
	Giveaways     *giveawaysBus
	ChatMessages  Queue[twitch.TwitchChatMessage, struct{}]
	RedemptionAdd Queue[twitch.ActivatedRedemption, struct{}]
	Events        *eventsBus
	YTSRSearch    Queue[ytsr.SearchRequest, ytsr.SearchResponse]
}

func NewNatsBus(nc *nats.Conn) *Bus {
	return &Bus{
		Giveaways: &giveawaysBus{
			TryAddParticipant: NewNatsQueue[giveaways.TryAddParticipantRequest, struct{}](
				nc,
				giveaways.TryAddParticipantSubject,
				1*time.Minute,
				GobEncoder,
			),
			ChooseWinner: NewNatsQueue[giveaways.ChooseWinnerRequest, giveaways.ChooseWinnerResponse](
				nc,
				giveaways.ChooseWinnerSubject,
				1*time.Minute,
				GobEncoder,
			),
			NewParticipants: NewNatsQueue[giveaways.NewParticipant, struct{}](
				nc,
				giveaways.NewParticipantsSubject,
				1*time.Minute,
				GobEncoder,
			),
		},
		AuditLogs: &auditLogsBus{
			Logs: NewNatsQueue[auditlog.NewAuditLogMessage, struct{}](
				nc,
				AUDIT_LOGS_SUBJECT,
				1*time.Minute,
				GobEncoder,
			),
		},

		Parser: &parserBus{
			GetCommandResponse: NewNatsQueue[twitch.TwitchChatMessage, parser.CommandParseResponse](
				nc,
				PARSER_COMMANDS_SUBJECT,
				30*time.Minute,
				GobEncoder,
			),

			ParseVariablesInText: NewNatsQueue[parser.ParseVariablesInTextRequest, parser.ParseVariablesInTextResponse](
				nc,
				PARSER_TEXT_VARIABLES_SUBJECT,
				1*time.Minute,
				GobEncoder,
			),

			ProcessMessageAsCommand: NewNatsQueue[twitch.TwitchChatMessage, struct{}](
				nc,
				PARSER_PROCESS_MESSAGE_AS_COMMAND_SUBJECT,
				30*time.Minute,
				GobEncoder,
			),

			GetBuiltInVariables: NewNatsQueue[struct{}, []parser.BuiltInVariable](
				nc,
				parser.GetBuiltInVariablesSubject,
				5*time.Second,
				GobEncoder,
			),
		},

		Bots: &botsBus{
			SendMessage: NewNatsQueue[botsservice.SendMessageRequest, struct{}](
				nc,
				botsservice.SendMessageSubject,
				1*time.Minute,
				GobEncoder,
			),
			DeleteMessage: NewNatsQueue[botsservice.DeleteMessageRequest, struct{}](
				nc,
				botsservice.DeleteMessageSubject,
				1*time.Minute,
				GobEncoder,
			),
			BanUser: NewNatsQueue[botsservice.BanRequest, struct{}](
				nc,
				botsservice.BanSubject,
				1*time.Minute,
				GobEncoder,
			),
			BanUsers: NewNatsQueue[[]botsservice.BanRequest, struct{}](
				nc,
				botsservice.BanMultipleSubject,
				5*time.Minute,
				GobEncoder,
			),
			ShoutOut: NewNatsQueue[botsservice.SentShoutOutRequest, struct{}](
				nc,
				botsservice.ShoutOutSubject,
				1*time.Minute,
				GobEncoder,
			),
			Vip: NewNatsQueue[botsservice.VipRequest, struct{}](
				nc,
				botsservice.VipSubject,
				1*time.Minute,
				GobEncoder,
			),
			UnVip: NewNatsQueue[botsservice.UnVipRequest, struct{}](
				nc,
				botsservice.UnVipSubject,
				1*time.Minute,
				GobEncoder,
			),
		},

		Websocket: &websocketBus{
			DudesGrow: NewNatsQueue[websockets.DudesGrowRequest, struct{}](
				nc,
				websockets.DudesGrowSubject,
				1*time.Minute,
				GobEncoder,
			),

			DudesUserSettings: NewNatsQueue[websockets.DudesChangeUserSettingsRequest, struct{}](
				nc,
				websockets.DudesUserSettingsSubjsect,
				1*time.Minute,
				GobEncoder,
			),

			DudesLeave: NewNatsQueue[websockets.DudesLeaveRequest, struct{}](
				nc,
				websockets.DudesLeaveSubject,
				1*time.Minute,
				GobEncoder,
			),
		},

		Channel: &channelBus{
			StreamOnline: NewNatsQueue[twitch.StreamOnlineMessage, struct{}](
				nc,
				STREAM_ONLINE_SUBJECT,
				1*time.Minute,
				GobEncoder,
			),
			StreamOffline: NewNatsQueue[twitch.StreamOfflineMessage, struct{}](
				nc,
				STREAM_OFFLINE_SUBJECT,
				1*time.Minute,
				GobEncoder,
			),
		},

		EmotesCacher: &emotesCacherBus{
			CacheGlobalEmotes: NewNatsQueue[struct{}, struct{}](
				nc,
				emotes_cacher.EMOTES_CACHER_GLOBAL_EMOTES_SUBJECT,
				1*time.Minute,
				GobEncoder,
			),
			CacheChannelEmotes: NewNatsQueue[emotes_cacher.EmotesCacheRequest, struct{}](
				nc,
				emotes_cacher.EMOTES_CACHER_CHANNEL_EMOTES_SUBJECT,
				1*time.Minute,
				GobEncoder,
			),
		},

		Timers: &timersBus{
			AddTimer: NewNatsQueue[timers.AddOrRemoveTimerRequest, struct{}](
				nc,
				timers.AddTimerSubject,
				1*time.Minute,
				GobEncoder,
			),
			RemoveTimer: NewNatsQueue[timers.AddOrRemoveTimerRequest, struct{}](
				nc,
				timers.RemoveTimerSubject,
				1*time.Minute,
				GobEncoder,
			),
		},

		EventSub: &eventSubBus{
			SubscribeToAllEvents: NewNatsQueue[eventsub.EventsubSubscribeToAllEventsRequest, struct{}](
				nc,
				eventsub.EventsubSubscribeAllSubject,
				1*time.Minute,
				GobEncoder,
			),
			Subscribe: NewNatsQueue[eventsub.EventsubSubscribeRequest, struct{}](
				nc,
				eventsub.EventsubSubscribeSubject,
				1*time.Minute,
				GobEncoder,
			),
			InitChannels: NewNatsQueue[struct{}, struct{}](
				nc,
				eventsub.EventsubInitChannelsSubject,
				1*time.Minute,
				GobEncoder,
			),
			Unsubscribe: NewNatsQueue[string, struct{}](
				nc,
				eventsub.EventsubUnsubscribeSubject,
				1*time.Minute,
				GobEncoder,
			),
		},

		Scheduler: &schedulerBus{
			CreateDefaultCommands: NewNatsQueue[scheduler.CreateDefaultCommandsRequest, struct{}](
				nc,
				scheduler.CreateDefaultCommandsSubject,
				1*time.Minute,
				GobEncoder,
			),
			CreateDefaultRoles: NewNatsQueue[scheduler.CreateDefaultRolesRequest, struct{}](
				nc,
				scheduler.CreateDefaultRolesSubject,
				1*time.Minute,
				GobEncoder,
			),
		},

		ChatMessages: NewNatsQueue[twitch.TwitchChatMessage, struct{}](
			nc,
			CHAT_MESSAGES_SUBJECT,
			30*time.Minute,
			JsonEncoder,
		),

		RedemptionAdd: NewNatsQueue[twitch.ActivatedRedemption, struct{}](
			nc,
			twitch.RedemptionAddSubject,
			30*time.Minute,
			GobEncoder,
		),

		Events: &eventsBus{
			Follow: NewNatsQueue[events.FollowMessage, struct{}](
				nc,
				events.FollowSubject,
				1*time.Minute,
				GobEncoder,
			),
			Subscribe: NewNatsQueue[events.SubscribeMessage, struct{}](
				nc,
				events.SubscribeSubject,
				1*time.Minute,
				GobEncoder,
			),
			SubGift: NewNatsQueue[events.SubGiftMessage, struct{}](
				nc,
				events.SubGiftSubject,
				1*time.Minute,
				GobEncoder,
			),
			ReSubscribe: NewNatsQueue[events.ReSubscribeMessage, struct{}](
				nc,
				events.ReSubscribeSubject,
				1*time.Minute,
				GobEncoder,
			),
			RedemptionCreated: NewNatsQueue[events.RedemptionCreatedMessage, struct{}](
				nc,
				events.RedemptionCreatedSubject,
				1*time.Minute,
				GobEncoder,
			),
			CommandUsed: NewNatsQueue[events.CommandUsedMessage, struct{}](
				nc,
				events.CommandUsedSubject,
				1*time.Minute,
				GobEncoder,
			),
			FirstUserMessage: NewNatsQueue[events.FirstUserMessageMessage, struct{}](
				nc,
				events.FirstUserMessageSubject,
				1*time.Minute,
				GobEncoder,
			),
			Raided: NewNatsQueue[events.RaidedMessage, struct{}](
				nc,
				events.RaidedSubject,
				1*time.Minute,
				GobEncoder,
			),
			TitleOrCategoryChanged: NewNatsQueue[events.TitleOrCategoryChangedMessage, struct{}](
				nc,
				events.TitleOrCategoryChangedSubject,
				1*time.Minute,
				GobEncoder,
			),
			ChatClear: NewNatsQueue[events.ChatClearMessage, struct{}](
				nc,
				events.ChatClearSubject,
				1*time.Minute,
				GobEncoder,
			),
			Donate: NewNatsQueue[events.DonateMessage, struct{}](
				nc,
				events.DonateSubject,
				1*time.Minute,
				GobEncoder,
			),
			KeywordMatched: NewNatsQueue[events.KeywordMatchedMessage, struct{}](
				nc,
				events.KeywordMatchedSubject,
				1*time.Minute,
				GobEncoder,
			),
			GreetingSended: NewNatsQueue[events.GreetingSendedMessage, struct{}](
				nc,
				events.GreetingSendedSubject,
				1*time.Minute,
				GobEncoder,
			),
			PollBegin: NewNatsQueue[events.PollBeginMessage, struct{}](
				nc,
				events.PollBeginSubject,
				1*time.Minute,
				GobEncoder,
			),
			PollProgress: NewNatsQueue[events.PollProgressMessage, struct{}](
				nc,
				events.PollProgressSubject,
				1*time.Minute,
				GobEncoder,
			),
			PollEnd: NewNatsQueue[events.PollEndMessage, struct{}](
				nc,
				events.PollEndSubject,
				1*time.Minute,
				GobEncoder,
			),
			PredictionBegin: NewNatsQueue[events.PredictionBeginMessage, struct{}](
				nc,
				events.PredictionBeginSubject,
				1*time.Minute,
				GobEncoder,
			),
			PredictionProgress: NewNatsQueue[events.PredictionProgressMessage, struct{}](
				nc,
				events.PredictionProgressSubject,
				1*time.Minute,
				GobEncoder,
			),
			PredictionLock: NewNatsQueue[events.PredictionLockMessage, struct{}](
				nc,
				events.PredictionLockSubject,
				1*time.Minute,
				GobEncoder,
			),
			PredictionEnd: NewNatsQueue[events.PredictionEndMessage, struct{}](
				nc,
				events.PredictionEndSubject,
				1*time.Minute,
				GobEncoder,
			),
			StreamFirstUserJoin: NewNatsQueue[events.StreamFirstUserJoinMessage, struct{}](
				nc,
				events.StreamFirstUserJoinSubject,
				1*time.Minute,
				GobEncoder,
			),
			ChannelBan: NewNatsQueue[events.ChannelBanMessage, struct{}](
				nc,
				events.ChannelBanSubject,
				1*time.Minute,
				GobEncoder,
			),
			ChannelUnbanRequestCreate: NewNatsQueue[events.ChannelUnbanRequestCreateMessage, struct{}](
				nc,
				events.ChannelUnbanRequestCreateSubject,
				1*time.Minute,
				GobEncoder,
			),
			ChannelUnbanRequestResolve: NewNatsQueue[events.ChannelUnbanRequestResolveMessage, struct{}](
				nc,
				events.ChannelUnbanRequestResolveSubject,
				1*time.Minute,
				GobEncoder,
			),
			ChannelMessageDelete: NewNatsQueue[events.ChannelMessageDeleteMessage, struct{}](
				nc,
				events.ChannelMessageDeleteSubject,
				1*time.Minute,
				GobEncoder,
			),
			VipAdded: NewNatsQueue[events.VipAddedMessage, struct{}](
				nc,
				events.VipAddedSubject,
				1*time.Minute,
				GobEncoder,
			),
			VipRemoved: NewNatsQueue[events.VipRemovedMessage, struct{}](
				nc,
				events.VipRemovedSubject,
				1*time.Minute,
				GobEncoder,
			),
			ModeratorAdded: NewNatsQueue[events.ModeratorAddedMessage, struct{}](
				nc,
				events.ModeratorAddedSubject,
				1*time.Minute,
				GobEncoder,
			),
			ModeratorRemoved: NewNatsQueue[events.ModeratorRemovedMessage, struct{}](
				nc,
				events.ModeratorRemovedSubject,
				1*time.Minute,
				GobEncoder,
			),
		},
		YTSRSearch: NewNatsQueue[ytsr.SearchRequest, ytsr.SearchResponse](
			nc,
			ytsr.SearchSubject,
			1*time.Minute,
			JsonEncoder,
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
