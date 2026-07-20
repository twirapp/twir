package chatalerts

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/kv"
	kvoptions "github.com/twirapp/kv/options"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/bots"
	busdota "github.com/twirapp/twir/libs/bus-core/dota"
	"github.com/twirapp/twir/libs/logger"
	dotarepository "github.com/twirapp/twir/libs/repositories/dota"
	dotamodel "github.com/twirapp/twir/libs/repositories/dota/model"
	"go.uber.org/fx"
)

const (
	dotaSubscriptionGroup = "dota"
	cooldownKeyPrefix     = "cache:twir:dota:chat-alert"
)

type eventKind string

const (
	eventMatchStarted eventKind = "match_started"
	eventMatchEnded   eventKind = "match_ended"
	eventRoshanKilled eventKind = "roshan_killed"
	eventAegisPickup  eventKind = "aegis_pickup"
)

type settingsRepository interface {
	GetByChannelID(context.Context, uuid.UUID) (dotamodel.ChannelDotaSettings, error)
}

type messagePublisher interface {
	Publish(context.Context, bots.SendMessageRequest) error
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

	Lifecycle fx.Lifecycle
	Bus       *buscore.Bus
	KV        kv.KV
	Logger    *slog.Logger

	SettingsRepository dotarepository.Repository
}

type ChatAlerts struct {
	repository    settingsRepository
	cache         kv.KV
	messages      messagePublisher
	logger        *slog.Logger
	subscriptions []subscription
}

func New(opts Opts) *ChatAlerts {
	alerts := newChatAlerts(
		opts.SettingsRepository,
		opts.KV,
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
	cache kv.KV,
	messages messagePublisher,
	logger *slog.Logger,
	subscriptions []subscription,
	lifecycle fx.Lifecycle,
) *ChatAlerts {
	alerts := &ChatAlerts{
		repository:    repository,
		cache:         cache,
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

func (c *ChatAlerts) Stop(_ context.Context) error {
	for _, subscription := range c.subscriptions {
		subscription.Unsubscribe()
	}

	return nil
}

func (c *ChatAlerts) handleMatchStarted(
	ctx context.Context,
	message busdota.MatchStartedMessage,
) (struct{}, error) {
	return struct{}{}, c.handle(
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
	return struct{}{}, c.handle(
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
	return struct{}{}, c.handle(
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
	return struct{}{}, c.handle(
		ctx,
		chatEvent{
			kind:      eventAegisPickup,
			channelID: message.ChannelID,
			player:    message.PlayerName,
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

	cooldownKey := fmt.Sprintf("%s:%s:%s", cooldownKeyPrefix, channelID, event.kind)
	if err := c.cache.Get(ctx, cooldownKey).Err(); err == nil {
		return nil
	} else if !errors.Is(err, kv.ErrKeyNil) {
		c.logger.ErrorContext(ctx, "dota chat alert: check cooldown", logger.Error(err))
		return fmt.Errorf("get dota chat alert cooldown: %w", err)
	}

	if event.kind != eventMatchEnded {
		event.mmr = settings.Mmr
		event.wins = settings.SessionWins
		event.losses = settings.SessionLosses
	}

	if err := c.messages.Publish(
		ctx,
		bots.SendMessageRequest{
			InternalChannelID: &channelID,
			Platform:          "twitch",
			Message:           renderTemplate(eventSettings.Template, event),
			SkipRateLimits:    true,
		},
	); err != nil {
		return fmt.Errorf("publish dota chat alert: %w", err)
	}

	if eventSettings.Cooldown > 0 {
		if err := c.cache.Set(
			ctx,
			cooldownKey,
			"1",
			kvoptions.WithExpire(time.Duration(eventSettings.Cooldown)*time.Second),
		); err != nil {
			return fmt.Errorf("set dota chat alert cooldown: %w", err)
		}
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
