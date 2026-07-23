package bus_listener

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/bus-core/eventsub"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/gomodels"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

func TestBusListenerUsesNarrowDependencies(t *testing.T) {
	listenerType := reflect.TypeOf(BusListener{})
	for _, name := range []string{"eventSubClient", "kickSubManager", "channelService"} {
		field, found := listenerType.FieldByName(name)
		if !found || field.Type.Kind() != reflect.Interface {
			t.Fatalf("BusListener field %s is not an interface dependency", name)
		}
	}
}

type recordingEventSubManager struct {
	unsubscribedChannels []string
}

func (m *recordingEventSubManager) UnsubscribeChannel(_ context.Context, channelID string) error {
	m.unsubscribedChannels = append(m.unsubscribedChannels, channelID)
	return nil
}

func (*recordingEventSubManager) SubscribeToNeededEvents(context.Context, []model.EventsubTopic, string, string) error {
	return nil
}

func (*recordingEventSubManager) SubscribeToEvent(context.Context, string, string, string) error {
	return nil
}

type recordingKickSubscriptionManager struct {
	unsubscribed []channelplatformsmodel.ChannelPlatform
}

func (*recordingKickSubscriptionManager) Subscribe(context.Context, channelplatformsmodel.ChannelPlatform) error {
	return nil
}

func (m *recordingKickSubscriptionManager) Unsubscribe(_ context.Context, binding channelplatformsmodel.ChannelPlatform) error {
	m.unsubscribed = append(m.unsubscribed, binding)
	return nil
}

func TestUnsubscribeUsesTwitchSnapshotWithoutChannelLookup(t *testing.T) {
	twitch := &recordingEventSubManager{}
	listener := &BusListener{eventSubClient: twitch}
	_, err := listener.unsubscribe(context.Background(), eventsub.EventsubUnsubscribeRequest{
		ChannelID: uuid.NewString(),
		Platform:  platform.PlatformTwitch,
		Binding:   &eventsub.EventsubBindingSnapshot{PlatformChannelID: "twitch-channel"},
	})
	if err != nil {
		t.Fatalf("unsubscribe snapshot: %v", err)
	}
	if !reflect.DeepEqual(twitch.unsubscribedChannels, []string{"twitch-channel"}) {
		t.Fatalf("twitch unsubscribe calls = %#v", twitch.unsubscribedChannels)
	}
}

func TestUnsubscribeUsesKickSnapshotWithoutChannelLookup(t *testing.T) {
	bindingID := uuid.New()
	userID := uuid.New()
	kick := &recordingKickSubscriptionManager{}
	listener := &BusListener{kickSubManager: kick}
	_, err := listener.unsubscribe(context.Background(), eventsub.EventsubUnsubscribeRequest{
		ChannelID: uuid.NewString(),
		Platform:  platform.PlatformKick,
		Binding: &eventsub.EventsubBindingSnapshot{
			ID:                bindingID.String(),
			UserID:            userID.String(),
			PlatformChannelID: "kick-channel",
		},
	})
	if err != nil {
		t.Fatalf("unsubscribe snapshot: %v", err)
	}
	if len(kick.unsubscribed) != 1 || kick.unsubscribed[0].ID != bindingID || kick.unsubscribed[0].UserID != userID {
		t.Fatalf("kick unsubscribe bindings = %#v", kick.unsubscribed)
	}
}

type unsubscribeChannelReader struct {
	channel channelsmodel.Channel
	calls   int
}

func (r *unsubscribeChannelReader) GetChannelByID(_ context.Context, channelID uuid.UUID) (channelsmodel.Channel, error) {
	r.calls++
	if channelID != r.channel.ID {
		return channelsmodel.Nil, errors.New("unexpected channel ID")
	}
	return r.channel, nil
}

func TestUnsubscribeWithoutSnapshotLoadsChannel(t *testing.T) {
	channelID := uuid.New()
	twitch := &recordingEventSubManager{}
	reader := &unsubscribeChannelReader{channel: channelsmodel.Channel{
		ID: channelID,
		Bindings: []channelplatformsmodel.ChannelPlatform{{
			ID:                uuid.New(),
			ChannelID:         channelID,
			Platform:          platform.PlatformTwitch,
			PlatformChannelID: "twitch-channel",
		}},
	}}
	listener := &BusListener{eventSubClient: twitch, channelService: reader}
	_, err := listener.unsubscribe(context.Background(), eventsub.EventsubUnsubscribeRequest{
		ChannelID: channelID.String(),
		Platform:  platform.PlatformTwitch,
	})
	if err != nil {
		t.Fatalf("unsubscribe without snapshot: %v", err)
	}
	if reader.calls != 1 || !reflect.DeepEqual(twitch.unsubscribedChannels, []string{"twitch-channel"}) {
		t.Fatalf("reader calls = %d, twitch unsubscribe calls = %#v", reader.calls, twitch.unsubscribedChannels)
	}
}

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

func TestReinitBoundChannelsLimitsConcurrentCallbacks(t *testing.T) {
	const maxConcurrentCallbacks = 10
	const uniqueChannelCount = maxConcurrentCallbacks + 5

	twitchChannels := make([]channelsmodel.Channel, 0, uniqueChannelCount)
	kickChannels := make([]channelsmodel.Channel, 0, 1)
	wantCalls := make(map[uuid.UUID]struct{}, uniqueChannelCount)

	sharedID := uuid.New()
	sharedChannel := reinitTestChannel(sharedID, platform.PlatformTwitch, platform.PlatformKick)
	twitchChannels = append(twitchChannels, sharedChannel)
	kickChannels = append(kickChannels, sharedChannel)
	wantCalls[sharedID] = struct{}{}

	for range uniqueChannelCount - 1 {
		id := uuid.New()
		twitchChannels = append(twitchChannels, reinitTestChannel(id, platform.PlatformTwitch))
		wantCalls[id] = struct{}{}
	}

	repo := &reinitChannelsRepo{
		channelsByPlatform: map[platform.Platform][]channelsmodel.Channel{
			platform.PlatformTwitch: twitchChannels,
			platform.PlatformKick:   kickChannels,
		},
	}
	listener := &BusListener{channelsRepo: repo}

	started := make(chan struct{}, uniqueChannelCount)
	release := make(chan struct{})
	type reinitResult struct {
		count int
		err   error
	}
	done := make(chan reinitResult, 1)

	var mu sync.Mutex
	inFlight := 0
	maxInFlight := 0
	calls := make(map[uuid.UUID]int, uniqueChannelCount)
	go func() {
		count, err := listener.reinitBoundChannels(context.Background(), func(id uuid.UUID) {
			mu.Lock()
			inFlight++
			if inFlight > maxInFlight {
				maxInFlight = inFlight
			}
			calls[id]++
			mu.Unlock()

			started <- struct{}{}
			<-release

			mu.Lock()
			inFlight--
			mu.Unlock()
		})
		done <- reinitResult{count: count, err: err}
	}()

	var releaseOnce sync.Once
	releaseCallbacks := func() {
		releaseOnce.Do(func() { close(release) })
	}
	completed := false
	defer func() {
		releaseCallbacks()
		if !completed {
			select {
			case <-done:
			case <-time.After(time.Second):
				t.Error("reinitBoundChannels did not finish after callbacks were released")
			}
		}
	}()

	for range maxConcurrentCallbacks {
		select {
		case <-started:
		case <-time.After(time.Second):
			t.Fatalf("fewer than %d callbacks started", maxConcurrentCallbacks)
		}
	}
	select {
	case <-started:
		t.Fatalf("started more than %d callbacks before release", maxConcurrentCallbacks)
	case <-time.After(100 * time.Millisecond):
	}

	releaseCallbacks()
	result := <-done
	completed = true
	if result.err != nil {
		t.Fatalf("reinitBoundChannels returned error: %v", result.err)
	}
	if result.count != uniqueChannelCount {
		t.Fatalf("reinitialized channels = %d, want %d", result.count, uniqueChannelCount)
	}

	mu.Lock()
	defer mu.Unlock()
	if maxInFlight > maxConcurrentCallbacks {
		t.Fatalf("max in-flight callbacks = %d, want at most %d", maxInFlight, maxConcurrentCallbacks)
	}
	if len(calls) != uniqueChannelCount {
		t.Fatalf("reinitialized unique channels = %d, want %d", len(calls), uniqueChannelCount)
	}
	for id := range wantCalls {
		if calls[id] != 1 {
			t.Errorf("channel %s reinitialized %d times, want 1", id, calls[id])
		}
	}
}
