package messagehandler

import (
	"context"
	"log/slog"

	"github.com/alitto/pond"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/bots/internal/twitchactions"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/shared"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Logger        logger.Logger
	Gorm          *gorm.DB
	Redis         *redis.Client
	TwitchActions *twitchactions.TwitchActions
	ParserGrpc    parser.ParserClient
}

func New(opts Opts) *MessageHandler {
	pool := pond.New(
		10,
		1000,
		pond.Strategy(pond.Balanced()),
		pond.PanicHandler(
			func(i interface{}) {
				opts.Logger.Error("paniced", slog.Any("err", i))
			},
		),
	)
	return &MessageHandler{
		logger:        opts.Logger,
		gorm:          opts.Gorm,
		redis:         opts.Redis,
		pool:          pool,
		twitchActions: opts.TwitchActions,
		parserGrpc:    opts.ParserGrpc,
	}
}

type MessageHandler struct {
	logger        logger.Logger
	gorm          *gorm.DB
	redis         *redis.Client
	pool          *pond.WorkerPool
	twitchActions *twitchactions.TwitchActions
	parserGrpc    parser.ParserClient
}

type handleMessage struct {
	*shared.TwitchChatMessage
	DbChannel *model.Channels
	DbStream  *model.ChannelsStreams
}

func (c *MessageHandler) Handle(ctx context.Context, req *shared.TwitchChatMessage) error {
	msg := handleMessage{
		TwitchChatMessage: req,
	}

	errwg, errWgCtx := errgroup.WithContext(ctx)

	errwg.Go(
		func() error {
			stream := &model.ChannelsStreams{}
			if err := c.gorm.WithContext(errWgCtx).Where(
				`"userId" = ?`,
				req.GetBroadcasterUserId(),
			).Find(stream).Error; err != nil {
				return err
			}
			msg.DbStream = stream
			return nil
		},
	)

	errwg.Go(
		func() error {
			dbChannel := &model.Channels{}
			if err := c.gorm.WithContext(errWgCtx).Where(
				"id = ?",
				req.GetBroadcasterUserId(),
			).Find(dbChannel).
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
		return nil
	}

	c.pool.SubmitAndWait(
		func() {
			c.handleCommand(context.TODO(), msg)
		},
	)

	return nil
}
