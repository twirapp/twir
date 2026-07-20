package chatalerts

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/twirapp/twir/libs/bus-core/bots"
	busdota "github.com/twirapp/twir/libs/bus-core/dota"
	dotarepository "github.com/twirapp/twir/libs/repositories/dota"
	dotamodel "github.com/twirapp/twir/libs/repositories/dota/model"
	"go.uber.org/fx"
)

type fakeCooldownStore struct {
	mu sync.Mutex

	values map[string]struct{}

	reserveErr error
	releaseErr error

	reserveCalls []string
	releaseCalls []string

	reserveBarrier *sync.WaitGroup
}

func newFakeCooldownStore() *fakeCooldownStore {
	return &fakeCooldownStore{values: make(map[string]struct{})}
}

func (c *fakeCooldownStore) Reserve(
	_ context.Context,
	key string,
	_ time.Duration,
) (bool, error) {
	c.mu.Lock()
	c.reserveCalls = append(c.reserveCalls, key)
	err := c.reserveErr
	barrier := c.reserveBarrier
	c.mu.Unlock()

	if barrier != nil {
		barrier.Done()
		barrier.Wait()
	}
	if err != nil {
		return false, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.values[key]; exists {
		return false, nil
	}
	c.values[key] = struct{}{}
	return true, nil
}

func (c *fakeCooldownStore) Release(_ context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.releaseCalls = append(c.releaseCalls, key)
	if c.releaseErr != nil {
		return c.releaseErr
	}
	delete(c.values, key)
	return nil
}

func (c *fakeCooldownStore) operationCalls() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.reserveCalls) + len(c.releaseCalls)
}

type fakeSettingsRepository struct {
	mu       sync.Mutex
	settings dotamodel.ChannelDotaSettings
	err      error
	channel  []uuid.UUID
}

func (r *fakeSettingsRepository) GetByChannelID(
	_ context.Context,
	channelID uuid.UUID,
) (dotamodel.ChannelDotaSettings, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.channel = append(r.channel, channelID)
	if r.err != nil {
		return dotamodel.Nil, r.err
	}
	return r.settings, nil
}

type fakeMessagePublisher struct {
	mu       sync.Mutex
	err      error
	requests []bots.SendMessageRequest
}

func (p *fakeMessagePublisher) Publish(
	_ context.Context,
	request bots.SendMessageRequest,
) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.err != nil {
		return p.err
	}
	p.requests = append(p.requests, request)
	return nil
}

type fakeSubscription struct {
	err          error
	groups       []string
	unsubscribes int
}

func (s *fakeSubscription) Subscribe(group string) error {
	s.groups = append(s.groups, group)
	return s.err
}

func (s *fakeSubscription) Unsubscribe() {
	s.unsubscribes++
}

type fakeLifecycle struct {
	hooks []fx.Hook
}

func (l *fakeLifecycle) Append(hook fx.Hook) {
	l.hooks = append(l.hooks, hook)
}

type fixture struct {
	alerts        *ChatAlerts
	cooldownStore *fakeCooldownStore
	channelID     uuid.UUID
	publisher     *fakeMessagePublisher
	repo          *fakeSettingsRepository
}

func newFixture(t *testing.T, settings dotamodel.ChannelDotaSettings) *fixture {
	t.Helper()

	cooldownStore := newFakeCooldownStore()
	publisher := &fakeMessagePublisher{}
	repo := &fakeSettingsRepository{settings: settings}

	return &fixture{
		alerts:        newChatAlerts(repo, cooldownStore, publisher, testLogger(t), nil, nil),
		cooldownStore: cooldownStore,
		channelID:     uuid.New(),
		publisher:     publisher,
		repo:          repo,
	}
}

func testLogger(t *testing.T) *slog.Logger {
	t.Helper()
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func enabledSettings() dotamodel.ChannelDotaSettings {
	return dotamodel.ChannelDotaSettings{
		Enabled:       true,
		Mmr:           3600,
		SessionWins:   8,
		SessionLosses: 5,
	}
}

func TestMatchStartedSkipsDisabledSettings(t *testing.T) {
	tests := []struct {
		name     string
		settings func() dotamodel.ChannelDotaSettings
	}{
		{
			name: "module disabled",
			settings: func() dotamodel.ChannelDotaSettings {
				settings := enabledSettings()
				settings.Enabled = false
				settings.ChatEvents.MatchStarted = dotamodel.ChatEventSettings{
					Enabled:  true,
					Template: "{hero}",
				}
				return settings
			},
		},
		{
			name: "event disabled",
			settings: func() dotamodel.ChannelDotaSettings {
				settings := enabledSettings()
				settings.ChatEvents.MatchStarted = dotamodel.ChatEventSettings{
					Enabled:  false,
					Template: "{hero}",
				}
				return settings
			},
		},
		{
			name: "empty template",
			settings: func() dotamodel.ChannelDotaSettings {
				settings := enabledSettings()
				settings.ChatEvents.MatchStarted = dotamodel.ChatEventSettings{
					Enabled:  true,
					Template: " \t ",
				}
				return settings
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFixture(t, tt.settings())

			_, err := f.alerts.handleMatchStarted(
				context.Background(),
				busdota.MatchStartedMessage{
					ChannelID: f.channelID.String(),
					HeroName:  "Juggernaut",
				},
			)

			require.NoError(t, err)
			require.Empty(t, f.publisher.requests)
			require.Zero(t, f.cooldownStore.operationCalls())
		})
	}
}

func TestMatchStartedUsesSettingsStats(t *testing.T) {
	settings := enabledSettings()
	settings.ChatEvents.MatchStarted = dotamodel.ChatEventSettings{
		Enabled:  true,
		Template: "{hero}|{mmr}|{wins}|{losses}|{team}|{player}|{time}",
	}
	f := newFixture(t, settings)

	_, err := f.alerts.handleMatchStarted(
		context.Background(),
		busdota.MatchStartedMessage{
			ChannelID: f.channelID.String(),
			HeroName:  "Juggernaut",
		},
	)

	require.NoError(t, err)
	require.Equal(t, "Juggernaut|3600|8|5|||", f.publisher.requests[0].Message)
}

func TestMatchEndedRendersAndPublishes(t *testing.T) {
	settings := enabledSettings()
	settings.ChatEvents.MatchEnded = dotamodel.ChatEventSettings{
		Enabled:  true,
		Template: "{hero}|{mmr}|{wins}|{losses}|{team}|{player}|{time}",
	}
	f := newFixture(t, settings)

	_, err := f.alerts.handleMatchEnded(
		context.Background(),
		busdota.MatchEndedMessage{
			ChannelID:     f.channelID.String(),
			HeroName:      "Earthshaker",
			Mmr:           4210,
			SessionWins:   9,
			SessionLosses: 4,
		},
	)

	require.NoError(t, err)
	require.Equal(
		t,
		bots.SendMessageRequest{
			InternalChannelID: &f.channelID,
			Platform:          "twitch",
			Message:           "Earthshaker|4210|9|4|||",
			SkipRateLimits:    true,
		},
		f.publisher.requests[0],
	)
}

func TestRoshanAndAegisRenderEventDetails(t *testing.T) {
	tests := []struct {
		name      string
		configure func(*dotamodel.ChannelDotaSettings)
		deliver   func(*ChatAlerts, uuid.UUID) error
		expected  string
	}{
		{
			name: "roshan",
			configure: func(settings *dotamodel.ChannelDotaSettings) {
				settings.ChatEvents.RoshanKilled = dotamodel.ChatEventSettings{
					Enabled:  true,
					Template: "{team}|{player}|{time}|{hero}|{mmr}|{wins}|{losses}",
				}
			},
			deliver: func(alerts *ChatAlerts, channelID uuid.UUID) error {
				_, err := alerts.handleRoshanKilled(
					context.Background(),
					busdota.RoshanKilledMessage{
						ChannelID: channelID.String(),
						Team:      "radiant",
						GameTime:  125,
					},
				)
				return err
			},
			expected: "radiant||2:05||3600|8|5",
		},
		{
			name: "aegis",
			configure: func(settings *dotamodel.ChannelDotaSettings) {
				settings.ChatEvents.AegisPickup = dotamodel.ChatEventSettings{
					Enabled:  true,
					Template: "{team}|{player}|{time}|{hero}|{mmr}|{wins}|{losses}",
				}
			},
			deliver: func(alerts *ChatAlerts, channelID uuid.UUID) error {
				_, err := alerts.handleAegisPickup(
					context.Background(),
					busdota.AegisPickupMessage{
						ChannelID:  channelID.String(),
						PlayerName: "Puppey",
						GameTime:   125,
					},
				)
				return err
			},
			expected: "|Puppey|2:05||3600|8|5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			settings := enabledSettings()
			tt.configure(&settings)
			f := newFixture(t, settings)

			require.NoError(t, tt.deliver(f.alerts, f.channelID))
			require.Equal(t, tt.expected, f.publisher.requests[0].Message)
		})
	}
}

func TestAegisPlayerTemplate(t *testing.T) {
	playerID := 2
	tests := []struct {
		name     string
		message  busdota.AegisPickupMessage
		expected string
	}{
		{
			name:     "uses player ID when name is unavailable",
			message:  busdota.AegisPickupMessage{PlayerID: &playerID},
			expected: "Aegis: player #2",
		},
		{
			name:     "prefers player name",
			message:  busdota.AegisPickupMessage{PlayerID: &playerID, PlayerName: "Puppey"},
			expected: "Aegis: Puppey",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			settings := enabledSettings()
			settings.ChatEvents.AegisPickup = dotamodel.ChatEventSettings{
				Enabled:  true,
				Template: "Aegis: {player}",
			}
			f := newFixture(t, settings)

			message := tt.message
			message.ChannelID = f.channelID.String()

			_, err := f.alerts.handleAegisPickup(context.Background(), message)
			require.NoError(t, err)
			require.Equal(t, tt.expected, f.publisher.requests[0].Message)
		})
	}
}

func TestCooldownSuppressesDuplicateEvent(t *testing.T) {
	settings := enabledSettings()
	settings.ChatEvents.MatchEnded = dotamodel.ChatEventSettings{
		Enabled:  true,
		Template: "match ended",
		Cooldown: 30,
	}
	f := newFixture(t, settings)
	event := busdota.MatchEndedMessage{ChannelID: f.channelID.String()}

	_, err := f.alerts.handleMatchEnded(context.Background(), event)
	require.NoError(t, err)
	_, err = f.alerts.handleMatchEnded(context.Background(), event)
	require.NoError(t, err)

	cooldownKey := "cache:twir:dota:chat-alert:" + f.channelID.String() + ":match_ended"
	require.Len(t, f.publisher.requests, 1)
	require.Equal(t, []string{cooldownKey, cooldownKey}, f.cooldownStore.reserveCalls)
}

func TestCooldownZeroBypassesStore(t *testing.T) {
	settings := enabledSettings()
	settings.ChatEvents.MatchEnded = dotamodel.ChatEventSettings{
		Enabled:  true,
		Template: "match ended",
	}
	f := newFixture(t, settings)

	_, err := f.alerts.handleMatchEnded(
		context.Background(),
		busdota.MatchEndedMessage{ChannelID: f.channelID.String()},
	)

	require.NoError(t, err)
	require.Len(t, f.publisher.requests, 1)
	require.Zero(t, f.cooldownStore.operationCalls())
}

func TestBlankRenderedMessageBypassesStoreAndPublish(t *testing.T) {
	settings := enabledSettings()
	settings.ChatEvents.AegisPickup = dotamodel.ChatEventSettings{
		Enabled:  true,
		Template: "{player}",
		Cooldown: 30,
	}
	f := newFixture(t, settings)

	_, err := f.alerts.handleAegisPickup(
		context.Background(),
		busdota.AegisPickupMessage{ChannelID: f.channelID.String()},
	)

	require.NoError(t, err)
	require.Empty(t, f.publisher.requests)
	require.Zero(t, f.cooldownStore.operationCalls())
}

func TestConcurrentCooldownReservationsPublishOnce(t *testing.T) {
	settings := enabledSettings()
	settings.ChatEvents.MatchEnded = dotamodel.ChatEventSettings{
		Enabled:  true,
		Template: "match ended",
		Cooldown: 30,
	}
	f := newFixture(t, settings)

	var barrier sync.WaitGroup
	barrier.Add(2)
	f.cooldownStore.reserveBarrier = &barrier

	event := busdota.MatchEndedMessage{ChannelID: f.channelID.String()}
	start := make(chan struct{})
	errs := make(chan error, 2)
	for range 2 {
		go func() {
			<-start
			_, err := f.alerts.handleMatchEnded(context.Background(), event)
			errs <- err
		}()
	}
	close(start)

	require.NoError(t, <-errs)
	require.NoError(t, <-errs)
	require.Len(t, f.publisher.requests, 1)
	require.Len(t, f.cooldownStore.reserveCalls, 2)
}

func TestPublishFailureReleasesReservation(t *testing.T) {
	settings := enabledSettings()
	settings.ChatEvents.MatchEnded = dotamodel.ChatEventSettings{
		Enabled:  true,
		Template: "match ended",
		Cooldown: 30,
	}
	f := newFixture(t, settings)
	publishErr := errors.New("bots unavailable")
	f.publisher.err = publishErr

	_, err := f.alerts.handleMatchEnded(
		context.Background(),
		busdota.MatchEndedMessage{ChannelID: f.channelID.String()},
	)

	require.ErrorIs(t, err, publishErr)
	require.Len(t, f.cooldownStore.releaseCalls, 1)
}

func TestPublishFailurePreservesErrorWhenReleaseFails(t *testing.T) {
	settings := enabledSettings()
	settings.ChatEvents.MatchEnded = dotamodel.ChatEventSettings{
		Enabled:  true,
		Template: "match ended",
		Cooldown: 30,
	}
	f := newFixture(t, settings)
	publishErr := errors.New("bots unavailable")
	f.publisher.err = publishErr
	f.cooldownStore.releaseErr = errors.New("cache unavailable")

	_, err := f.alerts.handleMatchEnded(
		context.Background(),
		busdota.MatchEndedMessage{ChannelID: f.channelID.String()},
	)

	require.ErrorIs(t, err, publishErr)
	require.Len(t, f.cooldownStore.releaseCalls, 1)
}

func TestHandlerIgnoresMissingSettings(t *testing.T) {
	f := newFixture(t, enabledSettings())
	f.repo.err = dotarepository.ErrNotFound

	_, err := f.alerts.handleMatchEnded(
		context.Background(),
		busdota.MatchEndedMessage{ChannelID: f.channelID.String()},
	)

	require.NoError(t, err)
	require.Empty(t, f.publisher.requests)
	require.Zero(t, f.cooldownStore.operationCalls())
}

func TestHandlerReturnsInvalidChannelAndRepositoryErrors(t *testing.T) {
	t.Run("invalid channel ID", func(t *testing.T) {
		f := newFixture(t, enabledSettings())

		_, err := f.alerts.handleMatchEnded(
			context.Background(),
			busdota.MatchEndedMessage{ChannelID: "not-a-uuid"},
		)

		require.Error(t, err)
		require.Empty(t, f.repo.channel)
	})

	t.Run("repository error", func(t *testing.T) {
		f := newFixture(t, enabledSettings())
		repositoryErr := errors.New("database unavailable")
		f.repo.err = repositoryErr

		_, err := f.alerts.handleMatchEnded(
			context.Background(),
			busdota.MatchEndedMessage{ChannelID: f.channelID.String()},
		)

		require.ErrorIs(t, err, repositoryErr)
	})
}

func TestHandlerReturnsCooldownReservationError(t *testing.T) {
	settings := enabledSettings()
	settings.ChatEvents.MatchEnded = dotamodel.ChatEventSettings{
		Enabled:  true,
		Template: "match ended",
		Cooldown: 30,
	}
	f := newFixture(t, settings)
	reserveErr := errors.New("cache unavailable")
	f.cooldownStore.reserveErr = reserveErr

	_, err := f.alerts.handleMatchEnded(
		context.Background(),
		busdota.MatchEndedMessage{ChannelID: f.channelID.String()},
	)

	require.ErrorIs(t, err, reserveErr)
	require.Empty(t, f.publisher.requests)
	require.Empty(t, f.cooldownStore.releaseCalls)
}

func TestLifecycleSubscribesAndUnsubscribesAllQueues(t *testing.T) {
	lifecycle := &fakeLifecycle{}
	subscriptions := []*fakeSubscription{{}, {}, {}, {}}
	newChatAlerts(
		&fakeSettingsRepository{},
		newFakeCooldownStore(),
		&fakeMessagePublisher{},
		testLogger(t),
		[]subscription{subscriptions[0], subscriptions[1], subscriptions[2], subscriptions[3]},
		lifecycle,
	)

	require.Len(t, lifecycle.hooks, 1)
	require.NotNil(t, lifecycle.hooks[0].OnStart)
	require.NotNil(t, lifecycle.hooks[0].OnStop)
	require.NoError(t, lifecycle.hooks[0].OnStart(context.Background()))

	for _, subscription := range subscriptions {
		require.Equal(t, []string{"dota"}, subscription.groups)
	}

	require.NoError(t, lifecycle.hooks[0].OnStop(context.Background()))
	for _, subscription := range subscriptions {
		require.Equal(t, 1, subscription.unsubscribes)
	}
}

func TestLifecycleCleansUpAfterPartialSubscriptionFailure(t *testing.T) {
	lifecycle := &fakeLifecycle{}
	subscriptionErr := errors.New("subscription failed")
	subscriptions := []*fakeSubscription{{}, {}, {err: subscriptionErr}, {}}
	newChatAlerts(
		&fakeSettingsRepository{},
		newFakeCooldownStore(),
		&fakeMessagePublisher{},
		testLogger(t),
		[]subscription{subscriptions[0], subscriptions[1], subscriptions[2], subscriptions[3]},
		lifecycle,
	)

	err := lifecycle.hooks[0].OnStart(context.Background())

	require.ErrorIs(t, err, subscriptionErr)
	require.Equal(t, []string{"dota"}, subscriptions[0].groups)
	require.Equal(t, []string{"dota"}, subscriptions[1].groups)
	require.Equal(t, []string{"dota"}, subscriptions[2].groups)
	require.Empty(t, subscriptions[3].groups)
	require.Equal(t, 1, subscriptions[0].unsubscribes)
	require.Equal(t, 1, subscriptions[1].unsubscribes)
	require.Zero(t, subscriptions[2].unsubscribes)
	require.Zero(t, subscriptions[3].unsubscribes)
}
