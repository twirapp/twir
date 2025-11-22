package buscore

import (
	"github.com/twirapp/twir/libs/bus-core/api"
	auditlogs "github.com/twirapp/twir/libs/bus-core/audit-logs"
	botsservice "github.com/twirapp/twir/libs/bus-core/bots"
	emotes_cacher "github.com/twirapp/twir/libs/bus-core/emotes-cacher"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/bus-core/eventsub"
	"github.com/twirapp/twir/libs/bus-core/giveaways"
	"github.com/twirapp/twir/libs/bus-core/integrations"
	"github.com/twirapp/twir/libs/bus-core/parser"
	"github.com/twirapp/twir/libs/bus-core/scheduler"
	"github.com/twirapp/twir/libs/bus-core/timers"
	"github.com/twirapp/twir/libs/bus-core/tokens"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/bus-core/websockets"
)

type auditLogsBus struct {
	Logs Queue[auditlogs.NewAuditLogMessage, struct{}]
}

type parserBus struct {
	GetCommandResponse      Queue[twitch.TwitchChatMessage, parser.CommandParseResponse]
	ProcessMessageAsCommand Queue[twitch.TwitchChatMessage, struct{}]
	ParseVariablesInText    Queue[parser.ParseVariablesInTextRequest, parser.ParseVariablesInTextResponse]
	GetBuiltInVariables     Queue[struct{}, []parser.BuiltInVariable]
	GetDefaultCommands      Queue[struct{}, parser.GetDefaultCommandsResponse]
}

type websocketBus struct {
	DudesGrow         Queue[websockets.DudesGrowRequest, struct{}]
	DudesUserSettings Queue[websockets.DudesChangeUserSettingsRequest, struct{}]
	DudesLeave        Queue[websockets.DudesLeaveRequest, struct{}]
}

type channelBus struct {
	StreamOnline  Queue[twitch.StreamOnlineMessage, struct{}]
	StreamOffline Queue[twitch.StreamOfflineMessage, struct{}]
}

type botsBus struct {
	SendMessage     Queue[botsservice.SendMessageRequest, struct{}]
	DeleteMessage   Queue[botsservice.DeleteMessageRequest, struct{}]
	BanUser         Queue[botsservice.BanRequest, struct{}]
	BanUsers        Queue[[]botsservice.BanRequest, struct{}]
	ShoutOut        Queue[botsservice.SentShoutOutRequest, struct{}]
	Vip             Queue[botsservice.VipRequest, struct{}]
	UnVip           Queue[botsservice.UnVipRequest, struct{}]
	ModeratorAdd    Queue[botsservice.ModeratorAddRequest, struct{}]
	ModeratorRemove Queue[botsservice.ModeratorRemoveRequest, struct{}]
}

type emotesCacherBus struct {
	GetGlobalEmotes  Queue[emotes_cacher.GetGlobalEmotesRequest, emotes_cacher.Response]
	GetChannelEmotes Queue[emotes_cacher.GetChannelEmotesRequest, emotes_cacher.Response]
}

type timersBus struct {
	AddTimer    Queue[timers.AddOrRemoveTimerRequest, struct{}]
	RemoveTimer Queue[timers.AddOrRemoveTimerRequest, struct{}]
}

type eventSubBus struct {
	SubscribeToAllEvents Queue[eventsub.EventsubSubscribeToAllEventsRequest, struct{}]
	Subscribe            Queue[eventsub.EventsubSubscribeRequest, struct{}]
	// Init channels is dangerous, only use it if you know what you're doing
	InitChannels Queue[struct{}, struct{}]
	Unsubscribe  Queue[string, struct{}]
}

type schedulerBus struct {
	CreateDefaultCommands Queue[scheduler.CreateDefaultCommandsRequest, struct{}]
	CreateDefaultRoles    Queue[scheduler.CreateDefaultRolesRequest, struct{}]
}

type giveawaysBus struct {
	TryAddParticipant Queue[giveaways.TryAddParticipantRequest, struct{}]
	ChooseWinner      Queue[giveaways.ChooseWinnerRequest, giveaways.ChooseWinnerResponse]

	NewParticipants Queue[giveaways.NewParticipant, struct{}]
}

type tokensBus struct {
	RequestAppToken  Queue[struct{}, tokens.TokenResponse]
	RequestUserToken Queue[tokens.GetUserTokenRequest, tokens.TokenResponse]
	RequestBotToken  Queue[tokens.GetBotTokenRequest, tokens.TokenResponse]
	UpdateToken      Queue[tokens.UpdateTokenRequest, struct{}]
}

type integrationsBus struct {
	Add    Queue[integrations.Request, struct{}]
	Remove Queue[integrations.Request, struct{}]
}

type eventsBus struct {
	Follow                     Queue[events.FollowMessage, struct{}]
	Subscribe                  Queue[events.SubscribeMessage, struct{}]
	SubGift                    Queue[events.SubGiftMessage, struct{}]
	ReSubscribe                Queue[events.ReSubscribeMessage, struct{}]
	RedemptionCreated          Queue[events.RedemptionCreatedMessage, struct{}]
	CommandUsed                Queue[events.CommandUsedMessage, struct{}]
	FirstUserMessage           Queue[events.FirstUserMessageMessage, struct{}]
	Raided                     Queue[events.RaidedMessage, struct{}]
	TitleOrCategoryChanged     Queue[events.TitleOrCategoryChangedMessage, struct{}]
	ChatClear                  Queue[events.ChatClearMessage, struct{}]
	Donate                     Queue[events.DonateMessage, struct{}]
	KeywordMatched             Queue[events.KeywordMatchedMessage, struct{}]
	GreetingSended             Queue[events.GreetingSendedMessage, struct{}]
	PollBegin                  Queue[events.PollBeginMessage, struct{}]
	PollProgress               Queue[events.PollProgressMessage, struct{}]
	PollEnd                    Queue[events.PollEndMessage, struct{}]
	PredictionBegin            Queue[events.PredictionBeginMessage, struct{}]
	PredictionProgress         Queue[events.PredictionProgressMessage, struct{}]
	PredictionLock             Queue[events.PredictionLockMessage, struct{}]
	PredictionEnd              Queue[events.PredictionEndMessage, struct{}]
	StreamFirstUserJoin        Queue[events.StreamFirstUserJoinMessage, struct{}]
	ChannelBan                 Queue[events.ChannelBanMessage, struct{}]
	ChannelUnban               Queue[events.ChannelUnbanMessage, struct{}]
	ChannelUnbanRequestCreate  Queue[events.ChannelUnbanRequestCreateMessage, struct{}]
	ChannelUnbanRequestResolve Queue[events.ChannelUnbanRequestResolveMessage, struct{}]
	ChannelMessageDelete       Queue[events.ChannelMessageDeleteMessage, struct{}]
	VipAdded                   Queue[events.VipAddedMessage, struct{}]
	VipRemoved                 Queue[events.VipRemovedMessage, struct{}]
	ModeratorAdded             Queue[events.ModeratorAddedMessage, struct{}]
	ModeratorRemoved           Queue[events.ModeratorRemovedMessage, struct{}]
}

type apiBus struct {
	TriggerKappagen Queue[api.TriggerKappagenMessage, struct{}]
	TriggerBrbStart Queue[api.TriggerBrbStart, struct{}]
	TriggerBrbStop  Queue[api.TriggerBrbStop, struct{}]
	TriggerTtsSay   Queue[api.TriggerTtsSay, struct{}]
	TriggerTtsSkip  Queue[api.TriggerTtsSkip, struct{}]
}
