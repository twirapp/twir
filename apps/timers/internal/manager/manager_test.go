package manager

import (
	"context"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	"github.com/twirapp/twir/libs/entities/platform"
	timersentity "github.com/twirapp/twir/libs/entities/timers"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	"github.com/twirapp/twir/libs/repositories/timers"
)

// Regression test: channels.GetMany defaults PerPage to 10 when it is not set,
// so initialize must paginate explicitly, otherwise timers on channels beyond
// the first page are never scheduled until the timer is re-saved.
func TestInitializePaginatesChannels(t *testing.T) {
	targetChannelID := uuid.New()
	targetTimerID := uuid.New()

	channelsRepo := &paginatedChannelsRepository{
		pages: map[int][]channelentity.Channel{
			0: makeChannelsPage(100),
			1: {
				{
					ID: targetChannelID,
					Bindings: []channelplatformentity.ChannelPlatform{
						{
							Platform:          platform.PlatformTwitch,
							PlatformChannelID: "twitch-channel",
							Enabled:           true,
						},
					},
				},
			},
		},
	}

	manager := &Manager{
		timers: make(map[TimerID]*Timer),
		repository: &initializeTimersRepository{
			timers: []timersentity.Timer{
				{
					ID:        targetTimerID,
					ChannelID: targetChannelID,
					Enabled:   true,
				},
			},
		},
		logger:       slog.Default(),
		channelsRepo: channelsRepo,
	}

	if err := manager.initialize(context.Background()); err != nil {
		t.Fatalf("initialize() error = %v", err)
	}

	if channelsRepo.calls < 2 {
		t.Fatalf("channels GetMany calls = %d, want at least 2 (pagination)", channelsRepo.calls)
	}

	if _, ok := manager.timers[TimerID(targetTimerID)]; !ok {
		t.Fatalf("timer on channel beyond first channels page was not scheduled")
	}
}

func makeChannelsPage(count int) []channelentity.Channel {
	channels := make([]channelentity.Channel, 0, count)
	for range count {
		channels = append(channels, channelentity.Channel{ID: uuid.New()})
	}

	return channels
}

type paginatedChannelsRepository struct {
	pages map[int][]channelentity.Channel
	calls int
}

func (r *paginatedChannelsRepository) GetMany(
	_ context.Context,
	input channelsrepository.GetManyInput,
) ([]channelentity.Channel, error) {
	r.calls++
	if input.PerPage == 0 {
		panic("initialize must set PerPage explicitly: channels.GetMany defaults it to 10")
	}

	return r.pages[input.Page], nil
}

func (r *paginatedChannelsRepository) GetAllByBindingPlatform(
	context.Context,
	platform.Platform,
) ([]channelentity.Channel, error) {
	return nil, nil
}

func (r *paginatedChannelsRepository) GetByID(
	context.Context,
	uuid.UUID,
) (channelentity.Channel, error) {
	return channelentity.Nil, nil
}

func (r *paginatedChannelsRepository) GetByApiKey(
	context.Context,
	string,
) (channelentity.Channel, error) {
	return channelentity.Nil, nil
}

func (r *paginatedChannelsRepository) GetByBindingUserID(
	context.Context,
	platform.Platform,
	uuid.UUID,
) (channelentity.Channel, error) {
	return channelentity.Nil, nil
}

func (r *paginatedChannelsRepository) GetByPlatformChannelID(
	context.Context,
	platform.Platform,
	string,
) (channelentity.Channel, error) {
	return channelentity.Nil, nil
}

func (r *paginatedChannelsRepository) GetBySlug(
	context.Context,
	channelsrepository.GetBySlugInput,
) (channelentity.Channel, error) {
	return channelentity.Nil, nil
}

func (r *paginatedChannelsRepository) Update(
	context.Context,
	uuid.UUID,
	channelsrepository.UpdateInput,
) (channelentity.Channel, error) {
	return channelentity.Nil, nil
}

func (r *paginatedChannelsRepository) Create(context.Context) (channelentity.Channel, error) {
	return channelentity.Nil, nil
}

type initializeTimersRepository struct {
	timers []timersentity.Timer
}

func (r *initializeTimersRepository) GetByID(
	context.Context,
	uuid.UUID,
) (timersentity.Timer, error) {
	return timersentity.Nil, nil
}

func (r *initializeTimersRepository) GetAllByChannelID(
	context.Context,
	string,
) ([]timersentity.Timer, error) {
	return nil, nil
}

func (r *initializeTimersRepository) CountByChannelID(context.Context, string) (int, error) {
	return 0, nil
}

func (r *initializeTimersRepository) Create(
	context.Context,
	timers.CreateInput,
) (timersentity.Timer, error) {
	return timersentity.Nil, nil
}

func (r *initializeTimersRepository) UpdateByID(
	context.Context,
	uuid.UUID,
	timers.UpdateInput,
) (timersentity.Timer, error) {
	return timersentity.Nil, nil
}

func (r *initializeTimersRepository) Delete(context.Context, uuid.UUID) error {
	return nil
}

func (r *initializeTimersRepository) Count(context.Context, timers.CountInput) (int64, error) {
	return int64(len(r.timers)), nil
}

func (r *initializeTimersRepository) GetMany(
	context.Context,
	timers.GetManyInput,
) ([]timersentity.Timer, error) {
	return r.timers, nil
}
