package messagehandler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"runtime"
	"sync"

	"github.com/alitto/pond"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/bots/internal/moderationhelpers"
	"github.com/satont/twir/apps/bots/internal/twitchactions"
	cfg "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/shared"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Logger            logger.Logger
	Gorm              *gorm.DB
	Redis             *redis.Client
	TwitchActions     *twitchactions.TwitchActions
	ParserGrpc        parser.ParserClient
	WebsocketsGrpc    websockets.WebsocketClient
	EventsGrpc        events.EventsClient
	ModerationHelpers *moderationhelpers.ModerationHelpers
	Config            cfg.Config
}

type MessageHandler struct {
	logger            logger.Logger
	gorm              *gorm.DB
	redis             *redis.Client
	pool              *pond.WorkerPool
	twitchActions     *twitchactions.TwitchActions
	parserGrpc        parser.ParserClient
	websocketsGrpc    websockets.WebsocketClient
	eventsGrpc        events.EventsClient
	moderationHelpers *moderationhelpers.ModerationHelpers
	config            cfg.Config
}

func New(opts Opts) *MessageHandler {
	pool := pond.New(
		10,
		1000,
		pond.Strategy(pond.Balanced()),
		pond.PanicHandler(
			func(i interface{}) {
				opts.Logger.Error("panic", slog.Any("err", i))
			},
		),
	)
	return &MessageHandler{
		logger:            opts.Logger,
		gorm:              opts.Gorm,
		redis:             opts.Redis,
		pool:              pool,
		twitchActions:     opts.TwitchActions,
		parserGrpc:        opts.ParserGrpc,
		websocketsGrpc:    opts.WebsocketsGrpc,
		eventsGrpc:        opts.EventsGrpc,
		moderationHelpers: opts.ModerationHelpers,
		config:            opts.Config,
	}
}

type handleMessage struct {
	*shared.TwitchChatMessage
	DbChannel *model.Channels
	DbStream  *model.ChannelsStreams
	DbUser    *model.Users
}

func (c *MessageHandler) Handle(ctx context.Context, req *shared.TwitchChatMessage) error {
	msg := handleMessage{
		TwitchChatMessage: req,
	}

	errwg, errWgCtx := errgroup.WithContext(context.TODO())

	errwg.Go(
		func() error {
			stream := &model.ChannelsStreams{}
			if err := c.gorm.WithContext(errWgCtx).Where(
				`"userId" = ?`,
				req.GetBroadcasterUserId(),
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
			dbChannel := &model.Channels{}
			if err := c.gorm.WithContext(errWgCtx).Where(
				"id = ?",
				req.GetBroadcasterUserId(),
			).First(dbChannel).
				Error; err != nil {
				return err
			}
			msg.DbChannel = dbChannel
			return nil
		},
	)

	if err := errwg.Wait(); err != nil {
		return err
	}

	if !msg.DbChannel.IsEnabled {
		fmt.Println("channel not enabled", msg.DbChannel.ID)
		return nil
	}

	dbUser, err := c.ensureUser(ctx, msg)
	if err != nil {
		return err
	}
	msg.DbUser = dbUser

	if req.GetChatterUserId() == msg.DbChannel.BotID && c.config.AppEnv == "production" {
		fmt.Println("same bot user", req.GetChatterUserId(), msg.DbChannel.BotID)
		return nil
	}

	var wg sync.WaitGroup

	funcsForExecute := [...]func(ctx context.Context, msg handleMessage) error{
		c.handleCommand,
		c.handleIncrementStreamMessages,
		c.handleGreetings,
		c.handleKeywords,
		c.handleEmotesUsages,
		c.handleStoreMessage,
		c.handleTts,
		c.handleRemoveLurker,
		c.handleModeration,
		c.handleFirstStreamUserJoin,
	}

	// TODO: i dont know why grpc context canceling before this function finished
	funcsCtx := context.Background()

	for _, f := range funcsForExecute {
		wg.Add(1)

		f := f

		c.pool.Submit(
			func() {
				if err := f(funcsCtx, msg); err != nil {
					c.logger.Error(
						"error when executing message handler function", slog.Any("err", err),
						slog.String("functionName", runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()),
					)
				}
				wg.Done()
			},
		)
	}

	wg.Wait()

	return nil
}
