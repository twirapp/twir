package handler

import (
	"context"
	"log/slog"
	"time"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/redis/go-redis/v9"
	batchprocessor "github.com/twirapp/batch-processor"
	user_creator "github.com/twirapp/twir/apps/eventsub/internal/services/user-creator"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	bus_core "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	cfg "github.com/twirapp/twir/libs/config"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	deprecatedmodel "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/websockets"
	channelplatforms "github.com/twirapp/twir/libs/repositories/channel_platforms"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelscommandsprefixmodel "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	channelsinfohistory "github.com/twirapp/twir/libs/repositories/channels_info_history"
	channelsredemptionshistory "github.com/twirapp/twir/libs/repositories/channels_redemptions_history"
	commandswithgroupsandresponsesmodel "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"
	scheduledvipsrepository "github.com/twirapp/twir/libs/repositories/scheduled_vips"
	"github.com/twirapp/twir/libs/repositories/streams"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"

	alertmodel "github.com/twirapp/twir/libs/repositories/alerts/model"
)

type Handler struct {
	logger *slog.Logger

	websocketsGrpc               websockets.WebsocketClient
	tracer                       trace.Tracer
	scheduledVipsRepo            scheduledvipsrepository.Repository
	channelPlatformsRepo         channelplatforms.Repository
	channelsRepo                 channelsrepository.Repository
	channelsCache                *generic_cacher.GenericCacher[channelentity.Channel]
	channelsInfoHistoryRepo      channelsinfohistory.Repository
	streamsrepository            streams.Repository
	redemptionsHistoryRepository channelsredemptionshistory.Repository
	eventsListRepository         channelseventslist.Repository

	userCreatorService *user_creator.UserCreatorService
	channelService     *channelservice.ChannelService

	gorm        *gorm.DB
	redisClient *redis.Client
	usersRepo   usersrepository.Repository

	twirBus                             *bus_core.Bus
	prefixCache                         *generic_cacher.GenericCacher[channelscommandsprefixmodel.ChannelsCommandsPrefix]
	alertsCache                         *generic_cacher.GenericCacher[[]alertmodel.Alert]
	commandsCache                       *generic_cacher.GenericCacher[[]commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses]
	channelSongRequestsSettingsCache    *generic_cacher.GenericCacher[deprecatedmodel.ChannelSongRequestsSettings]
	channelsIntegrationsSettingsSeventv *generic_cacher.GenericCacher[deprecatedmodel.ChannelsIntegrationsSettingsSeventv]
	config                              cfg.Config

	redemptionsBatcher *batchprocessor.BatchProcessor[eventsub.ChannelPointsCustomRewardRedemptionAddEvent]
}

func New(
	lc *lifecycle.Lifecycle,
	logger *slog.Logger,
	websocketsGrpc websockets.WebsocketClient,
	channelPlatformsRepository channelplatforms.Repository,
	channelsRepository channelsrepository.Repository,
	scheduledVipsRepo scheduledvipsrepository.Repository,
	channelsRepo *generic_cacher.GenericCacher[channelentity.Channel],
	channelsInfoHistoryRepo channelsinfohistory.Repository,
	streamsRepository streams.Repository,
	redemptionsHistoryRepository channelsredemptionshistory.Repository,
	eventsListRepository channelseventslist.Repository,
	commandsCache *generic_cacher.GenericCacher[[]commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses],
	channelSongRequestsSettingsCache *generic_cacher.GenericCacher[deprecatedmodel.ChannelSongRequestsSettings],
	channelsIntegrationsSettingsSeventv *generic_cacher.GenericCacher[deprecatedmodel.ChannelsIntegrationsSettingsSeventv],
	userCreatorService *user_creator.UserCreatorService,
	channelService *channelservice.ChannelService,
	tracer trace.Tracer,
	db *gorm.DB,
	redisClient *redis.Client,
	usersRepo usersrepository.Repository,
	bus *bus_core.Bus,
	prefixCache *generic_cacher.GenericCacher[channelscommandsprefixmodel.ChannelsCommandsPrefix],
	channelAlertsCache *generic_cacher.GenericCacher[[]alertmodel.Alert],
	config cfg.Config,
) *Handler {
	myHandler := &Handler{
		logger:                              logger,
		config:                              config,
		gorm:                                db,
		redisClient:                         redisClient,
		usersRepo:                           usersRepo,
		websocketsGrpc:                      websocketsGrpc,
		tracer:                              tracer,
		twirBus:                             bus,
		prefixCache:                         prefixCache,
		scheduledVipsRepo:                   scheduledVipsRepo,
		channelPlatformsRepo:                channelPlatformsRepository,
		channelsRepo:                        channelsRepository,
		channelsCache:                       channelsRepo,
		channelsInfoHistoryRepo:             channelsInfoHistoryRepo,
		streamsrepository:                   streamsRepository,
		redemptionsHistoryRepository:        redemptionsHistoryRepository,
		eventsListRepository:                eventsListRepository,
		alertsCache:                         channelAlertsCache,
		commandsCache:                       commandsCache,
		channelSongRequestsSettingsCache:    channelSongRequestsSettingsCache,
		channelsIntegrationsSettingsSeventv: channelsIntegrationsSettingsSeventv,
		userCreatorService:                  userCreatorService,
		channelService:                      channelService,
	}

	batcherCtx, batcherStop := context.WithCancel(context.Background())

	myHandler.redemptionsBatcher = batchprocessor.NewBatchProcessor[eventsub.ChannelPointsCustomRewardRedemptionAddEvent](
		batchprocessor.BatchProcessorOpts[eventsub.ChannelPointsCustomRewardRedemptionAddEvent]{
			Interval:  500 * time.Millisecond,
			BatchSize: 100,
			Callback:  myHandler.handleChannelPointsRewardRedemptionAddBatched,
		},
	)

	lc.Append(
		lifecycle.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					myHandler.redemptionsBatcher.Start(batcherCtx)
				}()

				logger.Info("Handler started")

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
