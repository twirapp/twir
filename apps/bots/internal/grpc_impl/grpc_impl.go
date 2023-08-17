package grpc_impl

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/samber/do"
	"github.com/satont/twir/apps/bots/internal/di"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/satont/twir/libs/twitch"
	"log/slog"

	"github.com/nicklaw5/helix/v2"
	internalBots "github.com/satont/twir/apps/bots/internal/bots"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type GrpcImplOpts struct {
	Db          *gorm.DB
	BotsService *internalBots.Service
	Logger      *zap.Logger
	Cfg         *cfg.Config
}

type BotsGrpcServer struct {
	bots.UnimplementedBotsServer

	db          *gorm.DB
	botsService *internalBots.Service
	logger      *zap.Logger
	cfg         *cfg.Config
}

func NewServer(opts *GrpcImplOpts) *BotsGrpcServer {
	return &BotsGrpcServer{
		db:          opts.Db,
		botsService: opts.BotsService,
		logger:      opts.Logger,
		cfg:         opts.Cfg,
	}
}

func (c *BotsGrpcServer) DeleteMessage(ctx context.Context, data *bots.DeleteMessagesRequest) (*emptypb.Empty, error) {
	channel := model.Channels{}
	err := c.db.Where("id = ?", data.ChannelId).Find(&channel).Error
	if err != nil {
		c.logger.Sugar().Error(err)
		return &emptypb.Empty{}, nil
	}

	if channel.ID == "" {
		return &emptypb.Empty{}, nil
	}

	bot, ok := c.botsService.Instances[channel.BotID]
	if !ok {
		return &emptypb.Empty{}, nil
	}

	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)
	twitchClient, err := twitch.NewBotClient(bot.Model.ID, *c.cfg, tokensGrpc)
	if err != nil {
		return nil, err
	}

	for _, m := range data.MessageIds {
		go twitchClient.DeleteChatMessage(
			&helix.DeleteChatMessageParams{
				BroadcasterID: channel.ID,
				ModeratorID:   channel.BotID,
				MessageID:     m,
			},
		)
	}
	return &emptypb.Empty{}, nil
}

func (c *BotsGrpcServer) SendMessage(ctx context.Context, data *bots.SendMessageRequest) (*emptypb.Empty, error) {
	if data.Message == "" {
		c.logger.Sugar().Error(
			"empty message",
			zap.String("channelId", data.ChannelId),
			zap.String("channelName", data.ChannelId),
			zap.String("text", data.Message),
		)
		return &emptypb.Empty{}, errors.New("empty message")
	}

	channel := model.Channels{}
	err := c.db.WithContext(ctx).Where("id = ?", data.ChannelId).Find(&channel).Error
	if err != nil {
		c.logger.Sugar().Error(err)
		return nil, err
	}

	if channel.ID == "" {
		return nil, errors.New("channel is empty")
	}

	bot, ok := c.botsService.Instances[channel.BotID]
	if !ok {
		return nil, errors.New("cannot find bot associated with this channel id")
	}

	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)
	twitchClient, err := twitch.NewBotClientWithContext(ctx, bot.Model.ID, *c.cfg, tokensGrpc)
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
			c.logger.Sugar().Error(err)
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
			c.logger.Sugar().Error(err, zap.String("channelId", channel.ID))
			return nil, err
		} else if announceReq.ErrorMessage != "" {
			slog.Error(
				"cannot do announce "+announceReq.ErrorMessage,
				slog.String("error", announceReq.Error),
				slog.String("channelId", channel.ID),
				slog.String("botId", channel.BotID),
				slog.String("message", data.Message),
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

func (c *BotsGrpcServer) Join(ctx context.Context, data *bots.JoinOrLeaveRequest) (*emptypb.Empty, error) {
	bot, ok := c.botsService.Instances[data.BotId]
	if !ok {
		return nil, errors.New("bot not found")
	}

	delete(bot.RateLimiters.Channels.Items, data.UserName)
	bot.Join(data.UserName)
	return &emptypb.Empty{}, nil
}

func (c *BotsGrpcServer) Leave(ctx context.Context, data *bots.JoinOrLeaveRequest) (*emptypb.Empty, error) {
	bot, ok := c.botsService.Instances[data.BotId]
	if !ok {
		return nil, errors.New("bot not found")
	}

	delete(bot.RateLimiters.Channels.Items, data.UserName)
	bot.Depart(data.UserName)
	return &emptypb.Empty{}, nil
}
