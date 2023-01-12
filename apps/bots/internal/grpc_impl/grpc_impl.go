package grpc_impl

import (
	"context"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/bots/internal/di"
	cfg "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/twitch"

	"github.com/satont/go-helix/v2"
	internalBots "github.com/satont/tsuwari/apps/bots/internal/bots"
	"github.com/satont/tsuwari/apps/bots/pkg/utils"
	"github.com/satont/tsuwari/apps/bots/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/bots"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type GrpcImplOpts struct {
	Db          *gorm.DB
	BotsService *internalBots.BotsService
	Logger      *zap.Logger
	Cfg         *cfg.Config
}

type botsGrpcServer struct {
	bots.UnimplementedBotsServer

	db          *gorm.DB
	botsService *internalBots.BotsService
	logger      *zap.Logger
	cfg         *cfg.Config
}

func NewServer(opts *GrpcImplOpts) *botsGrpcServer {
	return &botsGrpcServer{
		db:          opts.Db,
		botsService: opts.BotsService,
		logger:      opts.Logger,
		cfg:         opts.Cfg,
	}
}

func (c *botsGrpcServer) DeleteMessage(ctx context.Context, data *bots.DeleteMessagesRequest) (*emptypb.Empty, error) {
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
		go twitchClient.DeleteMessage(&helix.DeleteMessageParams{
			BroadcasterID: channel.ID,
			ModeratorID:   channel.BotID,
			MessageID:     m,
		})
	}
	return &emptypb.Empty{}, nil
}

func (c *botsGrpcServer) SendMessage(ctx context.Context, data *bots.SendMessageRequest) (*emptypb.Empty, error) {
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
		return &emptypb.Empty{}, err
	}

	channelName := data.ChannelName

	if channelName == nil {
		usersReq, err := twitchClient.GetUsers(&helix.UsersParams{
			IDs: []string{data.ChannelId},
		})
		if err != nil {
			c.logger.Sugar().Error(err)
			return &emptypb.Empty{}, nil
		}
		if len(usersReq.Data.Users) == 0 {
			return &emptypb.Empty{}, nil
		}
		channelName = &usersReq.Data.Users[0].Login
	}

	if data.IsAnnounce != nil && *data.IsAnnounce == true {
		twitchClient.SendChatAnnouncement(&helix.SendChatAnnouncementParams{
			BroadcasterID: channel.ID,
			ModeratorID:   channel.BotID,
			Message:       data.Message,
		})
	} else {
		bot.SayWithRateLimiting(*channelName, data.Message, nil)
	}

	return &emptypb.Empty{}, nil
}

func (c *botsGrpcServer) Join(ctx context.Context, data *bots.JoinOrLeaveRequest) (*emptypb.Empty, error) {
	bot, ok := c.botsService.Instances[data.BotId]
	if !ok {
		return &emptypb.Empty{}, nil
	}

	rateLimitedChannel := bot.RateLimiters.Channels.Items[data.UserName]
	if rateLimitedChannel == nil {
		bot.RateLimiters.Channels.Lock()
		defer bot.RateLimiters.Channels.Unlock()
		limiter := utils.CreateBotLimiter(false)
		bot.RateLimiters.Channels.Items[data.UserName] = &types.Channel{
			Limiter: limiter,
		}
	}
	bot.Join(data.UserName)
	return &emptypb.Empty{}, nil
}

func (c *botsGrpcServer) Leave(ctx context.Context, data *bots.JoinOrLeaveRequest) (*emptypb.Empty, error) {
	bot, ok := c.botsService.Instances[data.BotId]
	if !ok {
		return nil, nil
	}

	delete(bot.RateLimiters.Channels.Items, data.UserName)
	bot.Depart(data.UserName)
	return nil, nil
}
