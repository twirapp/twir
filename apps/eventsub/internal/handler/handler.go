package handler

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	batchprocessor "github.com/twirapp/batch-processor"
	bus_core "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	cfg "github.com/twirapp/twir/libs/config"
	deprecatedmodel "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/logger"
	channelmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	channelscommandsprefixmodel "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	channelsinfohistory "github.com/twirapp/twir/libs/repositories/channels_info_history"
	channelsredemptionshistory "github.com/twirapp/twir/libs/repositories/channels_redemptions_history"
	commandswithgroupsandresponsesmodel "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"
	scheduledvipsrepository "github.com/twirapp/twir/libs/repositories/scheduled_vips"
	"github.com/twirapp/twir/libs/repositories/streams"
	"github.com/twirapp/twitchy/eventsub"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"gorm.io/gorm"

	alertmodel "github.com/twirapp/twir/libs/repositories/alerts/model"
)

type Handler struct {
	logger logger.Logger

	websocketsGrpc               websockets.WebsocketClient
	tracer                       trace.Tracer
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
	commandsCache                       *generic_cacher.GenericCacher[[]commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses]
	channelSongRequestsSettingsCache    *generic_cacher.GenericCacher[deprecatedmodel.ChannelSongRequestsSettings]
	channelsIntegrationsSettingsSeventv *generic_cacher.GenericCacher[deprecatedmodel.ChannelsIntegrationsSettingsSeventv]
	config                              cfg.Config

	redemptionsBatcher *batchprocessor.BatchProcessor[eventsub.ChannelPointsCustomRewardRedemptionAddEvent]
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
	CommandsCache                       *generic_cacher.GenericCacher[[]commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses]
	ChannelSongRequestsSettingsCache    *generic_cacher.GenericCacher[deprecatedmodel.ChannelSongRequestsSettings]
	ChannelsIntegrationsSettingsSeventv *generic_cacher.GenericCacher[deprecatedmodel.ChannelsIntegrationsSettingsSeventv]

	Tracer trace.Tracer
	Gorm   *gorm.DB
	Redis  *redis.Client

	Bus                *bus_core.Bus
	PrefixCache        *generic_cacher.GenericCacher[channelscommandsprefixmodel.ChannelsCommandsPrefix]
	ChannelAlertsCache *generic_cacher.GenericCacher[[]alertmodel.Alert]

	Config cfg.Config
}

func New(opts Opts) *Handler {
	myHandler := &Handler{
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

	myHandler.redemptionsBatcher = batchprocessor.NewBatchProcessor[eventsub.ChannelPointsCustomRewardRedemptionAddEvent](
		batchprocessor.BatchProcessorOpts[eventsub.ChannelPointsCustomRewardRedemptionAddEvent]{
			Interval:  500 * time.Millisecond,
			BatchSize: 100,
			Callback:  myHandler.handleChannelPointsRewardRedemptionAddBatched,
		},
	)

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
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
