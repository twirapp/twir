package chatalerts

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/twirapp/twir/libs/bus-core/bots"
	busdota "github.com/twirapp/twir/libs/bus-core/dota"
	dotarepository "github.com/twirapp/twir/libs/repositories/dota"
	dotamodel "github.com/twirapp/twir/libs/repositories/dota/model"
	"go.uber.org/fx"
)

type fakeCooldownStore struct {
	mu sync.Mutex

	values map[string]string

	reserveErr error
	releaseErr error

	reserveCalls       []string
	releaseCalls       []string
	releaseCtxErr      error
	releaseDeadline    time.Time
	releaseHasDeadline bool

	reserveBarrier *sync.WaitGroup
}

func newFakeCooldownStore() *fakeCooldownStore {
	return &fakeCooldownStore{values: make(map[string]string)}
}

func (c *fakeCooldownStore) Reserve(
	_ context.Context,
	key string,
	token string,
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
	c.values[key] = token
	return true, nil
}

func (c *fakeCooldownStore) Release(ctx context.Context, key string, token string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.releaseCalls = append(c.releaseCalls, key)
	c.releaseCtxErr = ctx.Err()
	c.releaseDeadline, c.releaseHasDeadline = ctx.Deadline()
	if c.releaseErr != nil {
		return c.releaseErr
	}
	if c.values[key] == token {
		delete(c.values, key)
	}
	return nil
}

func (c *fakeCooldownStore) operationCalls() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.reserveCalls) + len(c.releaseCalls)
}

func (c *fakeCooldownStore) releaseContext() (error, time.Time, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.releaseCtxErr, c.releaseDeadline, c.releaseHasDeadline
}

type fakeRedisCooldownClient struct {
	values map[string]string

	evalScript string
	evalKeys   []string
	evalArgs   []any
}

func (c *fakeRedisCooldownClient) SetNX(
	_ context.Context,
	key string,
	value interface{},
	_ time.Duration,
) *redis.BoolCmd {
	if _, exists := c.values[key]; exists {
		return redis.NewBoolResult(false, nil)
	}

	c.values[key] = value.(string)
	return redis.NewBoolResult(true, nil)
}

func (c *fakeRedisCooldownClient) Eval(
	_ context.Context,
	script string,
	keys []string,
	args ...interface{},
) *redis.Cmd {
	c.evalScript = script
	c.evalKeys = append([]string(nil), keys...)
	c.evalArgs = append([]any(nil), args...)

	if len(keys) != 1 || len(args) != 1 {
		return redis.NewCmdResult(nil, errors.New("unexpected Redis Eval arguments"))
	}

	if c.values[keys[0]] == args[0] {
		delete(c.values, keys[0])
		return redis.NewCmdResult(int64(1), nil)
	}

	return redis.NewCmdResult(int64(0), nil)
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

type blockingFailedPublisher struct {
	started chan struct{}
	release chan struct{}
	err     error
}

func (p *blockingFailedPublisher) Publish(
	_ context.Context,
	_ bots.SendMessageRequest,
) error {
	close(p.started)
	<-p.release
	return p.err
}

type expiringCooldownStore struct {
	mu sync.Mutex

	key   string
	value string
}

func newExpiringCooldownStore() *expiringCooldownStore {
	return &expiringCooldownStore{}
}

func (c *expiringCooldownStore) Reserve(
	_ context.Context,
	key string,
	token string,
	_ time.Duration,
) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.value != "" {
		return false, nil
	}

	c.key = key
	c.value = token
	return true, nil
}

func (c *expiringCooldownStore) Release(_ context.Context, key string, token string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.key == key && c.value == token {
		c.value = ""
	}
	return nil
}

func (c *expiringCooldownStore) Expire(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.key == key {
		c.value = ""
	}
}

func (c *expiringCooldownStore) Value() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
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

type notifyingSubscription struct {
	unsubscribed chan struct{}
}

func (s *notifyingSubscription) Subscribe(_ string) error {
	return nil
}

func (s *notifyingSubscription) Unsubscribe() {
	close(s.unsubscribed)
}

type fakeLifecycle struct {
	hooks []fx.Hook
}

func (l *fakeLifecycle) Append(hook fx.Hook) {
	l.hooks = append(l.hooks, hook)
}

type capturedLogRecord struct {
	message string
	attrs   map[string]string
}

type captureHandler struct {
	mu      sync.Mutex
	records []capturedLogRecord
}

func (*captureHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *captureHandler) Handle(_ context.Context, record slog.Record) error {
	attrs := make(map[string]string)
	record.Attrs(func(attr slog.Attr) bool {
		attrs[attr.Key] = attr.Value.String()
		return true
	})

	h.mu.Lock()
	defer h.mu.Unlock()
	h.records = append(h.records, capturedLogRecord{
		message: record.Message,
		attrs:   attrs,
	})
	return nil
}

func (h *captureHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *captureHandler) WithGroup(_ string) slog.Handler {
	return h
}

func (h *captureHandler) Records() []capturedLogRecord {
	h.mu.Lock()
	defer h.mu.Unlock()

	records := make([]capturedLogRecord, len(h.records))
	for i, record := range h.records {
		attrs := make(map[string]string, len(record.attrs))
		for key, value := range record.attrs {
			attrs[key] = value
		}
		records[i] = capturedLogRecord{message: record.message, attrs: attrs}
	}
	return records
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

func TestOverflowCooldownReturnsBeforeReserveAndPublish(t *testing.T) {
	if int64(math.MaxInt) <= math.MaxInt64/int64(time.Second) {
		t.Skip("platform max int cannot overflow a duration in seconds")
	}

	settings := enabledSettings()
	settings.ChatEvents.MatchEnded = dotamodel.ChatEventSettings{
		Enabled:  true,
		Template: "match ended",
		Cooldown: math.MaxInt,
	}
	f := newFixture(t, settings)

	_, err := f.alerts.handleMatchEnded(
		context.Background(),
		busdota.MatchEndedMessage{ChannelID: f.channelID.String()},
	)

	require.Error(t, err)
	require.Empty(t, f.publisher.requests)
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

func TestPublishFailureDoesNotReleaseReacquiredCooldown(t *testing.T) {
	channelID := uuid.New()
	settings := enabledSettings()
	settings.ChatEvents.MatchEnded = dotamodel.ChatEventSettings{
		Enabled:  true,
		Template: "match ended",
		Cooldown: 30,
	}
	publishErr := errors.New("publish failed")
	publisher := &blockingFailedPublisher{
		started: make(chan struct{}),
		release: make(chan struct{}),
		err:     publishErr,
	}
	store := newExpiringCooldownStore()
	alerts := newChatAlerts(
		&fakeSettingsRepository{settings: settings},
		store,
		publisher,
		testLogger(t),
		nil,
		nil,
	)

	result := make(chan error, 1)
	go func() {
		_, err := alerts.handleMatchEnded(
			context.Background(),
			busdota.MatchEndedMessage{ChannelID: channelID.String()},
		)
		result <- err
	}()

	<-publisher.started
	cooldownKey := fmt.Sprintf("%s:%s:%s", cooldownKeyPrefix, channelID, eventMatchEnded)
	store.Expire(cooldownKey)
	reacquired, err := store.Reserve(context.Background(), cooldownKey, "token-b", time.Second)
	require.NoError(t, err)
	require.True(t, reacquired)
	close(publisher.release)

	require.ErrorIs(t, <-result, publishErr)
	require.Equal(t, "token-b", store.Value())
}

func TestRedisCooldownStoreReleaseDoesNotDeleteReplacementToken(t *testing.T) {
	const (
		key   = "cache:twir:dota:chat-alert:channel:match_ended"
		token = "token-a"
	)
	client := &fakeRedisCooldownClient{values: map[string]string{key: "token-b"}}
	store := &RedisCooldownStore{client: client}

	err := store.Release(context.Background(), key, token)

	require.NoError(t, err)
	require.Equal(t, "token-b", client.values[key])
	require.Equal(
		t,
		`if redis.call("GET", KEYS[1]) == ARGV[1] then return redis.call("DEL", KEYS[1]) end return 0`,
		client.evalScript,
	)
	require.Equal(t, []string{key}, client.evalKeys)
	require.Equal(t, []any{token}, client.evalArgs)
}

func TestPublishFailureUsesBoundedDetachedCleanupContext(t *testing.T) {
	settings := enabledSettings()
	settings.ChatEvents.MatchEnded = dotamodel.ChatEventSettings{
		Enabled:  true,
		Template: "match ended",
		Cooldown: 30,
	}
	f := newFixture(t, settings)
	publishErr := errors.New("bots unavailable")
	f.publisher.err = publishErr
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	before := time.Now()

	_, err := f.alerts.handleMatchEnded(
		ctx,
		busdota.MatchEndedMessage{ChannelID: f.channelID.String()},
	)

	require.ErrorIs(t, err, publishErr)
	releaseCtxErr, deadline, hasDeadline := f.cooldownStore.releaseContext()
	require.NoError(t, releaseCtxErr)
	require.True(t, hasDeadline)
	require.WithinDuration(t, before.Add(2*time.Second), deadline, time.Second)
}

func TestCallbackLogsGenericTerminalErrorWithoutSensitiveData(t *testing.T) {
	channelID := uuid.New()
	template := "private template"
	renderedMessage := "private rendered message"
	reserveErr := errors.New("reserve failed for " + channelID.String() + " " + template + " " + renderedMessage)
	settings := enabledSettings()
	settings.ChatEvents.MatchEnded = dotamodel.ChatEventSettings{
		Enabled:  true,
		Template: template,
		Cooldown: 30,
	}
	store := newFakeCooldownStore()
	store.reserveErr = reserveErr
	logs := &captureHandler{}
	alerts := newChatAlerts(
		&fakeSettingsRepository{settings: settings},
		store,
		&fakeMessagePublisher{},
		slog.New(logs),
		nil,
		nil,
	)

	_, err := alerts.handleMatchEnded(
		context.Background(),
		busdota.MatchEndedMessage{ChannelID: channelID.String()},
	)

	require.ErrorIs(t, err, reserveErr)
	require.Equal(t, "dota chat alert handler failed: match_ended", err.Error())
	require.NotContains(t, err.Error(), channelID.String())
	require.NotContains(t, err.Error(), template)
	require.NotContains(t, err.Error(), renderedMessage)
	records := logs.Records()
	require.Len(t, records, 1)
	require.Equal(t, "dota chat alert handler failed", records[0].message)
	require.Equal(t, map[string]string{"event_kind": "match_ended"}, records[0].attrs)
	require.NotContains(t, fmt.Sprint(records), channelID.String())
	require.NotContains(t, fmt.Sprint(records), template)
	require.NotContains(t, fmt.Sprint(records), renderedMessage)
}

func TestPublishFailureLogsGenericCleanupFailure(t *testing.T) {
	channelID := uuid.New()
	template := "private template"
	settings := enabledSettings()
	settings.ChatEvents.MatchEnded = dotamodel.ChatEventSettings{
		Enabled:  true,
		Template: template,
		Cooldown: 30,
	}
	publishErr := errors.New("publish failed for " + channelID.String() + " " + template)
	releaseErr := errors.New("release failed for " + channelID.String() + " " + template)
	store := newFakeCooldownStore()
	store.releaseErr = releaseErr
	logs := &captureHandler{}
	alerts := newChatAlerts(
		&fakeSettingsRepository{settings: settings},
		store,
		&fakeMessagePublisher{err: publishErr},
		slog.New(logs),
		nil,
		nil,
	)

	_, err := alerts.handleMatchEnded(
		context.Background(),
		busdota.MatchEndedMessage{ChannelID: channelID.String()},
	)

	require.ErrorIs(t, err, publishErr)
	records := logs.Records()
	require.GreaterOrEqual(t, len(records), 1)
	require.Equal(t, "dota chat alert: cooldown cleanup failed", records[0].message)
	require.Empty(t, records[0].attrs)
	require.NotContains(t, records[0].message, channelID.String())
	require.NotContains(t, records[0].message, template)
	require.NotContains(t, records[0].message, releaseErr.Error())
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

func TestLifecycleStopWaitsForActiveHandlerAfterUnsubscribe(t *testing.T) {
	settings := enabledSettings()
	settings.ChatEvents.MatchEnded = dotamodel.ChatEventSettings{
		Enabled:  true,
		Template: "match ended",
	}
	lifecycle := &fakeLifecycle{}
	publisher := &blockingFailedPublisher{
		started: make(chan struct{}),
		release: make(chan struct{}),
	}
	sub := &notifyingSubscription{unsubscribed: make(chan struct{})}
	alerts := newChatAlerts(
		&fakeSettingsRepository{settings: settings},
		newFakeCooldownStore(),
		publisher,
		testLogger(t),
		[]subscription{sub},
		lifecycle,
	)
	released := false
	defer func() {
		if !released {
			close(publisher.release)
		}
	}()

	require.NoError(t, lifecycle.hooks[0].OnStart(context.Background()))
	handlerResult := make(chan error, 1)
	go func() {
		_, err := alerts.handleMatchEnded(
			context.Background(),
			busdota.MatchEndedMessage{ChannelID: uuid.NewString()},
		)
		handlerResult <- err
	}()
	<-publisher.started

	stopResult := make(chan error, 1)
	go func() {
		stopResult <- lifecycle.hooks[0].OnStop(context.Background())
	}()
	<-sub.unsubscribed
	select {
	case err := <-stopResult:
		t.Fatalf("OnStop returned while the active handler was blocked: %v", err)
	default:
	}

	close(publisher.release)
	released = true
	require.NoError(t, <-handlerResult)
	require.NoError(t, <-stopResult)
}

func TestLifecycleStopRejectsCallbacksBeforeDependencies(t *testing.T) {
	settings := enabledSettings()
	settings.ChatEvents.MatchEnded = dotamodel.ChatEventSettings{
		Enabled:  true,
		Template: "match ended",
		Cooldown: 30,
	}
	lifecycle := &fakeLifecycle{}
	store := newFakeCooldownStore()
	publisher := &fakeMessagePublisher{}
	repository := &fakeSettingsRepository{settings: settings}
	alerts := newChatAlerts(
		repository,
		store,
		publisher,
		testLogger(t),
		[]subscription{&notifyingSubscription{unsubscribed: make(chan struct{})}},
		lifecycle,
	)

	require.NoError(t, lifecycle.hooks[0].OnStart(context.Background()))
	require.NoError(t, lifecycle.hooks[0].OnStop(context.Background()))
	_, err := alerts.handleMatchEnded(
		context.Background(),
		busdota.MatchEndedMessage{ChannelID: uuid.NewString()},
	)

	require.ErrorIs(t, err, context.Canceled)
	require.Empty(t, repository.channel)
	require.Zero(t, store.operationCalls())
	require.Empty(t, publisher.requests)
}
