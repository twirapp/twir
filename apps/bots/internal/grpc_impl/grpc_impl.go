package grpc_impl

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strings"
	"time"

	"github.com/satont/twir/libs/grpc/servers"
	"github.com/satont/twir/libs/logger"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/satont/twir/libs/twitch"

	"github.com/nicklaw5/helix/v2"
	internalBots "github.com/satont/twir/apps/bots/internal/bots"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"

	twirUtils "github.com/satont/twir/libs/utils"
)

type GrpcImplOpts struct {
	fx.In

	Db          *gorm.DB
	BotsService *internalBots.Service
	Logger      logger.Logger
	Cfg         cfg.Config
	LC          fx.Lifecycle

	TokensGrpc tokens.TokensClient
}

type BotsGrpcServer struct {
	bots.UnimplementedBotsServer

	db          *gorm.DB
	botsService *internalBots.Service
	logger      logger.Logger
	cfg         cfg.Config

	tokensGrpc tokens.TokensClient
}

func NewServer(opts GrpcImplOpts) error {
	server := &BotsGrpcServer{
		db:          opts.Db,
		botsService: opts.BotsService,
		logger:      opts.Logger,
		cfg:         opts.Cfg,
		tokensGrpc:  opts.TokensGrpc,
	}

	grpcNetListener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", servers.BOTS_SERVER_PORT))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				MaxConnectionAge: 1 * time.Minute,
			},
		),
	)
	bots.RegisterBotsServer(grpcServer, server)

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go grpcServer.Serve(grpcNetListener)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				grpcServer.Stop()
				return nil
			},
		},
	)

	return nil
}

func (c *BotsGrpcServer) DeleteMessage(
	ctx context.Context,
	data *bots.DeleteMessagesRequest,
) (*emptypb.Empty, error) {
	channel := model.Channels{}
	err := c.db.WithContext(ctx).Where("id = ?", data.ChannelId).Find(&channel).Error
	if err != nil {
		c.logger.Error(
			"cannot get channel",
			slog.String("channelId", data.ChannelId),
			slog.String("channelName", data.ChannelName),
		)
		return &emptypb.Empty{}, nil
	}

	if channel.ID == "" {
		return &emptypb.Empty{}, nil
	}

	bot, ok := c.botsService.Instances[channel.BotID]
	if !ok {
		return &emptypb.Empty{}, nil
	}

	twitchClient, err := twitch.NewBotClientWithContext(ctx, bot.Model.ID, c.cfg, c.tokensGrpc)
	if err != nil {
		return nil, err
	}

	wg := twirUtils.NewGoroutinesGroup()

	for _, m := range data.MessageIds {
		wg.Go(
			func() {
				req, err := twitchClient.DeleteChatMessage(
					&helix.DeleteChatMessageParams{
						BroadcasterID: channel.ID,
						ModeratorID:   channel.BotID,
						MessageID:     m,
					},
				)
				if err != nil {
					c.logger.Error("cannot delete message", slog.Any("error", err))
				} else if req.ErrorMessage != "" {
					c.logger.Error("cannot delete message", slog.String("error", req.ErrorMessage))
				}
			},
		)
	}

	wg.Wait()

	return &emptypb.Empty{}, nil
}

func (c *BotsGrpcServer) SendMessage(
	ctx context.Context,
	data *bots.SendMessageRequest,
) (*emptypb.Empty, error) {
	if data.Message == "" {
		c.logger.Error(
			"empty message",
			slog.String("channelId", data.ChannelId),
			slog.String("channelName", data.ChannelId),
			slog.String("text", data.Message),
		)
		return &emptypb.Empty{}, errors.New("empty message")
	}

	channel := model.Channels{}
	err := c.db.WithContext(ctx).Where("id = ?", data.ChannelId).Find(&channel).Error
	if err != nil {
		c.logger.Error("cannot find channel", slog.Any("error", err))
		return nil, err
	}

	if channel.ID == "" {
		return nil, errors.New("channel is empty")
	}

	bot, ok := c.botsService.Instances[channel.BotID]
	if !ok {
		return nil, errors.New("cannot find bot associated with this channel id")
	}

	twitchClient, err := twitch.NewBotClientWithContext(ctx, bot.Model.ID, c.cfg, c.tokensGrpc)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	channelName := data.ChannelName

	if channelName == nil || *channelName == "" {
		usersReq, err := twitchClient.GetUsers(
			&helix.UsersParams{
				IDs: []string{data.ChannelId},
			},
		)
		if err != nil {
			c.logger.Error("cannot get twitch user", slog.Any("error", err))
			return nil, err
		}
		if len(usersReq.Data.Users) == 0 {
			return nil, errors.New("channel not found")
		}
		channelName = &usersReq.Data.Users[0].Login
	}

	data.Message = strings.ReplaceAll(data.Message, "\n", " ")

	if data.IsAnnounce != nil && *data.IsAnnounce {
		announceReq, err := twitchClient.SendChatAnnouncement(
			&helix.SendChatAnnouncementParams{
				BroadcasterID: channel.ID,
				ModeratorID:   channel.BotID,
				Message:       data.Message,
			},
		)
		if err != nil {
			c.logger.Error(
				"cannot send announce",
				slog.String("channelId", channel.ID),
				slog.Any("error", err),
			)
			return nil, err
		} else if announceReq.ErrorMessage != "" {
			slog.Error(
				"cannot send announce",
				slog.String("error", announceReq.Error),
				slog.String("channelId", channel.ID),
				slog.String("botId", channel.BotID),
				slog.String("message", data.Message),
				slog.String("error", announceReq.ErrorMessage),
				slog.Int("code", announceReq.StatusCode),
			)
			return nil, fmt.Errorf(
				"cannot do announce, channelId: %s, message: %s, err: %s", channel.ID, data.Message,
				announceReq.ErrorMessage,
			)
		}
	} else {
		if data.SkipRateLimits {
			bot.Say(*channelName, data.Message)
		} else {
			bot.SayWithRateLimiting(*channelName, data.Message, nil)
		}
	}

	return &emptypb.Empty{}, nil
}

func (c *BotsGrpcServer) Join(_ context.Context, data *bots.JoinOrLeaveRequest) (
	*emptypb.Empty,
	error,
) {
	bot, ok := c.botsService.Instances[data.BotId]
	if !ok {
		return nil, errors.New("bot not found")
	}

	delete(bot.RateLimiters.Channels.Items, data.UserName)
	bot.Join(data.UserName)
	return &emptypb.Empty{}, nil
}

func (c *BotsGrpcServer) Leave(_ context.Context, data *bots.JoinOrLeaveRequest) (
	*emptypb.Empty,
	error,
) {
	bot, ok := c.botsService.Instances[data.BotId]
	if !ok {
		return nil, errors.New("bot not found")
	}

	delete(bot.RateLimiters.Channels.Items, data.UserName)
	bot.Reader.Leave(data.UserName)
	bot.Writer.Leave(data.UserName)
	return &emptypb.Empty{}, nil
}
