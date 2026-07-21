package bus_listener

import (
	"context"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

type reinitChannelsRepo struct {
	channelsByPlatform map[platform.Platform][]channelsmodel.Channel
	platformLookups    []platform.Platform
	getManyCalls       int
}

func (r *reinitChannelsRepo) GetAllByBindingPlatform(
	_ context.Context,
	p platform.Platform,
) ([]channelsmodel.Channel, error) {
	r.platformLookups = append(r.platformLookups, p)
	return r.channelsByPlatform[p], nil
}

func (r *reinitChannelsRepo) GetMany(
	context.Context,
	channelsrepo.GetManyInput,
) ([]channelsmodel.Channel, error) {
	r.getManyCalls++
	return nil, nil
}

func (r *reinitChannelsRepo) GetByID(context.Context, uuid.UUID) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (r *reinitChannelsRepo) GetByApiKey(context.Context, string) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (r *reinitChannelsRepo) GetByBindingUserID(
	context.Context,
	platform.Platform,
	uuid.UUID,
) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (r *reinitChannelsRepo) GetByPlatformChannelID(
	context.Context,
	platform.Platform,
	string,
) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (r *reinitChannelsRepo) GetBySlug(
	context.Context,
	channelsrepo.GetBySlugInput,
) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (r *reinitChannelsRepo) GetCount(context.Context, channelsrepo.GetCountInput) (int, error) {
	return 0, nil
}

func (r *reinitChannelsRepo) Update(
	context.Context,
	uuid.UUID,
	channelsrepo.UpdateInput,
) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (r *reinitChannelsRepo) Create(
	context.Context,
	channelsrepo.CreateInput,
) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func reinitTestChannel(id uuid.UUID, platforms ...platform.Platform) channelsmodel.Channel {
	bindings := make([]channelplatformsmodel.ChannelPlatform, 0, len(platforms))
	for _, p := range platforms {
		bindings = append(bindings, channelplatformsmodel.ChannelPlatform{Platform: p})
	}

	return channelsmodel.Channel{
		ID:       id,
		Bindings: bindings,
	}
}

func TestReinitBoundChannelsProcessesEveryUniquePlatformBinding(t *testing.T) {
	const twitchOnlyCount = 11
	const kickOnlyCount = 11

	twitchChannels := make([]channelsmodel.Channel, 0, twitchOnlyCount+1)
	kickChannels := make([]channelsmodel.Channel, 0, kickOnlyCount+1)
	wantCalls := make(map[uuid.UUID]struct{}, twitchOnlyCount+kickOnlyCount+1)

	sharedID := uuid.New()
	sharedChannel := reinitTestChannel(sharedID, platform.PlatformTwitch, platform.PlatformKick)
	twitchChannels = append(twitchChannels, sharedChannel)
	kickChannels = append(kickChannels, sharedChannel)
	wantCalls[sharedID] = struct{}{}

	for range twitchOnlyCount {
		id := uuid.New()
		twitchChannels = append(twitchChannels, reinitTestChannel(id, platform.PlatformTwitch))
		wantCalls[id] = struct{}{}
	}
	for range kickOnlyCount {
		id := uuid.New()
		kickChannels = append(kickChannels, reinitTestChannel(id, platform.PlatformKick))
		wantCalls[id] = struct{}{}
	}

	repo := &reinitChannelsRepo{
		channelsByPlatform: map[platform.Platform][]channelsmodel.Channel{
			platform.PlatformTwitch: twitchChannels,
			platform.PlatformKick:   kickChannels,
		},
	}
	listener := &BusListener{channelsRepo: repo}

	var mu sync.Mutex
	calls := make(map[uuid.UUID]int, len(wantCalls))
	count, err := listener.reinitBoundChannels(context.Background(), func(id uuid.UUID) {
		mu.Lock()
		defer mu.Unlock()
		calls[id]++
	})
	if err != nil {
		t.Fatalf("reinitBoundChannels returned error: %v", err)
	}

	if count != len(wantCalls) {
		t.Fatalf("reinitialized channels = %d, want %d", count, len(wantCalls))
	}
	if repo.getManyCalls != 0 {
		t.Fatalf("GetMany calls = %d, want 0", repo.getManyCalls)
	}
	if len(repo.platformLookups) != 2 ||
		repo.platformLookups[0] != platform.PlatformTwitch ||
		repo.platformLookups[1] != platform.PlatformKick {
		t.Fatalf("binding platform lookups = %v, want [%s %s]", repo.platformLookups, platform.PlatformTwitch, platform.PlatformKick)
	}
	if len(calls) != len(wantCalls) {
		t.Fatalf("reinitialized unique channels = %d, want %d", len(calls), len(wantCalls))
	}
	for id := range wantCalls {
		if calls[id] != 1 {
			t.Errorf("channel %s reinitialized %d times, want 1", id, calls[id])
		}
	}
}
