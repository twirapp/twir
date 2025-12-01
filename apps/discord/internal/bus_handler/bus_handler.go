package bus_handler

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	arikawa_discord "github.com/diamondburned/arikawa/v3/discord"
	arikawa_state "github.com/diamondburned/arikawa/v3/state"
	"github.com/twirapp/twir/apps/discord/internal/discord_go"
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

func New(opts Opts) (*Handler, error) {
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

	return h, nil
}

func (h *Handler) subscribe() error {
	if err := h.bus.Discord.GetGuildChannels.SubscribeGroup(
		"discord",
		h.handleGetGuildChannels,
	); err != nil {
		return fmt.Errorf("failed to subscribe to GetGuildChannels: %w", err)
	}

	if err := h.bus.Discord.GetGuildInfo.SubscribeGroup(
		"discord",
		h.handleGetGuildInfo,
	); err != nil {
		return fmt.Errorf("failed to subscribe to GetGuildInfo: %w", err)
	}

	if err := h.bus.Discord.LeaveGuild.SubscribeGroup(
		"discord",
		h.handleLeaveGuild,
	); err != nil {
		return fmt.Errorf("failed to subscribe to LeaveGuild: %w", err)
	}

	if err := h.bus.Discord.GetGuildRoles.SubscribeGroup(
		"discord",
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
	gUid, err := strconv.ParseUint(req.GuildID, 10, 64)
	if err != nil {
		return discord.GetGuildChannelsResponse{}, err
	}

	shard, _ := h.discord.FromGuildID(arikawa_discord.GuildID(gUid))
	if shard == nil {
		return discord.GetGuildChannelsResponse{}, fmt.Errorf("shard not found")
	}
	state := shard.(*arikawa_state.State)

	channels, err := state.Channels(arikawa_discord.GuildID(gUid))
	if err != nil {
		return discord.GetGuildChannelsResponse{}, err
	}

	discordUser, discordUserErr := state.Me()
	if discordUserErr != nil {
		return discord.GetGuildChannelsResponse{}, discordUserErr
	}

	resultedChannels := make([]discord.GuildChannel, 0, len(channels))
	for _, channel := range channels {
		var t discord.ChannelType

		switch channel.Type {
		case arikawa_discord.GuildText, arikawa_discord.GuildAnnouncement:
			t = discord.ChannelTypeText
		case arikawa_discord.GuildVoice:
			t = discord.ChannelTypeVoice
		default:
			t = -1
		}

		if t == -1 {
			continue
		}

		var hasSendMessagePerm bool

		if t == discord.ChannelTypeText {
			perms, permsErr := state.Permissions(channel.ID, discordUser.ID)
			if permsErr != nil {
				return discord.GetGuildChannelsResponse{}, permsErr
			}

			hasSendMessagePerm = perms.Has(arikawa_discord.PermissionSendMessages)
		}

		resultedChannels = append(
			resultedChannels, discord.GuildChannel{
				ID:              channel.ID.String(),
				Name:            channel.Name,
				Type:            t,
				CanSendMessages: hasSendMessagePerm,
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

	var guild *arikawa_discord.Guild
	var guildChannels []discord.GuildChannel
	var guildRoles []discord.Role

	gUid, err := strconv.ParseUint(req.GuildID, 10, 64)
	if err != nil {
		return discord.GetGuildInfoResponse{}, err
	}

	shard, _ := h.discord.FromGuildID(arikawa_discord.GuildID(gUid))
	if shard == nil {
		return discord.GetGuildInfoResponse{}, fmt.Errorf("shard not found")
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
		ID:       guild.ID.String(),
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
	gUid, err := strconv.ParseUint(req.GuildID, 10, 64)
	if err != nil {
		return struct{}{}, err
	}

	shard, _ := h.discord.FromGuildID(arikawa_discord.GuildID(gUid))
	if shard == nil {
		return struct{}{}, fmt.Errorf("shard not found")
	}
	state := shard.(*arikawa_state.State)

	if err := state.LeaveGuild(arikawa_discord.GuildID(gUid)); err != nil {
		return struct{}{}, err
	}

	return struct{}{}, nil
}

func (h *Handler) handleGetGuildRoles(
	ctx context.Context,
	req discord.GetGuildRolesRequest,
) (discord.GetGuildRolesResponse, error) {
	gUid, err := strconv.ParseUint(req.GuildID, 10, 64)
	if err != nil {
		return discord.GetGuildRolesResponse{}, err
	}

	shard, _ := h.discord.FromGuildID(arikawa_discord.GuildID(gUid))
	if shard == nil {
		return discord.GetGuildRolesResponse{}, fmt.Errorf("shard not found")
	}
	state := shard.(*arikawa_state.State)

	roles, err := state.Roles(arikawa_discord.GuildID(gUid))
	if err != nil {
		return discord.GetGuildRolesResponse{}, err
	}

	resultedRoles := make([]discord.Role, 0, len(roles))

	for _, role := range roles {
		resultedRoles = append(
			resultedRoles, discord.Role{
				ID:    role.ID.String(),
				Name:  role.Name,
				Color: fmt.Sprint(role.Color),
			},
		)
	}

	return discord.GetGuildRolesResponse{
		Roles: resultedRoles,
	}, nil
}
