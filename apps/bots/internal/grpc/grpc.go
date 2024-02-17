package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/satont/twir/apps/bots/internal/messagehandler"
	"github.com/satont/twir/apps/bots/internal/twitchactions"
	cfg "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/utils"
	"github.com/twirapp/twir/libs/grpc/bots"
	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/shared"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In
	Lc fx.In

	Gorm   *gorm.DB
	Logger logger.Logger
	Cfg    cfg.Config
	LC     fx.Lifecycle

	TokensGrpc     tokens.TokensClient
	TwitchActions  *twitchactions.TwitchActions
	MessageHandler *messagehandler.MessageHandler
	Tracer         trace.Tracer
}

func New(opts Opts) (*Grpc, error) {
	impl := &Grpc{
		gorm:           opts.Gorm,
		logger:         opts.Logger,
		config:         opts.Cfg,
		tokensGrpc:     opts.TokensGrpc,
		twitchactions:  opts.TwitchActions,
		messageHandler: opts.MessageHandler,
		tracer:         opts.Tracer,
	}

	grpcNetListener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", constants.BOTS_SERVER_PORT))
	if err != nil {
		return nil, err
	}
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	bots.RegisterBotsServer(grpcServer, impl)

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					if err := grpcServer.Serve(grpcNetListener); err != nil {
						panic(err)
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return nil
			},
		},
	)

	return impl, nil
}

type Grpc struct {
	bots.UnimplementedBotsServer

	gorm           *gorm.DB
	logger         logger.Logger
	config         cfg.Config
	tokensGrpc     tokens.TokensClient
	twitchactions  *twitchactions.TwitchActions
	messageHandler *messagehandler.MessageHandler
	tracer         trace.Tracer
}

var _ bots.BotsServer = (*Grpc)(nil)

func (c *Grpc) DeleteMessage(ctx context.Context, req *bots.DeleteMessagesRequest) (
	*emptypb.Empty,
	error,
) {
	channel := model.Channels{}
	err := c.gorm.WithContext(ctx).Where("id = ?", req.GetChannelId()).Find(&channel).Error
	if err != nil {
		c.logger.Error(
			"cannot get channel",
			slog.String("channelId", req.GetChannelId()),
			slog.String("channelName", req.GetChannelName()),
		)
		return &emptypb.Empty{}, nil
	}

	if channel.ID == "" {
		return &emptypb.Empty{}, nil
	}

	wg := utils.NewGoroutinesGroup()

	for _, m := range req.GetMessageIds() {
		wg.Go(
			func() {
				e := c.twitchactions.DeleteMessage(
					ctx,
					twitchactions.DeleteMessageOpts{
						BroadcasterID: req.GetChannelId(),
						ModeratorID:   channel.BotID,
						MessageID:     m,
					},
				)
				if e != nil {
					c.logger.Error("cannot delete message", slog.Any("err", e))
				}
			},
		)
	}

	wg.Wait()

	return &emptypb.Empty{}, nil
}

func (c *Grpc) SendMessage(ctx context.Context, req *bots.SendMessageRequest) (
	*emptypb.Empty,
	error,
) {
	channel := model.Channels{}
	err := c.gorm.WithContext(ctx).Where("id = ?", req.GetChannelId()).Find(&channel).Error
	if err != nil {
		c.logger.Error(
			"cannot get channel",
			slog.String("channelId", req.GetChannelId()),
			slog.String("channelName", req.GetChannelName()),
		)
		return &emptypb.Empty{}, nil
	}

	if channel.ID == "" {
		return &emptypb.Empty{}, nil
	}

	err = c.twitchactions.SendMessage(
		ctx,
		twitchactions.SendMessageOpts{
			BroadcasterID:        req.GetChannelId(),
			SenderID:             channel.BotID,
			Message:              req.GetMessage(),
			ReplyParentMessageID: req.GetReplyTo(),
			IsAnnounce:           req.GetIsAnnounce(),
		},
	)
	if err != nil {
		c.logger.Error("cannot send message", slog.Any("err", err))
	}
	return &emptypb.Empty{}, nil
}

func (c *Grpc) HandleChatMessage(ctx context.Context, req *shared.TwitchChatMessage) (
	*emptypb.Empty,
	error,
) {
	span := trace.SpanFromContext(ctx)
	// End the span when the operation we are measuring is done.
	defer span.End()
	span.SetAttributes(
		attribute.String("message_id", req.GetMessageId()),
		attribute.String("channel_id", req.GetBroadcasterUserId()),
	)

	err := c.messageHandler.Handle(ctx, req)
	if err != nil {
		c.logger.Error(
			"cannot handle message",
			slog.String("channelId", req.GetBroadcasterUserId()),
			slog.String("channelName", req.GetBroadcasterUserLogin()),
			slog.Any("err", err),
		)
	}

	return &emptypb.Empty{}, nil
}
