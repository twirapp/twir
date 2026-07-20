package chatalerts

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/bots"
	busdota "github.com/twirapp/twir/libs/bus-core/dota"
	dotarepository "github.com/twirapp/twir/libs/repositories/dota"
	dotamodel "github.com/twirapp/twir/libs/repositories/dota/model"
	"go.uber.org/fx"
)

const (
	dotaSubscriptionGroup  = "dota"
	cooldownKeyPrefix      = "cache:twir:dota:chat-alert"
	maxCooldownSeconds     = math.MaxInt64 / int64(time.Second)
	compareAndDeleteScript = `if redis.call("GET", KEYS[1]) == ARGV[1] then return redis.call("DEL", KEYS[1]) end return 0`
)

type eventKind string

const (
	eventMatchStarted eventKind = "match_started"
	eventMatchEnded   eventKind = "match_ended"
	eventRoshanKilled eventKind = "roshan_killed"
	eventAegisPickup  eventKind = "aegis_pickup"
)

type terminalHandlerError struct {
	kind  eventKind
	cause error
}

func (e terminalHandlerError) Error() string {
	return fmt.Sprintf("dota chat alert handler failed: %s", e.kind)
}

func (e terminalHandlerError) Unwrap() error {
	return e.cause
}

type settingsRepository interface {
	GetByChannelID(context.Context, uuid.UUID) (dotamodel.ChannelDotaSettings, error)
}

type messagePublisher interface {
	Publish(context.Context, bots.SendMessageRequest) error
}

type CooldownStore interface {
	Reserve(ctx context.Context, key string, token string, ttl time.Duration) (bool, error)
	Release(ctx context.Context, key string, token string) error
}

type redisCooldownClient interface {
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd
}

type RedisCooldownStore struct {
	client redisCooldownClient
}

func NewRedisCooldownStore(client *redis.Client) *RedisCooldownStore {
	return &RedisCooldownStore{client: client}
}

func (s *RedisCooldownStore) Reserve(
	ctx context.Context,
	key string,
	token string,
	ttl time.Duration,
) (bool, error) {
	return s.client.SetNX(ctx, key, token, ttl).Result()
}

func (s *RedisCooldownStore) Release(ctx context.Context, key string, token string) error {
	return s.client.Eval(
		ctx,
		compareAndDeleteScript,
		[]string{key},
		token,
	).Err()
}

type subscription interface {
	Subscribe(group string) error
	Unsubscribe()
}

type subscriptionFunc struct {
	subscribe   func(group string) error
	unsubscribe func()
}

func (s subscriptionFunc) Subscribe(group string) error {
	return s.subscribe(group)
}

func (s subscriptionFunc) Unsubscribe() {
	s.unsubscribe()
}

type busMessagePublisher struct {
	bus *buscore.Bus
}

func (p busMessagePublisher) Publish(ctx context.Context, request bots.SendMessageRequest) error {
	return p.bus.Bots.SendMessage.Publish(ctx, request)
}

type Opts struct {
	fx.In

	Lifecycle     fx.Lifecycle
	Bus           *buscore.Bus
	CooldownStore CooldownStore
	Logger        *slog.Logger

	SettingsRepository dotarepository.Repository
}

type ChatAlerts struct {
	repository    settingsRepository
	cooldownStore CooldownStore
	messages      messagePublisher
	logger        *slog.Logger
	subscriptions []subscription

	handlersMu sync.Mutex
	handlers   sync.WaitGroup
	stopping   bool
}

func New(opts Opts) *ChatAlerts {
	alerts := newChatAlerts(
		opts.SettingsRepository,
		opts.CooldownStore,
		busMessagePublisher{bus: opts.Bus},
		opts.Logger,
		nil,
		opts.Lifecycle,
	)
	alerts.subscriptions = newBusSubscriptions(opts.Bus, alerts)

	return alerts
}

func newChatAlerts(
	repository settingsRepository,
	cooldownStore CooldownStore,
	messages messagePublisher,
	logger *slog.Logger,
	subscriptions []subscription,
	lifecycle fx.Lifecycle,
) *ChatAlerts {
	alerts := &ChatAlerts{
		repository:    repository,
		cooldownStore: cooldownStore,
		messages:      messages,
		logger:        logger,
		subscriptions: subscriptions,
	}

	if lifecycle != nil {
		lifecycle.Append(
			fx.Hook{
				OnStart: alerts.Start,
				OnStop:  alerts.Stop,
			},
		)
	}

	return alerts
}

func newBusSubscriptions(bus *buscore.Bus, alerts *ChatAlerts) []subscription {
	return []subscription{
		subscriptionFunc{
			subscribe: func(group string) error {
				return bus.Dota.MatchStarted.SubscribeGroup(group, alerts.handleMatchStarted)
			},
			unsubscribe: bus.Dota.MatchStarted.Unsubscribe,
		},
		subscriptionFunc{
			subscribe: func(group string) error {
				return bus.Dota.MatchEnded.SubscribeGroup(group, alerts.handleMatchEnded)
			},
			unsubscribe: bus.Dota.MatchEnded.Unsubscribe,
		},
		subscriptionFunc{
			subscribe: func(group string) error {
				return bus.Dota.RoshanKilled.SubscribeGroup(group, alerts.handleRoshanKilled)
			},
			unsubscribe: bus.Dota.RoshanKilled.Unsubscribe,
		},
		subscriptionFunc{
			subscribe: func(group string) error {
				return bus.Dota.AegisPickup.SubscribeGroup(group, alerts.handleAegisPickup)
			},
			unsubscribe: bus.Dota.AegisPickup.Unsubscribe,
		},
	}
}

func (c *ChatAlerts) Start(_ context.Context) error {
	started := make([]subscription, 0, len(c.subscriptions))
	for _, subscription := range c.subscriptions {
		if err := subscription.Subscribe(dotaSubscriptionGroup); err != nil {
			for i := len(started) - 1; i >= 0; i-- {
				started[i].Unsubscribe()
			}
			return fmt.Errorf("subscribe to dota chat alerts: %w", err)
		}
		started = append(started, subscription)
	}

	return nil
}

func (c *ChatAlerts) Stop(ctx context.Context) error {
	c.handlersMu.Lock()
	c.stopping = true
	for _, subscription := range c.subscriptions {
		subscription.Unsubscribe()
	}
	c.handlersMu.Unlock()

	handlersDone := make(chan struct{})
	go func() {
		c.handlers.Wait()
		close(handlersDone)
	}()

	select {
	case <-handlersDone:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}

}

func (c *ChatAlerts) handleMatchStarted(
	ctx context.Context,
	message busdota.MatchStartedMessage,
) (struct{}, error) {
	return struct{}{}, c.handleTracked(
		ctx,
		chatEvent{
			kind:      eventMatchStarted,
			channelID: message.ChannelID,
			hero:      message.HeroName,
		},
	)
}

func (c *ChatAlerts) handleMatchEnded(
	ctx context.Context,
	message busdota.MatchEndedMessage,
) (struct{}, error) {
	return struct{}{}, c.handleTracked(
		ctx,
		chatEvent{
			kind:      eventMatchEnded,
			channelID: message.ChannelID,
			hero:      message.HeroName,
			mmr:       message.Mmr,
			wins:      message.SessionWins,
			losses:    message.SessionLosses,
		},
	)
}

func (c *ChatAlerts) handleRoshanKilled(
	ctx context.Context,
	message busdota.RoshanKilledMessage,
) (struct{}, error) {
	return struct{}{}, c.handleTracked(
		ctx,
		chatEvent{
			kind:      eventRoshanKilled,
			channelID: message.ChannelID,
			team:      message.Team,
			gameTime:  formatGameTime(message.GameTime),
		},
	)
}

func (c *ChatAlerts) handleAegisPickup(
	ctx context.Context,
	message busdota.AegisPickupMessage,
) (struct{}, error) {
	player := message.PlayerName
	if player == "" && message.PlayerID != nil {
		// GSI has a slot ID only; roster lookup is deferred.
		player = fmt.Sprintf("player #%d", *message.PlayerID)
	}

	return struct{}{}, c.handleTracked(
		ctx,
		chatEvent{
			kind:      eventAegisPickup,
			channelID: message.ChannelID,
			player:    player,
			gameTime:  formatGameTime(message.GameTime),
		},
	)
}

type chatEvent struct {
	kind      eventKind
	channelID string
	hero      string
	team      string
	player    string
	gameTime  string
	mmr       int
	wins      int
	losses    int
}

func (c *ChatAlerts) handleTracked(ctx context.Context, event chatEvent) error {
	c.handlersMu.Lock()
	if c.stopping {
		c.handlersMu.Unlock()
		return context.Canceled
	}
	c.handlers.Add(1)
	c.handlersMu.Unlock()
	defer c.handlers.Done()

	err := c.handle(ctx, event)
	if err != nil {
		c.logger.ErrorContext(
			ctx,
			"dota chat alert handler failed",
			slog.String("event_kind", string(event.kind)),
		)
		return terminalHandlerError{kind: event.kind, cause: err}
	}
	return nil
}

func (c *ChatAlerts) handle(ctx context.Context, event chatEvent) error {
	channelID, err := uuid.Parse(event.channelID)
	if err != nil {
		return fmt.Errorf("parse channel ID: %w", err)
	}

	settings, err := c.repository.GetByChannelID(ctx, channelID)
	if errors.Is(err, dotarepository.ErrNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("get dota settings: %w", err)
	}
	if !settings.Enabled {
		return nil
	}

	eventSettings := selectEventSettings(settings.ChatEvents, event.kind)
	if !eventSettings.Enabled || strings.TrimSpace(eventSettings.Template) == "" {
		return nil
	}

	if event.kind != eventMatchEnded {
		event.mmr = settings.Mmr
		event.wins = settings.SessionWins
		event.losses = settings.SessionLosses
	}

	message := strings.TrimSpace(renderTemplate(eventSettings.Template, event))
	if message == "" {
		return nil
	}

	request := bots.SendMessageRequest{
		InternalChannelID: &channelID,
		Platform:          "twitch",
		Message:           message,
		SkipRateLimits:    true,
	}

	if eventSettings.Cooldown <= 0 {
		if err := c.messages.Publish(ctx, request); err != nil {
			return fmt.Errorf("publish dota chat alert: %w", err)
		}

		return nil
	}
	if int64(eventSettings.Cooldown) > maxCooldownSeconds {
		return fmt.Errorf("dota chat alert cooldown exceeds maximum duration")
	}

	cooldownKey := fmt.Sprintf("%s:%s:%s", cooldownKeyPrefix, channelID, event.kind)
	cooldownToken := uuid.NewString()
	reserved, err := c.cooldownStore.Reserve(
		ctx,
		cooldownKey,
		cooldownToken,
		time.Duration(eventSettings.Cooldown)*time.Second,
	)
	if err != nil {
		return fmt.Errorf("reserve dota chat alert cooldown: %w", err)
	}
	if !reserved {
		return nil
	}

	if err := c.messages.Publish(ctx, request); err != nil {
		cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
		defer cancel()
		if releaseErr := c.cooldownStore.Release(cleanupCtx, cooldownKey, cooldownToken); releaseErr != nil {
			c.logger.ErrorContext(cleanupCtx, "dota chat alert: cooldown cleanup failed")
		}
		return fmt.Errorf("publish dota chat alert: %w", err)
	}

	return nil
}

func selectEventSettings(events dotamodel.ChatEvents, kind eventKind) dotamodel.ChatEventSettings {
	switch kind {
	case eventMatchStarted:
		return events.MatchStarted
	case eventMatchEnded:
		return events.MatchEnded
	case eventRoshanKilled:
		return events.RoshanKilled
	case eventAegisPickup:
		return events.AegisPickup
	default:
		return dotamodel.ChatEventSettings{}
	}
}

func renderTemplate(template string, event chatEvent) string {
	template = strings.ReplaceAll(template, "{hero}", event.hero)
	template = strings.ReplaceAll(template, "{mmr}", strconv.Itoa(event.mmr))
	template = strings.ReplaceAll(template, "{wins}", strconv.Itoa(event.wins))
	template = strings.ReplaceAll(template, "{losses}", strconv.Itoa(event.losses))
	template = strings.ReplaceAll(template, "{team}", event.team)
	template = strings.ReplaceAll(template, "{player}", event.player)
	return strings.ReplaceAll(template, "{time}", event.gameTime)
}

func formatGameTime(seconds int) string {
	return fmt.Sprintf("%d:%02d", seconds/60, seconds%60)
}
