package handler

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/redis/go-redis/v9"
	duplicate_tracker "github.com/satont/twir/apps/eventsub/internal/duplicate-tracker"
	"github.com/satont/twir/apps/eventsub/internal/manager"
	"github.com/satont/twir/apps/eventsub/internal/tunnel"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	bus_core "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/grpc/websockets"
	channelmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	"github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
	channelsinfohistory "github.com/twirapp/twir/libs/repositories/channels_info_history"
	scheduledvipsrepository "github.com/twirapp/twir/libs/repositories/scheduled_vips"
	"github.com/twirapp/twir/libs/repositories/streams"
	eventsub_framework "github.com/twirapp/twitch-eventsub-framework"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Handler struct {
	logger logger.Logger

	eventsGrpc              events.EventsClient
	parserGrpc              parser.ParserClient
	websocketsGrpc          websockets.WebsocketClient
	tokensGrpc              tokens.TokensClient
	tracer                  trace.Tracer
	manager                 *manager.Manager
	scheduledVipsRepo       scheduledvipsrepository.Repository
	channelsCache           *generic_cacher.GenericCacher[channelmodel.Channel]
	channelsInfoHistoryRepo channelsinfohistory.Repository
	streamsrepository       streams.Repository

	gorm        *gorm.DB
	redisClient *redis.Client

	bus         *bus_core.Bus
	prefixCache *generic_cacher.GenericCacher[model.ChannelsCommandsPrefix]
	config      cfg.Config
}

type Opts struct {
	fx.In
	Lc fx.Lifecycle

	Logger logger.Logger

	EventsGrpc              events.EventsClient
	ParserGrpc              parser.ParserClient
	WebsocketsGrpc          websockets.WebsocketClient
	TokensGrpc              tokens.TokensClient
	ScheduledVipsRepo       scheduledvipsrepository.Repository
	ChannelsRepo            *generic_cacher.GenericCacher[channelmodel.Channel]
	ChannelsInfoHistoryRepo channelsinfohistory.Repository
	StreamsRepository       streams.Repository

	Tracer  trace.Tracer
	Tunn    *tunnel.AppTunnel
	Manager *manager.Manager
	Gorm    *gorm.DB
	Redis   *redis.Client

	Bus         *bus_core.Bus
	PrefixCache *generic_cacher.GenericCacher[model.ChannelsCommandsPrefix]

	Config cfg.Config
}

func New(opts Opts) *Handler {
	handler := eventsub_framework.NewSubHandler(true, []byte(opts.Config.TwitchClientSecret))
	handler.IDTracker = duplicate_tracker.New(duplicate_tracker.Opts{Redis: opts.Redis})

	myHandler := &Handler{
		manager:                 opts.Manager,
		logger:                  opts.Logger,
		config:                  opts.Config,
		gorm:                    opts.Gorm,
		redisClient:             opts.Redis,
		eventsGrpc:              opts.EventsGrpc,
		parserGrpc:              opts.ParserGrpc,
		websocketsGrpc:          opts.WebsocketsGrpc,
		tokensGrpc:              opts.TokensGrpc,
		tracer:                  opts.Tracer,
		bus:                     opts.Bus,
		prefixCache:             opts.PrefixCache,
		scheduledVipsRepo:       opts.ScheduledVipsRepo,
		channelsCache:           opts.ChannelsRepo,
		channelsInfoHistoryRepo: opts.ChannelsInfoHistoryRepo,
		streamsrepository:       opts.StreamsRepository,
	}

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

	httpHandler := otelhttp.NewHandler(handler, "")

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

				opts.Logger.Info("Handler started")

				return nil
			},
		},
	)

	return myHandler
}
