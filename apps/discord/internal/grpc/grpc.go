package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/satont/twir/apps/discord/internal/discord_go"
	"github.com/satont/twir/libs/grpc/generated/discord"
	"github.com/satont/twir/libs/grpc/servers"
	"github.com/satont/twir/libs/logger"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Opts struct {
	fx.In

	LC      fx.Lifecycle
	Logger  logger.Logger
	Discord *discord_go.Discord
}

func New(opts Opts) (discord.DiscordServer, error) {
	service := &Impl{
		discord: opts.Discord,
	}

	grpcNetListener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", servers.DISCORD_SERVER_PORT))
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

	discord.RegisterDiscordServer(grpcServer, service)

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go grpcServer.Serve(grpcNetListener)
				opts.Logger.Info("Grpc server is running")
				return nil
			},
			OnStop: func(ctx context.Context) error {
				grpcServer.GracefulStop()
				return nil
			},
		},
	)

	return service, nil
}

type Impl struct {
	discord.UnimplementedDiscordServer

	discord *discord_go.Discord
}

func (c *Impl) GetGuildChannels(
	ctx context.Context,
	req *discord.GetGuildChannelsRequest,
) (*discord.GetGuildChannelsResponse, error) {
	// channelsReq, err := c.discord.GuildChannels(
	// 	req.GuildId,
	// 	discordgo.WithContext(ctx),
	// )
	// if err != nil {
	// 	return nil, err
	// }
	//
	// channels := make([]*discord.GuildChannel, 0, len(channelsReq))
	// for _, channel := range channelsReq {
	// 	var t discord.ChannelType
	//
	// 	switch channel.Type {
	// 	case discordgo.ChannelTypeGuildText, discordgo.ChannelTypeGuildNews:
	// 		t = discord.ChannelType_TEXT
	// 	case discordgo.ChannelTypeGuildVoice:
	// 		t = discord.ChannelType_VOICE
	// 	default:
	// 		t = -1
	// 	}
	//
	// 	if t == -1 {
	// 		continue
	// 	}
	//
	// 	channels = append(
	// 		channels, &discord.GuildChannel{
	// 			Id:   channel.ID,
	// 			Name: channel.Name,
	// 			Type: t,
	// 		},
	// 	)
	// }

	return &discord.GetGuildChannelsResponse{
		Channels: nil,
	}, nil
}
