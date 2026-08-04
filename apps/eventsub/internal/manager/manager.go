package manager

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/google/uuid"
	"github.com/kvizyx/twitchy/eventsub"
	goredislib "github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/apps/eventsub/internal/handler"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	buscore "github.com/twirapp/twir/libs/bus-core"
	cfg "github.com/twirapp/twir/libs/config"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
	twitchconduits "github.com/twirapp/twir/libs/repositories/twitch_conduits"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	twitchlib "github.com/twirapp/twir/libs/twitch"
	"go.uber.org/atomic"
	"gorm.io/gorm"
)

type Manager struct {
	config             cfg.Config
	logger             *slog.Logger
	gorm               *gorm.DB
	twirBus            *buscore.Bus
	channelsRepo       channelsrepo.Repository
	channelService     *channelservice.ChannelService
	conduitsRepository twitchconduits.Repository
	redSync            *redsync.Redsync
	eventsub           eventsub.EventSub
	handler            *handler.Handler

	httpClient *http.Client
	apiBaseUrl string
	wsOpts     []eventsub.WebsocketOption

	wsCurrentSessionId *string
	currentConduit     *conduitsResponseConduit

	workerCancel context.CancelFunc
	workerDone   chan struct{}
}

type Opts struct {
	Lc *lifecycle.Lifecycle

	Config             cfg.Config
	Logger             *slog.Logger
	Gorm               *gorm.DB
	TwirBus            *buscore.Bus
	ChannelsRepo       channelsrepo.Repository
	ChannelService     *channelservice.ChannelService
	ConduitsRepository twitchconduits.Repository
	Redis              *goredislib.Client
	Handler            *handler.Handler
}

func NewManager(opts Opts) (*Manager, error) {
	var httpClient *http.Client
	var apiBaseUrl string
	var wsOpts []eventsub.WebsocketOption

	if opts.Config.TwitchMockEnabled {
		httpClient = &http.Client{
			Transport: twitchlib.NewMockRoundTripper(http.DefaultTransport, opts.Config),
		}
		apiBaseUrl = strings.TrimSuffix(opts.Config.TwitchMockApiUrl, "/helix")
		wsOpts = append(wsOpts, eventsub.WebsocketWithServerURL(opts.Config.TwitchMockWsUrl))
	} else {
		httpClient = http.DefaultClient
		apiBaseUrl = "https://api.twitch.tv"
	}

	workerCtx, workerCancel := context.WithCancel(context.Background())

	manager := &Manager{
		config:             opts.Config,
		logger:             opts.Logger,
		gorm:               opts.Gorm,
		twirBus:            opts.TwirBus,
		channelsRepo:       opts.ChannelsRepo,
		channelService:     opts.ChannelService,
		conduitsRepository: opts.ConduitsRepository,
		redSync:            redsync.New(goredis.NewPool(opts.Redis)),
		eventsub:           eventsub.New(),
		handler:            opts.Handler,
		httpClient:         httpClient,
		apiBaseUrl:         apiBaseUrl,
		wsOpts:             wsOpts,
		wsCurrentSessionId: nil,
		currentConduit:     nil,
		workerCancel:       workerCancel,
		workerDone:         make(chan struct{}),
	}

	opts.Lc.Append(
		lifecycle.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					defer close(manager.workerDone)
					manager.runStartup(workerCtx)
				}()

				return nil
			},
			OnStop: func(ctx context.Context) error {
				workerCancel()

				select {
				case <-manager.workerDone:
					return nil
				case <-ctx.Done():
					return ctx.Err()
				}
			},
		},
	)

	return manager, nil
}

func (c *Manager) SubscribeToNeededEvents(
	ctx context.Context,
	topics []model.EventsubTopic,
	broadcasterId,
	botId string,
) error {
	if c.currentConduit == nil {
		return errors.New("current conduit is not set")
	}

	var wg sync.WaitGroup
	newSubsCount := atomic.NewInt64(0)

	for _, topic := range topics {
		wg.Add(1)

		topic := topic
		go func() {
			defer wg.Done()

			err := c.SubscribeWithLimits(
				ctx,
				eventsub.EventType(topic.Topic),
				eventsub.ConduitTransport{
					Method:    "conduit",
					ConduitId: c.currentConduit.Id,
				},
				topic.Version,
				broadcasterId,
				botId,
			)

			if err != nil {
				c.logger.Error(
					"failed to subscribe to event",
					logger.Error(err),
					slog.Any("topic", topic.Topic),
					slog.String("version", topic.Version),
				)
			}

			newSubsCount.Inc()
		}()
	}

	wg.Wait()

	if newSubsCount.Load() > 0 {
		c.logger.Info(
			"New subscriptions created for channel",
			slog.String("channel_id", broadcasterId),
			slog.String("bot_id", botId),
			slog.Int64("count", newSubsCount.Load()),
		)
	}

	return nil
}

func (c *Manager) SubscribeToEvent(
	ctx context.Context,
	topic,
	version,
	channelId string,
) error {
	broadcasterID, botID, err := c.resolveTwitchSubscriptionIdentities(ctx, channelId)
	if err != nil {
		return err
	}

	err = c.SubscribeWithLimits(
		ctx,
		eventsub.EventType(topic),
		eventsub.ConduitTransport{
			Method:    "conduit",
			ConduitId: c.currentConduit.Id,
		},
		version,
		broadcasterID,
		botID,
	)

	return err
}

func (c *Manager) resolveTwitchSubscriptionIdentities(ctx context.Context, channelID string) (string, string, error) {
	channelUUID, err := uuid.Parse(channelID)
	if err != nil {
		return "", "", err
	}

	channel, err := c.channelService.GetChannelByID(ctx, channelUUID)
	if err != nil {
		return "", "", err
	}

	binding, ok := channel.Binding(platformentity.PlatformTwitch)
	if !ok || binding.PlatformChannelID == "" {
		return "", "", errors.New("channel has no twitch platform id")
	}

	botConfig, err := binding.ParseTwitchBotConfig()
	if err != nil {
		return "", "", fmt.Errorf("parse Twitch bot config: %w", err)
	}

	return binding.PlatformChannelID, botConfig.BotID, nil
}
