package manager

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	kvinmemory "github.com/twirapp/kv/stores/inmemory"
	buscore "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/entities/platform"
	timersentity "github.com/twirapp/twir/libs/entities/timers"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	streamsrepository "github.com/twirapp/twir/libs/repositories/streams"
	streamsmodel "github.com/twirapp/twir/libs/repositories/streams/model"
	channelservice "github.com/twirapp/twir/libs/services/channels"
)

func TestTryTickLockedDoesNotAdvanceOfflineIntervalWithoutEligibleTarget(t *testing.T) {
	ctx := context.Background()
	channelID := uuid.New()
	channelCache := generic_cacher.New[channelmodel.Channel](
		generic_cacher.Opts[channelmodel.Channel]{
			KV:        kvinmemory.New(),
			KeyPrefix: "timer-test-channel:",
			LoadFn: func(context.Context, string) (channelmodel.Channel, error) {
				return channelmodel.Nil, nil
			},
		},
	)

	ineligibleChannel := channelmodel.Channel{
		ID: channelID,
		Bindings: []channelplatformsmodel.ChannelPlatform{
			{
				Platform:          platform.PlatformTwitch,
				PlatformChannelID: "twitch-channel",
				Enabled:           false,
				BotConfig:         json.RawMessage(`{"is_bot_mod":true}`),
			},
		},
	}
	if err := channelCache.SetValue(ctx, channelID.String(), ineligibleChannel); err != nil {
		t.Fatalf("set ineligible channel cache value: %v", err)
	}

	streamsRepository := &offlineTimerStreamsRepository{}
	manager := &Manager{
		channelCachedRepo: channelCache,
		channelservice: channelservice.NewChannelService(
			nil,
			&buscore.Bus{},
			cfg.Config{},
			kvinmemory.New(),
			streamsRepository,
		),
	}
	timer := &Timer{
		id: TimerID(uuid.New()),
		dbRow: timersentity.Timer{
			ChannelID:       channelID,
			OfflineEnabled:  true,
			MessageInterval: 3,
		},
	}

	manager.tryTickLocked(timer.id, timer)

	if streamsRepository.getListByChannelIDCalls != 0 {
		t.Fatalf("stream lookups while no target is eligible = %d, want 0", streamsRepository.getListByChannelIDCalls)
	}
	if timer.offlineMessageNumber != 0 {
		t.Fatalf("offline message number while no target is eligible = %d, want 0", timer.offlineMessageNumber)
	}

	eligibleChannel := ineligibleChannel
	eligibleChannel.Bindings = append([]channelplatformsmodel.ChannelPlatform(nil), ineligibleChannel.Bindings...)
	eligibleChannel.Bindings[0].Enabled = true
	if err := channelCache.SetValue(ctx, channelID.String(), eligibleChannel); err != nil {
		t.Fatalf("set eligible channel cache value: %v", err)
	}

	manager.tryTickLocked(timer.id, timer)

	if streamsRepository.getListByChannelIDCalls != 1 {
		t.Fatalf("stream lookups after target becomes eligible = %d, want 1", streamsRepository.getListByChannelIDCalls)
	}
	if timer.offlineMessageNumber != 1 {
		t.Fatalf("offline message number after target becomes eligible = %d, want 1", timer.offlineMessageNumber)
	}
	if timer.lastTriggerOfflineNumber != 0 {
		t.Fatalf("last triggered offline number = %d, want 0 before interval is reached", timer.lastTriggerOfflineNumber)
	}
}

type offlineTimerStreamsRepository struct {
	getListByChannelIDCalls int
}

func (r *offlineTimerStreamsRepository) GetByChannelID(
	context.Context,
	uuid.UUID,
	platform.Platform,
) (streamsmodel.Stream, error) {
	return streamsmodel.Stream{}, nil
}

func (r *offlineTimerStreamsRepository) GetByUserID(
	context.Context,
	string,
	platform.Platform,
) (streamsmodel.Stream, error) {
	return streamsmodel.Stream{}, nil
}

func (r *offlineTimerStreamsRepository) GetListByChannelID(
	context.Context,
	uuid.UUID,
) ([]streamsmodel.Stream, error) {
	r.getListByChannelIDCalls++
	return nil, nil
}

func (r *offlineTimerStreamsRepository) GetList(context.Context) ([]streamsmodel.Stream, error) {
	return nil, nil
}

func (r *offlineTimerStreamsRepository) Count(context.Context) (uint64, error) {
	return 0, nil
}

func (r *offlineTimerStreamsRepository) Save(context.Context, streamsrepository.SaveInput) error {
	return nil
}

func (r *offlineTimerStreamsRepository) DeleteByChannelID(
	context.Context,
	uuid.UUID,
	platform.Platform,
) error {
	return nil
}

func (r *offlineTimerStreamsRepository) Update(
	context.Context,
	uuid.UUID,
	platform.Platform,
	streamsrepository.UpdateInput,
) error {
	return nil
}
