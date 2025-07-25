package grpc

import (
	"context"
	"fmt"
	"net"
	"strconv"

	arikawa_state "github.com/diamondburned/arikawa/v3/state"
	"github.com/twirapp/twir/apps/discord/internal/discord_go"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/discord"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	arikawa_discord "github.com/diamondburned/arikawa/v3/discord"
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

	grpcNetListener, err := net.Listen(
		"tcp",
		fmt.Sprintf("0.0.0.0:%d", constants.DISCORD_SERVER_PORT),
	)
	if err != nil {
		return nil, err
	}
	grpcServer := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))

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
	gUid, err := strconv.ParseUint(req.GuildId, 10, 64)
	if err != nil {
		return nil, err
	}

	shard, _ := c.discord.FromGuildID(arikawa_discord.GuildID(gUid))
	if shard == nil {
		return nil, fmt.Errorf("shard not found")
	}
	state := shard.(*arikawa_state.State)

	channels, err := state.Channels(arikawa_discord.GuildID(gUid))
	if err != nil {
		return nil, err
	}

	discordUser, discordUserErr := state.Me()
	if discordUserErr != nil {
		return nil, discordUserErr
	}

	resultedChannels := make([]*discord.GuildChannel, 0, len(channels))
	for _, channel := range channels {
		var t discord.ChannelType

		switch channel.Type {
		case arikawa_discord.GuildText, arikawa_discord.GuildAnnouncement:
			t = discord.ChannelType_TEXT
		case arikawa_discord.GuildVoice:
			t = discord.ChannelType_VOICE
		default:
			t = -1
		}

		if t == -1 {
			continue
		}

		var hasSendMessagePerm bool

		if t == discord.ChannelType_TEXT {
			perms, permsErr := state.Permissions(channel.ID, discordUser.ID)
			if permsErr != nil {
				return nil, permsErr
			}

			hasSendMessagePerm = perms.Has(arikawa_discord.PermissionSendMessages)
		}

		resultedChannels = append(
			resultedChannels, &discord.GuildChannel{
				Id:              channel.ID.String(),
				Name:            channel.Name,
				Type:            t,
				CanSendMessages: hasSendMessagePerm,
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

	var guild *arikawa_discord.Guild
	var guildIcon string
	var guildChannels []*discord.GuildChannel
	var guildRoles []*discord.Role

	gUid, err := strconv.ParseUint(req.GuildId, 10, 64)
	if err != nil {
		return nil, err
	}

	shard, _ := c.discord.FromGuildID(arikawa_discord.GuildID(gUid))
	if shard == nil {
		return nil, fmt.Errorf("shard not found")
	}
	state := shard.(*arikawa_state.State)

	errWg.Go(
		func() error {
			disGuild, err := state.Guild(arikawa_discord.GuildID(gUid))
			if err != nil {
				return err
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
		Id:       guild.ID.String(),
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
	gUid, err := strconv.ParseUint(req.GuildId, 10, 64)
	if err != nil {
		return nil, err
	}

	shard, _ := c.discord.FromGuildID(arikawa_discord.GuildID(gUid))
	if shard == nil {
		return nil, fmt.Errorf("shard not found")
	}
	state := shard.(*arikawa_state.State)

	if err := state.LeaveGuild(arikawa_discord.GuildID(gUid)); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Impl) GetGuildRoles(ctx context.Context, req *discord.GetGuildRolesRequest) (
	*discord.
		GetGuildRolesResponse, error,
) {
	gUid, err := strconv.ParseUint(req.GuildId, 10, 64)
	if err != nil {
		return nil, err
	}

	shard, _ := c.discord.FromGuildID(arikawa_discord.GuildID(gUid))
	if shard == nil {
		return nil, fmt.Errorf("shard not found")
	}
	state := shard.(*arikawa_state.State)

	roles, err := state.Roles(arikawa_discord.GuildID(gUid))
	if err != nil {
		return nil, err
	}

	resultedRoles := make([]*discord.Role, 0, len(roles))

	for _, role := range roles {
		resultedRoles = append(
			resultedRoles, &discord.Role{
				Id:    role.ID.String(),
				Name:  role.Name,
				Color: fmt.Sprint(role.Color),
			},
		)
	}

	return &discord.GetGuildRolesResponse{
		Roles: resultedRoles,
	}, nil
}
