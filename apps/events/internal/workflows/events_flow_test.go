package workflows

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	kvinmemory "github.com/twirapp/kv/stores/inmemory"
	eventsactivity "github.com/twirapp/twir/apps/events/internal/activities/events"
	"github.com/twirapp/twir/apps/events/internal/shared"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/repositories/events/model"
	streamsrepository "github.com/twirapp/twir/libs/repositories/streams"
	streamsmodel "github.com/twirapp/twir/libs/repositories/streams/model"
	"go.temporal.io/sdk/testsuite"
)

func TestGetEventChannelBindingsSelectsEventAndTwitchBindingsByPlatform(t *testing.T) {
	twitchUserID := uuid.New()
	channel := channelentity.Channel{
		Bindings: []channelplatformentity.ChannelPlatform{
			{
				Platform:          platform.PlatformTwitch,
				PlatformChannelID: "twitch-channel",
				UserID:            twitchUserID,
				Enabled:           false,
				BotConfig: json.RawMessage(
					`{"bot_id":"twitch-bot","is_bot_mod":true,"is_twitch_banned":true}`,
				),
			},
			{
				Platform:          platform.PlatformKick,
				PlatformChannelID: "kick-channel",
				Enabled:           true,
			},
		},
	}

	bindings, err := getEventChannelBindings(channel, platform.PlatformKick)
	if err != nil {
		t.Fatalf("getEventChannelBindings returned error: %v", err)
	}
	if bindings.event.PlatformChannelID != "kick-channel" {
		t.Errorf("event PlatformChannelID = %q, want %q", bindings.event.PlatformChannelID, "kick-channel")
	}
	if !bindings.event.Enabled {
		t.Error("event binding Enabled = false, want true")
	}
	if !bindings.hasTwitch {
		t.Fatal("expected Twitch binding")
	}
	if bindings.twitch.UserID != twitchUserID {
		t.Errorf("Twitch UserID = %s, want %s", bindings.twitch.UserID, twitchUserID)
	}
	if bindings.twitchBotConfig.BotID != "twitch-bot" {
		t.Errorf("Twitch bot ID = %q, want %q", bindings.twitchBotConfig.BotID, "twitch-bot")
	}
	if !bindings.twitchBotConfig.IsBotMod {
		t.Error("Twitch bot mod state = false, want true")
	}
	if !bindings.twitchBotConfig.IsTwitchBanned {
		t.Error("Twitch ban state = false, want true")
	}

	data := bindings.applyTo(shared.EventData{
		ChannelID: "incoming-kick-channel",
		Platform:  platform.PlatformKick,
	})
	if data.ChannelID != "kick-channel" {
		t.Errorf("event ChannelID = %q, want %q", data.ChannelID, "kick-channel")
	}
	if data.ChannelTwitchPlatformID != "twitch-channel" {
		t.Errorf("Twitch broadcaster ID = %q, want %q", data.ChannelTwitchPlatformID, "twitch-channel")
	}
	if data.ChannelTwitchUserID != twitchUserID.String() {
		t.Errorf("Twitch user ID = %q, want %q", data.ChannelTwitchUserID, twitchUserID)
	}
}

func TestIsTwitchBannedEventScopesBanStateToTwitch(t *testing.T) {
	tests := []struct {
		name           string
		eventPlatform  platform.Platform
		wantSuppressed bool
	}{
		{
			name:           "Twitch event is suppressed",
			eventPlatform:  platform.PlatformTwitch,
			wantSuppressed: true,
		},
		{
			name:           "Kick event flows",
			eventPlatform:  platform.PlatformKick,
			wantSuppressed: false,
		},
		{
			name:           "VK event flows",
			eventPlatform:  platform.PlatformVKVideoLive,
			wantSuppressed: false,
		},
		{
			name:           "platform-neutral event flows",
			eventPlatform:  "",
			wantSuppressed: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isTwitchBanned := true
			if got := isTwitchBannedEvent(test.eventPlatform, isTwitchBanned); got != test.wantSuppressed {
				t.Errorf("isTwitchBannedEvent(%q, true) = %t, want %t", test.eventPlatform, got, test.wantSuppressed)
			}
		})
	}
}

func TestFlowAppliesTwitchBanGateBeforeSelectingOperations(t *testing.T) {
	channelID := uuid.New()
	channel := channelentity.Channel{
		ID: channelID,
		Bindings: []channelplatformentity.ChannelPlatform{
			{
				Platform:          platform.PlatformTwitch,
				PlatformChannelID: "twitch-channel",
				Enabled:           true,
				BotConfig:         json.RawMessage(`{"is_twitch_banned":true}`),
			},
			{
				Platform:          platform.PlatformKick,
				PlatformChannelID: "kick-channel",
				Enabled:           true,
			},
			{
				Platform:          platform.PlatformVKVideoLive,
				PlatformChannelID: "vk-channel",
				Enabled:           true,
			},
		},
	}
	channelCache := generic_cacher.New[channelentity.Channel](generic_cacher.Opts[channelentity.Channel]{
		KV:        kvinmemory.New(),
		KeyPrefix: "events-flow-channel:",
		LoadFn: func(context.Context, string) (channelentity.Channel, error) {
			return channelentity.Nil, nil
		},
	})
	if err := channelCache.SetValue(context.Background(), channelID.String(), channel); err != nil {
		t.Fatalf("set channel cache value: %v", err)
	}

	channelEventsCache := generic_cacher.New[[]model.Event](generic_cacher.Opts[[]model.Event]{
		KV:        kvinmemory.New(),
		KeyPrefix: "events-flow-events:",
		LoadFn: func(context.Context, string) ([]model.Event, error) {
			return nil, nil
		},
	})
	if err := channelEventsCache.SetValue(context.Background(), channelID.String(), []model.Event{{
		ID:      "event-id",
		Type:    model.EventTypeFollow,
		Enabled: true,
		Operations: []model.EventOperation{{
			ID:      "operation-id",
			Type:    model.EventOperationTypeSendMessage,
			Enabled: true,
		}},
	}}); err != nil {
		t.Fatalf("set event cache value: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{Addr: "unused"})
	redisClient.AddHook(eventFlowRedisHook{})
	t.Cleanup(func() { _ = redisClient.Close() })

	workflow := &EventWorkflow{
		channelsCache:                     channelCache,
		channelsEventsWithOperationsCache: channelEventsCache,
		streamsRepo:                       eventFlowStreamsRepository{},
		redis:                             redisClient,
		eventsActivity:                    &eventsactivity.Activity{},
	}

	tests := []struct {
		name          string
		eventPlatform platform.Platform
		wantOperation bool
	}{
		{
			name:          "Twitch event is suppressed",
			eventPlatform: platform.PlatformTwitch,
			wantOperation: false,
		},
		{
			name:          "Kick event flows",
			eventPlatform: platform.PlatformKick,
			wantOperation: true,
		},
		{
			name:          "VK event flows",
			eventPlatform: platform.PlatformVKVideoLive,
			wantOperation: true,
		},
		{
			name:          "platform-neutral event flows",
			eventPlatform: "",
			wantOperation: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			suite := testsuite.WorkflowTestSuite{}
			env := suite.NewTestWorkflowEnvironment()
			activityCall := env.OnActivity(
				workflow.eventsActivity.SendMessage,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(nil)
			if test.wantOperation {
				activityCall.Once()
			} else {
				activityCall.Never()
			}

			env.ExecuteWorkflow(
				workflow.Flow,
				model.EventTypeFollow,
				shared.EventData{
					ChannelID:   "incoming-channel",
					ChannelDBID: channelID.String(),
					Platform:    test.eventPlatform,
				},
			)
			if err := env.GetWorkflowError(); err != nil {
				t.Fatalf("Flow returned error: %v", err)
			}
			env.AssertExpectations(t)
		})
	}
}

type eventFlowRedisHook struct{}

func (eventFlowRedisHook) DialHook(next redis.DialHook) redis.DialHook { return next }

func (eventFlowRedisHook) ProcessHook(redis.ProcessHook) redis.ProcessHook {
	return func(context.Context, redis.Cmder) error { return nil }
}

func (eventFlowRedisHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return next
}

type eventFlowStreamsRepository struct{ streamsrepository.Repository }

func (eventFlowStreamsRepository) GetByChannelID(
	context.Context,
	uuid.UUID,
	platform.Platform,
) (streamsmodel.Stream, error) {
	return streamsmodel.Stream{}, nil
}
