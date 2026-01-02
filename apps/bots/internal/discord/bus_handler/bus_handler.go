package bus_handler

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/twirapp/twir/apps/bots/internal/discord/discord_go"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/discord"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
)

type Opts struct {
	fx.In

	LC      fx.Lifecycle
	Logger  *slog.Logger
	Discord *discord_go.Discord
	Bus     *buscore.Bus
}

type Handler struct {
	discord *discord_go.Discord
	bus     *buscore.Bus
	logger  *slog.Logger
}

func New(opts Opts) error {
	h := &Handler{
		discord: opts.Discord,
		bus:     opts.Bus,
		logger:  opts.Logger,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := h.subscribe(); err != nil {
					return err
				}
				opts.Logger.Info("Discord bus handler is running")
				return nil
			},
			OnStop: func(ctx context.Context) error {
				h.unsubscribe()
				return nil
			},
		},
	)

	return nil
}

func (h *Handler) subscribe() error {
	if err := h.bus.Discord.GetGuildChannels.SubscribeGroup(
		"bots",
		h.handleGetGuildChannels,
	); err != nil {
		return fmt.Errorf("failed to subscribe to GetGuildChannels: %w", err)
	}

	if err := h.bus.Discord.GetGuildInfo.SubscribeGroup(
		"bots",
		h.handleGetGuildInfo,
	); err != nil {
		return fmt.Errorf("failed to subscribe to GetGuildInfo: %w", err)
	}

	if err := h.bus.Discord.LeaveGuild.SubscribeGroup(
		"bots",
		h.handleLeaveGuild,
	); err != nil {
		return fmt.Errorf("failed to subscribe to LeaveGuild: %w", err)
	}

	if err := h.bus.Discord.GetGuildRoles.SubscribeGroup(
		"bots",
		h.handleGetGuildRoles,
	); err != nil {
		return fmt.Errorf("failed to subscribe to GetGuildRoles: %w", err)
	}

	return nil
}

func (h *Handler) unsubscribe() {
	h.bus.Discord.GetGuildChannels.Unsubscribe()
	h.bus.Discord.GetGuildInfo.Unsubscribe()
	h.bus.Discord.LeaveGuild.Unsubscribe()
	h.bus.Discord.GetGuildRoles.Unsubscribe()
}

func (h *Handler) handleGetGuildChannels(
	ctx context.Context,
	req discord.GetGuildChannelsRequest,
) (discord.GetGuildChannelsResponse, error) {
	channels, err := h.discord.GetGuildChannels(ctx, req.GuildID)
	if err != nil {
		return discord.GetGuildChannelsResponse{}, err
	}

	resultedChannels := make([]discord.GuildChannel, 0, len(channels))
	for _, channel := range channels {
		var t discord.ChannelType

		switch channel.Type {
		case discord_go.DiscordChannelTypeText, discord_go.DiscordChannelTypeAnnouncement:
			t = discord.ChannelTypeText
		case discord_go.DiscordChannelTypeVoice:
			t = discord.ChannelTypeVoice
		default:
			t = -1
		}

		if t == -1 {
			continue
		}

		resultedChannels = append(
			resultedChannels, discord.GuildChannel{
				ID:              channel.ID,
				Name:            channel.Name,
				Type:            t,
				CanSendMessages: channel.CanSendMessages,
			},
		)
	}

	return discord.GetGuildChannelsResponse{
		Channels: resultedChannels,
	}, nil
}

func (h *Handler) handleGetGuildInfo(
	ctx context.Context,
	req discord.GetGuildInfoRequest,
) (discord.GetGuildInfoResponse, error) {
	errWg, _ := errgroup.WithContext(ctx)

	var guild discord.GetGuildInfoResponse
	var guildChannels []discord.GuildChannel
	var guildRoles []discord.Role

	errWg.Go(
		func() error {
			disGuild, err := h.discord.GetGuild(ctx, req.GuildID)
			if err != nil {
				return err
			}

			guild = discord.GetGuildInfoResponse{
				ID:   disGuild.ID,
				Name: disGuild.Name,
				Icon: disGuild.Icon,
			}
			return nil
		},
	)

	errWg.Go(
		func() error {
			channels, err := h.handleGetGuildChannels(
				ctx,
				discord.GetGuildChannelsRequest{GuildID: req.GuildID},
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
			roles, err := h.handleGetGuildRoles(
				ctx,
				discord.GetGuildRolesRequest{GuildID: req.GuildID},
			)
			if err != nil {
				return err
			}

			guildRoles = roles.Roles

			return nil
		},
	)

	if err := errWg.Wait(); err != nil {
		return discord.GetGuildInfoResponse{}, fmt.Errorf("failed to get guild info: %w", err)
	}

	return discord.GetGuildInfoResponse{
		ID:       guild.ID,
		Name:     guild.Name,
		Icon:     guild.Icon,
		Channels: guildChannels,
		Roles:    guildRoles,
	}, nil
}

func (h *Handler) handleLeaveGuild(
	ctx context.Context,
	req discord.LeaveGuildRequest,
) (struct{}, error) {
	if err := h.discord.LeaveGuild(ctx, req.GuildID); err != nil {
		return struct{}{}, err
	}

	return struct{}{}, nil
}

func (h *Handler) handleGetGuildRoles(
	ctx context.Context,
	req discord.GetGuildRolesRequest,
) (discord.GetGuildRolesResponse, error) {
	roles, err := h.discord.GetGuildRoles(ctx, req.GuildID)
	if err != nil {
		return discord.GetGuildRolesResponse{}, err
	}

	resultedRoles := make([]discord.Role, 0, len(roles))

	for _, role := range roles {
		resultedRoles = append(
			resultedRoles, discord.Role{
				ID:    role.ID,
				Name:  role.Name,
				Color: fmt.Sprint(role.Color),
			},
		)
	}

	return discord.GetGuildRolesResponse{
		Roles: resultedRoles,
	}, nil
}
