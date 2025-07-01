package handler

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	duplicate_tracker "github.com/satont/twir/apps/eventsub/internal/duplicate-tracker"
	"github.com/satont/twir/apps/eventsub/internal/manager"
	"github.com/satont/twir/apps/eventsub/internal/tunnel"
	cfg "github.com/satont/twir/libs/config"
	deprecatedmodel "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	batchprocessor "github.com/twirapp/batch-processor"
	bus_core "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/grpc/websockets"
	channelmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	channelscommandsprefixmodel "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	channelsinfohistory "github.com/twirapp/twir/libs/repositories/channels_info_history"
	channelsredemptionshistory "github.com/twirapp/twir/libs/repositories/channels_redemptions_history"
	scheduledvipsrepository "github.com/twirapp/twir/libs/repositories/scheduled_vips"
	"github.com/twirapp/twir/libs/repositories/streams"
	eventsub_framework "github.com/twirapp/twitch-eventsub-framework"
	eventsub_bindings "github.com/twirapp/twitch-eventsub-framework/esb"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"gorm.io/gorm"

	alertmodel "github.com/twirapp/twir/libs/repositories/alerts/model"
)

type Handler struct {
	logger logger.Logger

	websocketsGrpc               websockets.WebsocketClient
	tracer                       trace.Tracer
	manager                      *manager.Manager
	scheduledVipsRepo            scheduledvipsrepository.Repository
	channelsCache                *generic_cacher.GenericCacher[channelmodel.Channel]
	channelsInfoHistoryRepo      channelsinfohistory.Repository
	streamsrepository            streams.Repository
	redemptionsHistoryRepository channelsredemptionshistory.Repository
	eventsListRepository         channelseventslist.Repository

	gorm        *gorm.DB
	redisClient *redis.Client

	twirBus                             *bus_core.Bus
	prefixCache                         *generic_cacher.GenericCacher[channelscommandsprefixmodel.ChannelsCommandsPrefix]
	alertsCache                         *generic_cacher.GenericCacher[[]alertmodel.Alert]
	commandsCache                       *generic_cacher.GenericCacher[[]deprecatedmodel.ChannelsCommands]
	channelSongRequestsSettingsCache    *generic_cacher.GenericCacher[deprecatedmodel.ChannelSongRequestsSettings]
	channelsIntegrationsSettingsSeventv *generic_cacher.GenericCacher[deprecatedmodel.ChannelsIntegrationsSettingsSeventv]
	config                              cfg.Config

	redemptionsBatcher *batchprocessor.BatchProcessor[eventsub_bindings.EventChannelPointsRewardRedemptionAdd]
}

type Opts struct {
	fx.In
	Lc fx.Lifecycle

	Logger logger.Logger

	WebsocketsGrpc                      websockets.WebsocketClient
	ScheduledVipsRepo                   scheduledvipsrepository.Repository
	ChannelsRepo                        *generic_cacher.GenericCacher[channelmodel.Channel]
	ChannelsInfoHistoryRepo             channelsinfohistory.Repository
	StreamsRepository                   streams.Repository
	RedemptionsHistoryRepository        channelsredemptionshistory.Repository
	EventsListRepository                channelseventslist.Repository
	CommandsCache                       *generic_cacher.GenericCacher[[]deprecatedmodel.ChannelsCommands]
	ChannelSongRequestsSettingsCache    *generic_cacher.GenericCacher[deprecatedmodel.ChannelSongRequestsSettings]
	ChannelsIntegrationsSettingsSeventv *generic_cacher.GenericCacher[deprecatedmodel.ChannelsIntegrationsSettingsSeventv]

	Tracer  trace.Tracer
	Tunn    *tunnel.AppTunnel
	Manager *manager.Manager
	Gorm    *gorm.DB
	Redis   *redis.Client

	Bus                *bus_core.Bus
	PrefixCache        *generic_cacher.GenericCacher[channelscommandsprefixmodel.ChannelsCommandsPrefix]
	ChannelAlertsCache *generic_cacher.GenericCacher[[]alertmodel.Alert]

	Config cfg.Config
}

func New(opts Opts) *Handler {
	var doSignatureVerification = true
	if opts.Config.EventSubDisableSignatureVerification {
		doSignatureVerification = false
	}

	handler := eventsub_framework.NewSubHandler(
		doSignatureVerification,
		[]byte(opts.Config.TwitchClientSecret),
	)
	handler.IDTracker = duplicate_tracker.New(duplicate_tracker.Opts{Redis: opts.Redis})

	myHandler := &Handler{
		manager:                             opts.Manager,
		logger:                              opts.Logger,
		config:                              opts.Config,
		gorm:                                opts.Gorm,
		redisClient:                         opts.Redis,
		websocketsGrpc:                      opts.WebsocketsGrpc,
		tracer:                              opts.Tracer,
		twirBus:                             opts.Bus,
		prefixCache:                         opts.PrefixCache,
		scheduledVipsRepo:                   opts.ScheduledVipsRepo,
		channelsCache:                       opts.ChannelsRepo,
		channelsInfoHistoryRepo:             opts.ChannelsInfoHistoryRepo,
		streamsrepository:                   opts.StreamsRepository,
		redemptionsHistoryRepository:        opts.RedemptionsHistoryRepository,
		eventsListRepository:                opts.EventsListRepository,
		alertsCache:                         opts.ChannelAlertsCache,
		commandsCache:                       opts.CommandsCache,
		channelSongRequestsSettingsCache:    opts.ChannelSongRequestsSettingsCache,
		channelsIntegrationsSettingsSeventv: opts.ChannelsIntegrationsSettingsSeventv,
	}

	batcherCtx, batcherStop := context.WithCancel(context.Background())

	myHandler.redemptionsBatcher = batchprocessor.NewBatchProcessor[eventsub_bindings.EventChannelPointsRewardRedemptionAdd](
		batchprocessor.BatchProcessorOpts[eventsub_bindings.EventChannelPointsRewardRedemptionAdd]{
			Interval:  500 * time.Millisecond,
			BatchSize: 100,
			Callback:  myHandler.handleChannelPointsRewardRedemptionAddBatched,
		},
	)

	handler.OnNotification = myHandler.onNotification
	handler.HandleUserAuthorizationRevoke = myHandler.handleUserAuthorizationRevoke
	handler.OnRevocate = myHandler.handleSubRevocate

	handler.HandleChannelUpdate = myHandler.handleChannelUpdate
	handler.HandleStreamOnline = myHandler.handleStreamOnline
	handler.HandleStreamOffline = myHandler.handleStreamOffline
	handler.HandleUserUpdate = myHandler.handleUserUpdate
	handler.HandleChannelFollow = myHandler.handleChannelFollow
	handler.HandleChannelModeratorAdd = myHandler.handleChannelModeratorAdd
	handler.HandleChannelModeratorRemove = myHandler.handleChannelModeratorRemove
	handler.HandleChannelPointsRewardRedemptionAdd = myHandler.handleChannelPointsRewardRedemptionAdd
	handler.HandleChannelPointsRewardRedemptionUpdate = myHandler.handleChannelPointsRewardRedemptionUpdate
	handler.HandleChannelPollBegin = myHandler.handleChannelPollBegin
	handler.HandleChannelPollProgress = myHandler.handleChannelPollProgress
	handler.HandleChannelPollEnd = myHandler.handleChannelPollEnd
	handler.HandleChannelPredictionBegin = myHandler.handleChannelPredictionBegin
	handler.HandleChannelPredictionProgress = myHandler.handleChannelPredictionProgress
	handler.HandleChannelPredictionLock = myHandler.handleChannelPredictionLock
	handler.HandleChannelPredictionEnd = myHandler.handleChannelPredictionEnd
	handler.HandleChannelBan = myHandler.handleBan
	handler.HandleChannelSubscribe = myHandler.handleChannelSubscribe
	handler.HandleChannelSubscriptionMessage = myHandler.handleChannelSubscriptionMessage
	handler.HandleChannelRaid = myHandler.handleChannelRaid
	handler.HandleChannelChatClear = myHandler.handleChannelChatClear
	handler.HandleChannelChatNotification = myHandler.handleChannelChatNotification
	handler.HandleChannelChatMessage = myHandler.handleChannelChatMessage
	handler.HandleChannelUnbanRequestCreate = myHandler.handleChannelUnbanRequestCreate
	handler.HandleChannelUnbanRequestResolve = myHandler.handleChannelUnbanRequestResolve
	handler.HandleChannelChatMessageDelete = myHandler.handleChannelChatMessageDelete
	handler.HandleChannelPointsRewardAdd = myHandler.handleChannelPointsRewardAdd
	handler.HandleChannelPointsRewardUpdate = myHandler.handleChannelPointsRewardUpdate
	handler.HandleChannelPointsRewardRemove = myHandler.handleChannelPointsRewardRemove
	handler.HandleChannelVipAdd = myHandler.handleChannelVipAdd
	handler.HandleChannelVipRemove = myHandler.handleChannelVipRemove

	mux := http.NewServeMux()
	// middleware
	mux.Handle(
		"POST /", http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				span := trace.SpanFromContext(r.Context())
				defer span.End()

				headers := r.Header

				span.SetAttributes(
					attribute.String("twitcheventsub.retry", headers.Get("Twitch-Eventsub-Message-Retry")),
					attribute.String(
						"twitcheventsub.message_type",
						headers.Get("Twitch-Eventsub-Message-Type"),
					),
					attribute.String(
						"twitcheventsub.subscription_type",
						headers.Get("Twitch-Eventsub-Subscription-Type"),
					),
					attribute.String(
						"twitcheventsub.subscription_version",
						headers.Get("Twitch-Eventsub-Subscription-Version"),
					),
				)

				if r.Method != "POST" {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
					return
				}

				r = r.WithContext(context.WithoutCancel(r.Context()))

				handler.ServeHTTP(w, r)
			},
		),
	)

	httpHandler := otelhttp.NewHandler(mux, "eventsub-server")

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					if err := http.Serve(opts.Tunn, httpHandler); err != nil && !errors.Is(
						err,
						net.ErrClosed,
					) {
						panic(err)
					}
				}()

				go func() {
					myHandler.redemptionsBatcher.Start(batcherCtx)
				}()

				opts.Logger.Info("Handler started")

				return nil
			},
			OnStop: func(ctx context.Context) error {
				if err := myHandler.redemptionsBatcher.Shutdown(ctx); err != nil {
					return err
				}
				batcherStop()
				return nil
			},
		},
	)

	return myHandler
}
