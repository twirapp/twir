package messagehandler

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"runtime"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/bots/internal/moderationhelpers"
	"github.com/satont/twir/apps/bots/internal/services/keywords"
	"github.com/satont/twir/apps/bots/internal/services/tts"
	"github.com/satont/twir/apps/bots/internal/twitchactions"
	"github.com/satont/twir/apps/bots/internal/workers"
	cfg "github.com/satont/twir/libs/config"
	deprecatedgormmodel "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	batchprocessor "github.com/twirapp/batch-processor"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/websockets"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	"github.com/twirapp/twir/libs/repositories/channels_emotes_usages"
	channelsmoderationsettingsmodel "github.com/twirapp/twir/libs/repositories/channels_moderation_settings/model"
	"github.com/twirapp/twir/libs/repositories/chat_messages"
	chatwallrepository "github.com/twirapp/twir/libs/repositories/chat_wall"
	chatwallmodel "github.com/twirapp/twir/libs/repositories/chat_wall/model"
	giveawaysmodel "github.com/twirapp/twir/libs/repositories/giveaways/model"
	"github.com/twirapp/twir/libs/repositories/greetings"
	greetingsmodel "github.com/twirapp/twir/libs/repositories/greetings/model"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Logger                           logger.Logger
	ParserGrpc                       parser.ParserClient
	WebsocketsGrpc                   websockets.WebsocketClient
	GreetingsRepository              greetings.Repository
	ChatMessagesRepository           chat_messages.Repository
	ChannelsEmotesUsagesRepository   channels_emotes_usages.Repository
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

	WorkersPool *workers.Pool
}

type MessageHandler struct {
	logger                           logger.Logger
	parserGrpc                       parser.ParserClient
	websocketsGrpc                   websockets.WebsocketClient
	greetingsRepository              greetings.Repository
	chatMessagesRepository           chat_messages.Repository
	channelsEmotesUsagesRepository   channels_emotes_usages.Repository
	gorm                             *gorm.DB
	redis                            *redis.Client
	twitchActions                    *twitchactions.TwitchActions
	moderationHelpers                *moderationhelpers.ModerationHelpers
	twirBus                          *buscore.Bus
	votebanMutex                     *redsync.Mutex
	greetingsCache                   *generic_cacher.GenericCacher[[]greetingsmodel.Greeting]
	chatWallCacher                   *generic_cacher.GenericCacher[[]chatwallmodel.ChatWall]
	chatWallRepository               chatwallrepository.Repository
	chatWallSettingsCacher           *generic_cacher.GenericCacher[chatwallmodel.ChatWallSettings]
	giveawaysCacher                  *generic_cacher.GenericCacher[[]giveawaysmodel.ChannelGiveaway]
	channelsModerationSettingsCacher *generic_cacher.GenericCacher[[]channelsmoderationsettingsmodel.ChannelModerationSettings]

	keywordsService *keywords.Service
	ttsService      *tts.Service
	config          cfg.Config
	workersPool     *workers.Pool

	messagesSaveBatcher    *batchprocessor.BatchProcessor[handleMessage]
	messagesLurkersBatcher *batchprocessor.BatchProcessor[handleMessage]
	messagesEmotesBatcher  *batchprocessor.BatchProcessor[handleMessage]
}

func New(opts Opts) *MessageHandler {
	votebanLock := redsync.New(goredis.NewPool(opts.Redis))

	handler := &MessageHandler{
		logger:                           opts.Logger,
		gorm:                             opts.Gorm,
		redis:                            opts.Redis,
		twitchActions:                    opts.TwitchActions,
		parserGrpc:                       opts.ParserGrpc,
		websocketsGrpc:                   opts.WebsocketsGrpc,
		moderationHelpers:                opts.ModerationHelpers,
		config:                           opts.Config,
		twirBus:                          opts.Bus,
		votebanMutex:                     votebanLock.NewMutex("bots:voteban_handle_message"),
		keywordsService:                  opts.KeywordsService,
		greetingsRepository:              opts.GreetingsRepository,
		chatMessagesRepository:           opts.ChatMessagesRepository,
		channelsEmotesUsagesRepository:   opts.ChannelsEmotesUsagesRepository,
		greetingsCache:                   opts.GreetingsCache,
		ttsService:                       opts.TTSService,
		chatWallCacher:                   opts.ChatWallCacher,
		chatWallRepository:               opts.ChatWallRepository,
		chatWallSettingsCacher:           opts.ChatWallSettingsCacher,
		giveawaysCacher:                  opts.GiveawaysCacher,
		channelsModerationSettingsCacher: opts.ChannelsModerationSettingsCacher,

		workersPool: opts.WorkersPool,
	}

	batcherCtx, batcherCancel := context.WithCancel(context.Background())

	handler.messagesSaveBatcher = batchprocessor.NewBatchProcessor[handleMessage](
		batchprocessor.BatchProcessorOpts[handleMessage]{
			Interval:  500 * time.Millisecond,
			BatchSize: 1000,
			Callback:  handler.handleSaveMessageBatched,
		},
	)
	handler.messagesLurkersBatcher = batchprocessor.NewBatchProcessor[handleMessage](
		batchprocessor.BatchProcessorOpts[handleMessage]{
			Interval:  100 * time.Millisecond,
			BatchSize: 100,
			Callback:  handler.handleRemoveLurkerBatched,
		},
	)
	handler.messagesEmotesBatcher = batchprocessor.NewBatchProcessor[handleMessage](
		batchprocessor.BatchProcessorOpts[handleMessage]{
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

	return handler
}

type handleMessage struct {
	DbUser *deprecatedgormmodel.Users
	twitch.TwitchChatMessage
}

var handlersForExecute = []func(
	c *MessageHandler,
	ctx context.Context,
	msg handleMessage,
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
	(*MessageHandler).handleGamesVoteban,
	(*MessageHandler).handleChatWall,
	(*MessageHandler).handleGiveaways,
}

func (c *MessageHandler) Handle(ctx context.Context, req twitch.TwitchChatMessage) error {
	msg := handleMessage{
		TwitchChatMessage: req,
	}

	if !msg.EnrichedData.DbChannel.IsEnabled {
		return nil
	}

	dbUser, err := c.ensureUser(ctx, msg)
	if err != nil {
		return err
	}
	msg.DbUser = dbUser

	if req.ChatterUserId == msg.EnrichedData.DbChannel.BotID && c.config.AppEnv == "production" {
		return nil
	}

	// tasks will be stopped if context is canceled
	handleTask := c.workersPool.NewGroupContext(ctx)

	for _, f := range handlersForExecute {
		handleTask.SubmitErr(
			func() error {
				handlerError := f(c, ctx, msg)
				if handlerError != nil {
					return fmt.Errorf(
						"error executing %s handler: %w",
						runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(),
						handlerError,
					)
				}

				return nil
			},
		)
	}

	if err := handleTask.Wait(); err != nil {
		c.logger.Error("error on execution all handlers", slog.Any("err", err))
		return err
	}

	return nil
}
