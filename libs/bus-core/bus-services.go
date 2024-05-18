package buscore

import (
	botsservice "github.com/twirapp/twir/libs/bus-core/bots"
	emotes_cacher "github.com/twirapp/twir/libs/bus-core/emotes-cacher"
	"github.com/twirapp/twir/libs/bus-core/eval"
	"github.com/twirapp/twir/libs/bus-core/eventsub"
	"github.com/twirapp/twir/libs/bus-core/parser"
	"github.com/twirapp/twir/libs/bus-core/scheduler"
	"github.com/twirapp/twir/libs/bus-core/timers"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/bus-core/websockets"
)

type parserBus struct {
	GetCommandResponse      Queue[twitch.TwitchChatMessage, parser.CommandParseResponse]
	ProcessMessageAsCommand Queue[twitch.TwitchChatMessage, struct{}]
	ParseVariablesInText    Queue[parser.ParseVariablesInTextRequest, parser.ParseVariablesInTextResponse]
	GetBuiltInVariables     Queue[struct{}, []parser.BuiltInVariable]
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
	ProcessMessage Queue[twitch.TwitchChatMessage, struct{}]
	SendMessage    Queue[botsservice.SendMessageRequest, struct{}]
	DeleteMessage  Queue[botsservice.DeleteMessageRequest, struct{}]
	BanUser        Queue[botsservice.BanRequest, struct{}]
}

type emotesCacherBus struct {
	CacheGlobalEmotes  Queue[struct{}, struct{}]
	CacheChannelEmotes Queue[emotes_cacher.EmotesCacheRequest, struct{}]
}

type timersBus struct {
	AddTimer    Queue[timers.AddOrRemoveTimerRequest, struct{}]
	RemoveTimer Queue[timers.AddOrRemoveTimerRequest, struct{}]
}

type evalBus struct {
	Evaluate Queue[eval.EvalRequest, eval.EvalResponse]
}

type eventSubBus struct {
	Subscribe Queue[eventsub.EventsubSubscribeRequest, struct{}]
}

type schedulerBus struct {
	CreateDefaultCommands Queue[scheduler.CreateDefaultCommandsRequest, struct{}]
	CreateDefaultRoles    Queue[scheduler.CreateDefaultRolesRequest, struct{}]
}
