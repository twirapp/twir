package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/kr/pretty"
	"github.com/nicklaw5/helix/v2"
	cfg "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/twitch"
	"github.com/satont/twir/libs/utils"
	"github.com/twirapp/twir/libs/grpc/bots"
	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/shared"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
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

	TokensGrpc tokens.TokensClient
}

func New(opts Opts) (*Grpc, error) {
	impl := &Grpc{
		gorm:   opts.Gorm,
		logger: opts.Logger,
	}

	grpcNetListener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", constants.BOTS_SERVER_PORT))
	if err != nil {
		return nil, err
	}
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				MaxConnectionAge: 1 * time.Minute,
			},
		),
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

	gorm       *gorm.DB
	logger     logger.Logger
	config     cfg.Config
	tokensGrpc tokens.TokensClient
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

	twitchClient, err := twitch.NewBotClientWithContext(ctx, channel.BotID, c.config, c.tokensGrpc)
	if err != nil {
		return nil, err
	}

	wg := utils.NewGoroutinesGroup()

	for _, m := range req.GetMessageIds() {
		wg.Go(
			func() {
				request, e := twitchClient.DeleteChatMessage(
					&helix.DeleteChatMessageParams{
						BroadcasterID: channel.ID,
						ModeratorID:   channel.BotID,
						MessageID:     m,
					},
				)
				if e != nil {
					c.logger.Error("cannot delete message", slog.Any("err", e))
				} else if request.ErrorMessage != "" {
					c.logger.Error("cannot delete message", slog.String("err", request.ErrorMessage))
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
	pretty.Println(req.GetMessage())
	return &emptypb.Empty{}, nil
}

func (c *Grpc) HandleChatMessage(ctx context.Context, req *shared.TwitchChatMessage) (
	*emptypb.Empty,
	error,
) {
	pretty.Println(req.Message.GetFragments())
	return &emptypb.Empty{}, nil
}
