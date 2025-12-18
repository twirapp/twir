package messagehandler

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/redis/go-redis/v9"
	batchprocessor "github.com/twirapp/batch-processor"
	"github.com/twirapp/twir/apps/bots/internal/moderationhelpers"
	chattranslationsservice "github.com/twirapp/twir/apps/bots/internal/services/chat_translations"
	"github.com/twirapp/twir/apps/bots/internal/services/giveaways"
	"github.com/twirapp/twir/apps/bots/internal/services/keywords"
	"github.com/twirapp/twir/apps/bots/internal/services/tts"
	"github.com/twirapp/twir/apps/bots/internal/services/voteban"
	"github.com/twirapp/twir/apps/bots/internal/twitchactions"
	"github.com/twirapp/twir/apps/bots/internal/workers"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/logger"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	"github.com/twirapp/twir/libs/repositories/channels_emotes_usages"
	channelsmoderationsettingsmodel "github.com/twirapp/twir/libs/repositories/channels_moderation_settings/model"
	"github.com/twirapp/twir/libs/repositories/chat_messages"
	chatwallrepository "github.com/twirapp/twir/libs/repositories/chat_wall"
	chatwallmodel "github.com/twirapp/twir/libs/repositories/chat_wall/model"
	giveawaysmodel "github.com/twirapp/twir/libs/repositories/giveaways/model"
	"github.com/twirapp/twir/libs/repositories/greetings"
	greetingsmodel "github.com/twirapp/twir/libs/repositories/greetings/model"
	"github.com/twirapp/twir/libs/repositories/users"
	usersstats "github.com/twirapp/twir/libs/repositories/users_stats"
	"github.com/twirapp/twir/libs/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Logger                           *slog.Logger
	WebsocketsGrpc                   websockets.WebsocketClient
	GreetingsRepository              greetings.Repository
	ChatMessagesRepository           chat_messages.Repository
	ChannelsEmotesUsagesRepository   channels_emotes_usages.Repository
	UsersstatsRepository             usersstats.Repository
	UsersRepository                  users.Repository
	Gorm                             *gorm.DB
	Redis                            *redis.Client
	TwitchActions                    *twitchactions.TwitchActions
	ModerationHelpers                *moderationhelpers.ModerationHelpers
	Bus                              *buscore.Bus
	KeywordsService                  *keywords.Service
	GreetingsCache                   *generic_cacher.GenericCacher[[]greetingsmodel.Greeting]
	TTSService                       *tts.Service
	Config                           cfg.Config
	ChatWallCacher                   *generic_cacher.GenericCacher[[]chatwallmodel.ChatWall]
	ChatWallRepository               chatwallrepository.Repository
	ChatWallSettingsCacher           *generic_cacher.GenericCacher[chatwallmodel.ChatWallSettings]
	ChannelsRepository               channelsrepository.Repository
	GiveawaysCacher                  *generic_cacher.GenericCacher[[]giveawaysmodel.ChannelGiveaway]
	ChannelsModerationSettingsCacher *generic_cacher.GenericCacher[[]channelsmoderationsettingsmodel.ChannelModerationSettings]
	VotebanService                   *voteban.Service
	ChatTranslatorService            *chattranslationsservice.Service
	GiveawaysService                 *giveaways.Service

	TrmManager trm.Manager

	WorkersPool *workers.Pool
}

type MessageHandler struct {
	logger                           *slog.Logger
	websocketsGrpc                   websockets.WebsocketClient
	greetingsRepository              greetings.Repository
	chatMessagesRepository           chat_messages.Repository
	channelsEmotesUsagesRepository   channels_emotes_usages.Repository
	usersstatsRepository             usersstats.Repository
	usersRepository                  users.Repository
	gorm                             *gorm.DB
	redis                            *redis.Client
	twitchActions                    *twitchactions.TwitchActions
	moderationHelpers                *moderationhelpers.ModerationHelpers
	twirBus                          *buscore.Bus
	greetingsCache                   *generic_cacher.GenericCacher[[]greetingsmodel.Greeting]
	chatWallCacher                   *generic_cacher.GenericCacher[[]chatwallmodel.ChatWall]
	chatWallRepository               chatwallrepository.Repository
	chatWallSettingsCacher           *generic_cacher.GenericCacher[chatwallmodel.ChatWallSettings]
	giveawaysCacher                  *generic_cacher.GenericCacher[[]giveawaysmodel.ChannelGiveaway]
	channelsModerationSettingsCacher *generic_cacher.GenericCacher[[]channelsmoderationsettingsmodel.ChannelModerationSettings]
	votebanService                   *voteban.Service
	chatTranslatorService            *chattranslationsservice.Service
	giveawaysService                 *giveaways.Service

	keywordsService *keywords.Service
	ttsService      *tts.Service
	config          cfg.Config
	workersPool     *workers.Pool
	trmManager      trm.Manager

	messagesSaveBatcher    *batchprocessor.BatchProcessor[twitch.TwitchChatMessage]
	messagesLurkersBatcher *batchprocessor.BatchProcessor[twitch.TwitchChatMessage]
	messagesEmotesBatcher  *batchprocessor.BatchProcessor[twitch.TwitchChatMessage]
}

var messageHandlerTracer = otel.Tracer("message-handler")

var handlersForExecute = []func(
	c *MessageHandler,
	ctx context.Context,
	msg twitch.TwitchChatMessage,
) error{
	(*MessageHandler).handleSaveMessage,
	(*MessageHandler).handleIncrementStreamMessages,
	(*MessageHandler).handleGreetings,
	(*MessageHandler).handleKeywords,
	(*MessageHandler).handleEmotesUsages,
	(*MessageHandler).handleTts,
	(*MessageHandler).handleRemoveLurker,
	(*MessageHandler).handleModeration,
	(*MessageHandler).handleFirstStreamUserJoin,
	(*MessageHandler).handleChatWall,
	(*MessageHandler).handleGiveaways,
}

func New(opts Opts) *MessageHandler {
	handler := &MessageHandler{
		logger:            opts.Logger,
		gorm:              opts.Gorm,
		redis:             opts.Redis,
		twitchActions:     opts.TwitchActions,
		websocketsGrpc:    opts.WebsocketsGrpc,
		moderationHelpers: opts.ModerationHelpers,
		config:            opts.Config,
		twirBus:           opts.Bus,

		keywordsService:                  opts.KeywordsService,
		greetingsRepository:              opts.GreetingsRepository,
		chatMessagesRepository:           opts.ChatMessagesRepository,
		channelsEmotesUsagesRepository:   opts.ChannelsEmotesUsagesRepository,
		usersstatsRepository:             opts.UsersstatsRepository,
		greetingsCache:                   opts.GreetingsCache,
		ttsService:                       opts.TTSService,
		chatWallCacher:                   opts.ChatWallCacher,
		chatWallRepository:               opts.ChatWallRepository,
		chatWallSettingsCacher:           opts.ChatWallSettingsCacher,
		giveawaysCacher:                  opts.GiveawaysCacher,
		channelsModerationSettingsCacher: opts.ChannelsModerationSettingsCacher,
		trmManager:                       opts.TrmManager,
		usersRepository:                  opts.UsersRepository,
		votebanService:                   opts.VotebanService,
		chatTranslatorService:            opts.ChatTranslatorService,
		giveawaysService:                 opts.GiveawaysService,

		workersPool: opts.WorkersPool,
	}

	batcherCtx, batcherCancel := context.WithCancel(context.Background())

	handler.messagesSaveBatcher = batchprocessor.NewBatchProcessor[twitch.TwitchChatMessage](
		batchprocessor.BatchProcessorOpts[twitch.TwitchChatMessage]{
			Interval:  500 * time.Millisecond,
			BatchSize: 1000,
			Callback:  handler.handleSaveMessageBatched,
		},
	)
	handler.messagesLurkersBatcher = batchprocessor.NewBatchProcessor[twitch.TwitchChatMessage](
		batchprocessor.BatchProcessorOpts[twitch.TwitchChatMessage]{
			Interval:  100 * time.Millisecond,
			BatchSize: 100,
			Callback:  handler.handleRemoveLurkerBatched,
		},
	)
	handler.messagesEmotesBatcher = batchprocessor.NewBatchProcessor[twitch.TwitchChatMessage](
		batchprocessor.BatchProcessorOpts[twitch.TwitchChatMessage]{
			Interval:  500 * time.Millisecond,
			BatchSize: 1000,
			Callback:  handler.handleEmotesUsagesBatched,
		},
	)

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					handler.messagesSaveBatcher.Start(batcherCtx)
				}()

				go func() {
					handler.messagesEmotesBatcher.Start(batcherCtx)
				}()

				go func() {
					handler.messagesLurkersBatcher.Start(batcherCtx)
				}()

				return nil
			},
			OnStop: func(ctx context.Context) error {
				if err := handler.messagesSaveBatcher.Shutdown(ctx); err != nil {
					return err
				}

				if err := handler.messagesLurkersBatcher.Shutdown(ctx); err != nil {
					return err
				}

				if err := handler.messagesEmotesBatcher.Shutdown(ctx); err != nil {
					return err
				}

				batcherCancel()
				return nil
			},
		},
	)

	handlersForExecute = append(
		handlersForExecute,
		func(c *MessageHandler, ctx context.Context, msg twitch.TwitchChatMessage) error {
			c.votebanService.TryRegisterVote(msg)
			return nil
		},
		func(c *MessageHandler, ctx context.Context, msg twitch.TwitchChatMessage) error {
			return c.chatTranslatorService.Handle(ctx, msg)
		},
	)

	return handler
}

func (c *MessageHandler) Handle(ctx context.Context, req twitch.TwitchChatMessage) error {
	newCtx, span := messageHandlerTracer.Start(ctx, "handle")
	ctx = newCtx
	defer span.End()

	span.SetAttributes(
		attribute.String("function.name", utils.GetFuncName()),
		attribute.String("channel.id", req.BroadcasterUserId),
		attribute.String("channel.login", req.BroadcasterUserLogin),
		attribute.String("user.id", req.ChatterUserId),
		attribute.String("user.login", req.ChatterUserLogin),
	)

	if !req.EnrichedData.DbChannel.IsEnabled {
		return nil
	}

	if req.EnrichedData.DbUser == nil {
		return fmt.Errorf("db user not found after ensureUser")
	}

	if req.ChatterUserId == req.EnrichedData.DbChannel.BotID && c.config.AppEnv == "production" {
		return nil
	}

	// tasks will be stopped if context is canceled
	handleTask := c.workersPool.NewGroupContext(ctx)

	for _, handlerFunc := range handlersForExecute {
		funcName := runtime.FuncForPC(reflect.ValueOf(handlerFunc).Pointer()).Name()
		splitName := strings.Split(funcName, ".")
		shortFuncName := splitName[len(splitName)-1]

		handleTask.Submit(
			func() {
				handlerCtx, handlerSpan := messageHandlerTracer.Start(
					ctx, shortFuncName, trace.WithAttributes(
						attribute.String("handler.name", funcName),
					),
				)
				defer handlerSpan.End()

				handlerError := handlerFunc(c, handlerCtx, req)
				if handlerError != nil {
					handlerSpan.RecordError(handlerError)
					c.logger.Error(
						"error executing handler",
						slog.String("shortFuncName", shortFuncName),
						logger.Error(handlerError),
					)
				}
			},
		)
	}

	if err := handleTask.Wait(); err != nil {
		return err
	}

	return nil
}
