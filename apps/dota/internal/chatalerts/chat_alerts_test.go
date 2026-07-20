package chatalerts

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/twirapp/kv"
	kvoptions "github.com/twirapp/kv/options"
	"github.com/twirapp/twir/libs/bus-core/bots"
	busdota "github.com/twirapp/twir/libs/bus-core/dota"
	dotarepository "github.com/twirapp/twir/libs/repositories/dota"
	dotamodel "github.com/twirapp/twir/libs/repositories/dota/model"
	"go.uber.org/fx"
)

type fakeValuer struct {
	err error
}

func (v fakeValuer) Int() (int64, error)     { return 0, v.err }
func (v fakeValuer) String() (string, error) { return "", v.err }
func (v fakeValuer) Bytes() ([]byte, error)  { return nil, v.err }
func (v fakeValuer) Bool() (bool, error)     { return false, v.err }
func (v fakeValuer) Float() (float64, error) { return 0, v.err }
func (v fakeValuer) Scan(any) error          { return v.err }
func (v fakeValuer) Err() error              { return v.err }

type fakeCache struct {
	values   map[string]struct{}
	getErr   error
	setErr   error
	getCalls []string
	setCalls []string
}

func newFakeCache() *fakeCache {
	return &fakeCache{values: make(map[string]struct{})}
}

func (c *fakeCache) Get(_ context.Context, key string) kv.Valuer {
	c.getCalls = append(c.getCalls, key)
	if c.getErr != nil {
		return fakeValuer{err: c.getErr}
	}
	if _, ok := c.values[key]; !ok {
		return fakeValuer{err: kv.ErrKeyNil}
	}
	return fakeValuer{}
}

func (c *fakeCache) Set(
	_ context.Context,
	key string,
	_ any,
	_ ...kvoptions.Option,
) error {
	c.setCalls = append(c.setCalls, key)
	if c.setErr != nil {
		return c.setErr
	}
	c.values[key] = struct{}{}
	return nil
}

func (c *fakeCache) SetMany(_ context.Context, _ []kv.SetMany) error { return nil }
func (c *fakeCache) Delete(_ context.Context, key string) error {
	delete(c.values, key)
	return nil
}
func (c *fakeCache) DeleteMany(_ context.Context, _ []string) error { return nil }
func (c *fakeCache) Exists(_ context.Context, key string) (bool, error) {
	_, ok := c.values[key]
	return ok, nil
}
func (c *fakeCache) ExistsMany(_ context.Context, keys []string) ([]bool, error) {
	values := make([]bool, len(keys))
	for i, key := range keys {
		values[i], _ = c.Exists(context.Background(), key)
	}
	return values, nil
}
func (c *fakeCache) GetKeysByPattern(_ context.Context, _ string) ([]string, error) {
	return nil, nil
}

type fakeSettingsRepository struct {
	settings dotamodel.ChannelDotaSettings
	err      error
	channel  []uuid.UUID
}

func (r *fakeSettingsRepository) GetByChannelID(
	_ context.Context,
	channelID uuid.UUID,
) (dotamodel.ChannelDotaSettings, error) {
	r.channel = append(r.channel, channelID)
	if r.err != nil {
		return dotamodel.Nil, r.err
	}
	return r.settings, nil
}

type fakeMessagePublisher struct {
	err      error
	requests []bots.SendMessageRequest
}

func (p *fakeMessagePublisher) Publish(
	_ context.Context,
	request bots.SendMessageRequest,
) error {
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
	alerts    *ChatAlerts
	cache     *fakeCache
	channelID uuid.UUID
	publisher *fakeMessagePublisher
	repo      *fakeSettingsRepository
}

func newFixture(t *testing.T, settings dotamodel.ChannelDotaSettings) *fixture {
	t.Helper()

	cache := newFakeCache()
	publisher := &fakeMessagePublisher{}
	repo := &fakeSettingsRepository{settings: settings}

	return &fixture{
		alerts:    newChatAlerts(repo, cache, publisher, testLogger(t), nil, nil),
		cache:     cache,
		channelID: uuid.New(),
		publisher: publisher,
		repo:      repo,
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
			require.Empty(t, f.cache.getCalls)
			require.Empty(t, f.cache.setCalls)
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
	require.Equal(t, []string{cooldownKey}, f.cache.setCalls)
}

func TestPublishFailureDoesNotSetCooldown(t *testing.T) {
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
	require.Empty(t, f.cache.setCalls)
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
	require.Empty(t, f.cache.getCalls)
	require.Empty(t, f.cache.setCalls)
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

func TestHandlerReturnsCooldownLookupError(t *testing.T) {
	settings := enabledSettings()
	settings.ChatEvents.MatchEnded = dotamodel.ChatEventSettings{
		Enabled:  true,
		Template: "match ended",
	}
	f := newFixture(t, settings)
	cacheErr := errors.New("cache unavailable")
	f.cache.getErr = cacheErr

	_, err := f.alerts.handleMatchEnded(
		context.Background(),
		busdota.MatchEndedMessage{ChannelID: f.channelID.String()},
	)

	require.ErrorIs(t, err, cacheErr)
	require.Empty(t, f.publisher.requests)
	require.Empty(t, f.cache.setCalls)
}

func TestLifecycleSubscribesAndUnsubscribesAllQueues(t *testing.T) {
	lifecycle := &fakeLifecycle{}
	subscriptions := []*fakeSubscription{{}, {}, {}, {}}
	newChatAlerts(
		&fakeSettingsRepository{},
		newFakeCache(),
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
		newFakeCache(),
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
