package handler

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/eventsub/internal/manager"
	"github.com/satont/twir/apps/eventsub/internal/tunnel"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	bus_core "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/cache/7tv"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/grpc/websockets"
	seventvintegration "github.com/twirapp/twir/libs/integrations/seventv"
	eventsub_framework "github.com/twirapp/twitch-eventsub-framework"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Handler struct {
	manager *manager.Manager

	logger      logger.Logger
	config      cfg.Config
	gorm        *gorm.DB
	redisClient *redis.Client

	eventsGrpc     events.EventsClient
	parserGrpc     parser.ParserClient
	websocketsGrpc websockets.WebsocketClient
	tokensGrpc     tokens.TokensClient
	tracer         trace.Tracer
	bus            *bus_core.Bus
	seventvCache   *generic_cacher.GenericCacher[*seventvintegration.ProfileResponse]
}

type Opts struct {
	fx.In
	Lc fx.Lifecycle

	Config  cfg.Config
	Tunn    *tunnel.AppTunnel
	Manager *manager.Manager
	Logger  logger.Logger
	Gorm    *gorm.DB
	Redis   *redis.Client

	EventsGrpc     events.EventsClient
	ParserGrpc     parser.ParserClient
	WebsocketsGrpc websockets.WebsocketClient
	TokensGrpc     tokens.TokensClient
	Bus            *bus_core.Bus

	Tracer trace.Tracer
}

func New(opts Opts) *Handler {
	handler := eventsub_framework.NewSubHandler(true, []byte(opts.Config.TwitchClientSecret))

	myHandler := &Handler{
		manager:        opts.Manager,
		logger:         opts.Logger,
		config:         opts.Config,
		gorm:           opts.Gorm,
		redisClient:    opts.Redis,
		eventsGrpc:     opts.EventsGrpc,
		parserGrpc:     opts.ParserGrpc,
		websocketsGrpc: opts.WebsocketsGrpc,
		tokensGrpc:     opts.TokensGrpc,
		tracer:         opts.Tracer,
		bus:            opts.Bus,
		seventvCache:   seventv.New(opts.Redis),
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
