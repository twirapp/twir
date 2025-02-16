package messagehandler

import (
	"context"
	"errors"
	"log/slog"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/bots/internal/moderationhelpers"
	"github.com/satont/twir/apps/bots/internal/services/commands"
	"github.com/satont/twir/apps/bots/internal/services/keywords"
	"github.com/satont/twir/apps/bots/internal/services/tts"
	"github.com/satont/twir/apps/bots/internal/twitchactions"
	cfg "github.com/satont/twir/libs/config"
	deprecatedgormmodel "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/repositories/chat_messages"
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
	CommandService         *commands.Service
	TTSService             *tts.Service
	Config                 cfg.Config
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

	keywordsService *keywords.Service
	commandsService *commands.Service
	ttsService      *tts.Service
	config          cfg.Config
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
		commandsService:        opts.CommandService,
		ttsService:             opts.TTSService,
	}

	return handler
}

type handleMessage struct {
	DbChannel *deprecatedgormmodel.Channels
	DbStream  *deprecatedgormmodel.ChannelsStreams
	DbUser    *deprecatedgormmodel.Users
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
}

func (c *MessageHandler) Handle(ctx context.Context, req twitch.TwitchChatMessage) error {
	msg := handleMessage{
		TwitchChatMessage: req,
	}

	errwg, errWgCtx := errgroup.WithContext(context.TODO())

	errwg.Go(
		func() error {
			stream := &deprecatedgormmodel.ChannelsStreams{}
			if err := c.gorm.WithContext(errWgCtx).Where(
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

	errwg.Go(
		func() error {
			cacheKey := "cache:bots:channels:" + req.BroadcasterUserId

			cachedData, err := c.redis.Get(ctx, cacheKey).Bytes()
			if err != nil && !errors.Is(err, redis.Nil) {
				return err
			}

			if len(cachedData) > 0 {
				dbChannel := &deprecatedgormmodel.Channels{}
				if err := json.Unmarshal(cachedData, dbChannel); err != nil {
					return err
				}
				msg.DbChannel = dbChannel
				return nil
			}

			dbChannel := &deprecatedgormmodel.Channels{}
			if err := c.gorm.WithContext(errWgCtx).Where(
				"id = ?",
				req.BroadcasterUserId,
			).First(dbChannel).
				Error; err != nil {
				return err
			}
			msg.DbChannel = dbChannel

			cacheBytes, err := json.Marshal(dbChannel)
			if err != nil {
				return err
			}

			if err := c.redis.Set(
				ctx,
				cacheKey,
				cacheBytes,
				5*time.Minute,
			).Err(); err != nil {
				return err
			}

			return nil
		},
	)

	if err := errwg.Wait(); err != nil {
		return err
	}

	if !msg.DbChannel.IsEnabled {
		return nil
	}

	dbUser, err := c.ensureUser(ctx, msg)
	if err != nil {
		return err
	}
	msg.DbUser = dbUser

	if req.ChatterUserId == msg.DbChannel.BotID && c.config.AppEnv == "production" {
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(len(handlersForExecute))

	// TODO: i dont know why grpc context canceling before this function finished
	funcsCtx := context.Background()

	for _, f := range handlersForExecute {
		f := f

		go func() {
			defer wg.Done()
			if err := f(c, funcsCtx, msg); err != nil {
				c.logger.Error(
					"error when executing message handler function",
					slog.Any("err", err),
					slog.String("functionName", runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()),
				)
			}
		}()
	}

	wg.Wait()

	return nil
}
