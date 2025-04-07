package messagehandler

import (
	"context"
	"errors"
	"log/slog"

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
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/websockets"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	"github.com/twirapp/twir/libs/repositories/chat_messages"
	chatwallrepository "github.com/twirapp/twir/libs/repositories/chat_wall"
	chatwallmodel "github.com/twirapp/twir/libs/repositories/chat_wall/model"
	"github.com/twirapp/twir/libs/repositories/greetings"
	greetingsmodel "github.com/twirapp/twir/libs/repositories/greetings/model"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Logger                 logger.Logger
	ParserGrpc             parser.ParserClient
	WebsocketsGrpc         websockets.WebsocketClient
	EventsGrpc             events.EventsClient
	GreetingsRepository    greetings.Repository
	ChatMessagesRepository chat_messages.Repository
	Gorm                   *gorm.DB
	Redis                  *redis.Client
	TwitchActions          *twitchactions.TwitchActions
	ModerationHelpers      *moderationhelpers.ModerationHelpers
	Bus                    *buscore.Bus
	KeywordsService        *keywords.Service
	GreetingsCache         *generic_cacher.GenericCacher[[]greetingsmodel.Greeting]
	TTSService             *tts.Service
	Config                 cfg.Config
	ChatWallCacher         *generic_cacher.GenericCacher[[]chatwallmodel.ChatWall]
	ChatWallRepository     chatwallrepository.Repository
	ChatWallSettingsCacher *generic_cacher.GenericCacher[chatwallmodel.ChatWallSettings]
	ChannelsRepository     channelsrepository.Repository

	WorkersPool *workers.Pool
}

type MessageHandler struct {
	logger                 logger.Logger
	parserGrpc             parser.ParserClient
	websocketsGrpc         websockets.WebsocketClient
	eventsGrpc             events.EventsClient
	greetingsRepository    greetings.Repository
	chatMessagesRepository chat_messages.Repository
	gorm                   *gorm.DB
	redis                  *redis.Client
	twitchActions          *twitchactions.TwitchActions
	moderationHelpers      *moderationhelpers.ModerationHelpers
	bus                    *buscore.Bus
	votebanMutex           *redsync.Mutex
	greetingsCache         *generic_cacher.GenericCacher[[]greetingsmodel.Greeting]
	chatWallCacher         *generic_cacher.GenericCacher[[]chatwallmodel.ChatWall]
	chatWallRepository     chatwallrepository.Repository
	chatWallSettingsCacher *generic_cacher.GenericCacher[chatwallmodel.ChatWallSettings]

	keywordsService *keywords.Service
	ttsService      *tts.Service
	config          cfg.Config
	workersPool     *workers.Pool
}

func New(opts Opts) *MessageHandler {
	votebanLock := redsync.New(goredis.NewPool(opts.Redis))

	handler := &MessageHandler{
		logger:                 opts.Logger,
		gorm:                   opts.Gorm,
		redis:                  opts.Redis,
		twitchActions:          opts.TwitchActions,
		parserGrpc:             opts.ParserGrpc,
		websocketsGrpc:         opts.WebsocketsGrpc,
		eventsGrpc:             opts.EventsGrpc,
		moderationHelpers:      opts.ModerationHelpers,
		config:                 opts.Config,
		bus:                    opts.Bus,
		votebanMutex:           votebanLock.NewMutex("bots:voteban_handle_message"),
		keywordsService:        opts.KeywordsService,
		greetingsRepository:    opts.GreetingsRepository,
		chatMessagesRepository: opts.ChatMessagesRepository,
		greetingsCache:         opts.GreetingsCache,
		ttsService:             opts.TTSService,
		chatWallCacher:         opts.ChatWallCacher,
		chatWallRepository:     opts.ChatWallRepository,
		chatWallSettingsCacher: opts.ChatWallSettingsCacher,
		workersPool:            opts.WorkersPool,
	}

	return handler
}

type handleMessage struct {
	DbStream *deprecatedgormmodel.ChannelsStreams
	DbUser   *deprecatedgormmodel.Users
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
}

func (c *MessageHandler) Handle(ctx context.Context, req twitch.TwitchChatMessage) error {
	msg := handleMessage{
		TwitchChatMessage: req,
	}

	var errwg errgroup.Group
	errwg.Go(
		func() error {
			stream := &deprecatedgormmodel.ChannelsStreams{}
			if err := c.gorm.WithContext(ctx).Where(
				`"userId" = ?`,
				req.BroadcasterUserId,
			).First(stream).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if stream.ID == "" {
				msg.DbStream = nil
			} else {
				msg.DbStream = stream
			}
			return nil
		},
	)

	if err := errwg.Wait(); err != nil {
		return err
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
				return f(c, ctx, msg)
			},
		)
	}

	if err := handleTask.Wait(); err != nil {
		c.logger.Error("error on execution all handlers", slog.Any("err", err))
		return err
	}

	return nil
}
