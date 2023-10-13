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
	"github.com/switchupcb/disgo"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/types/known/emptypb"
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
	channelsReq := disgo.GetGuildChannels{GuildID: req.GuildId}
	channels, err := channelsReq.Send(c.discord.Client)
	if err != nil {
		return nil, err
	}

	resultedChannels := make([]*discord.GuildChannel, 0, len(channels))
	for _, channel := range channels {
		var t discord.ChannelType

		if channel.Type == nil {
			continue
		}

		switch *channel.Type {
		case disgo.FlagChannelTypeGUILD_TEXT, disgo.FlagChannelTypeGUILD_ANNOUNCEMENT:
			t = discord.ChannelType_TEXT
		case disgo.FlagChannelTypeGUILD_VOICE:
			t = discord.ChannelType_VOICE
		default:
			t = -1
		}

		if t == -1 {
			continue
		}

		var name string
		if channel.Name != nil {
			name = **channel.Name
		}

		resultedChannels = append(
			resultedChannels, &discord.GuildChannel{
				Id:   channel.ID,
				Name: name,
				Type: t,
			},
		)
	}

	return &discord.GetGuildChannelsResponse{
		Channels: resultedChannels,
	}, nil
}

func (c *Impl) GetGuildInfo(
	ctx context.Context,
	req *discord.GetGuildInfoRequest,
) (*discord.GetGuildInfoResponse, error) {
	errWg, _ := errgroup.WithContext(ctx)

	var guild *disgo.Guild
	var guildIcon string
	var guildChannels []*discord.GuildChannel
	var guildRoles []*discord.Role

	errWg.Go(
		func() error {
			guildReq := disgo.GetGuild{GuildID: req.GuildId}
			var err error
			disGuild, err := guildReq.Send(c.discord.Client)
			if err != nil {
				return err
			}

			if disGuild.Icon != nil {
				guildIcon = *disGuild.Icon
			}

			guild = disGuild
			return nil
		},
	)

	errWg.Go(
		func() error {
			channels, err := c.GetGuildChannels(
				ctx,
				&discord.GetGuildChannelsRequest{GuildId: req.GuildId},
			)
			if err != nil {
				return err
			}

			guildChannels = channels.Channels

			return nil
		},
	)

	errWg.Go(
		func() error {
			roles, err := c.GetGuildRoles(
				ctx,
				&discord.GetGuildRolesRequest{GuildId: req.GuildId},
			)
			if err != nil {
				return err
			}

			guildRoles = roles.Roles

			return nil
		},
	)

	if err := errWg.Wait(); err != nil {
		return nil, fmt.Errorf("failed to get guild info: %w", err)
	}

	return &discord.GetGuildInfoResponse{
		Id:       guild.ID,
		Name:     guild.Name,
		Icon:     guildIcon,
		Channels: guildChannels,
		Roles:    guildRoles,
	}, nil
}

func (c *Impl) LeaveGuild(
	ctx context.Context,
	req *discord.LeaveGuildRequest,
) (*emptypb.Empty, error) {
	leaveGuildReq := disgo.LeaveGuild{GuildID: req.GuildId}
	if err := leaveGuildReq.Send(c.discord.Client); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Impl) GetGuildRoles(ctx context.Context, req *discord.GetGuildRolesRequest) (
	*discord.
		GetGuildRolesResponse, error,
) {
	rolesReq := disgo.GetGuildRoles{GuildID: req.GuildId}
	roles, err := rolesReq.Send(c.discord.Client)
	if err != nil {
		return nil, err
	}

	resultedRoles := make([]*discord.Role, 0, len(roles))

	for _, role := range roles {
		resultedRoles = append(
			resultedRoles, &discord.Role{
				Id:    role.ID,
				Name:  role.Name,
				Color: fmt.Sprint(role.Color),
			},
		)
	}

	return &discord.GetGuildRolesResponse{
		Roles: resultedRoles,
	}, nil
}
