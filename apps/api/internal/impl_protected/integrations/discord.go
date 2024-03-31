package integrations

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strings"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/helpers"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/integrations_discord"
	"github.com/twirapp/twir/libs/grpc/discord"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Integrations) IntegrationsDiscordGetAuthLink(
	_ context.Context,
	_ *empty.Empty,
) (*integrations_discord.GetAuthLink, error) {
	u, _ := url.Parse("https://discord.com/oauth2/authorize")

	if c.Config.DiscordClientID == "" || c.Config.DiscordClientSecret == "" {
		return nil, errors.New("discord not enabled on our side, please be patient")
	}

	redirectUrl := fmt.Sprintf("http://%s/dashboard/integrations/discord", c.Config.SiteBaseUrl)

	q := u.Query()
	q.Add("client_id", c.Config.DiscordClientID)
	q.Add("response_type", "code")
	q.Add("permissions", "1497333180438")
	q.Add("scope", "bot applications.commands")
	q.Add("redirect_uri", redirectUrl)
	u.RawQuery = q.Encode()

	str := u.String()

	return &integrations_discord.GetAuthLink{Link: str}, nil
}

func (c *Integrations) IntegrationsDiscordGetData(
	ctx context.Context,
	_ *empty.Empty,
) (*integrations_discord.GetDataResponse, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	channelIntegration, err := c.getChannelIntegrationByService(
		ctx,
		model.IntegrationServiceDiscord,
		dashboardId,
	)
	if err != nil {
		return nil, err
	}

	if channelIntegration.ID == "" {
		return nil, errors.New("integration not found")
	}

	if channelIntegration.Data == nil || channelIntegration.Data.Discord == nil {
		return &integrations_discord.GetDataResponse{}, nil
	}

	guilds := make(
		[]*integrations_discord.DiscordGuild,
		0,
		len(channelIntegration.Data.Discord.Guilds),
	)
	guildsMu := sync.Mutex{}

	guildsGroup, gCtx := errgroup.WithContext(ctx)
	for _, guild := range channelIntegration.Data.Discord.Guilds {
		guild := guild

		guildsGroup.Go(
			func() error {
				g, err := c.Grpc.Discord.GetGuildInfo(
					gCtx, &discord.GetGuildInfoRequest{
						GuildId: guild.ID,
					},
				)
				if err != nil {
					return err
				}

				channels := make([]*integrations_discord.GuildChannel, 0, len(g.GetChannels()))
				for _, channel := range g.GetChannels() {
					channels = append(
						channels,
						&integrations_discord.GuildChannel{
							Id:              channel.GetId(),
							Name:            channel.GetName(),
							Type:            integrations_discord.ChannelType(channel.GetType().Number()),
							CanSendMessages: channel.GetCanSendMessages(),
						},
					)
				}

				roles := make([]*integrations_discord.GuildRole, 0, len(g.GetRoles()))
				for _, role := range g.Roles {
					roles = append(
						roles,
						&integrations_discord.GuildRole{
							Id:    role.GetId(),
							Name:  role.GetName(),
							Color: role.GetColor(),
						},
					)
				}

				guildsMu.Lock()
				guilds = append(
					guilds,
					&integrations_discord.DiscordGuild{
						Id:                                       g.GetId(),
						Name:                                     g.GetName(),
						Icon:                                     g.GetIcon(),
						LiveNotificationEnabled:                  guild.LiveNotificationEnabled,
						LiveNotificationChannelsIds:              guild.LiveNotificationChannelsIds,
						LiveNotificationShowTitle:                guild.LiveNotificationShowTitle,
						LiveNotificationShowCategory:             guild.LiveNotificationShowCategory,
						LiveNotificationShowViewers:              guild.LiveNotificationShowViewers,
						LiveNotificationMessage:                  guild.LiveNotificationMessage,
						LiveNotificationAdditionalTwitchUsersIds: guild.LiveNotificationChannelsIds,
						LiveNotificationShowPreview:              guild.LiveNotificationShowPreview,
						LiveNotificationShowProfileImage:         guild.LiveNotificationShowProfileImage,
						Channels:                                 channels,
						Roles:                                    roles,
						OfflineNotificationMessage:               guild.OfflineNotificationMessage,
						ShouldDeleteMessageOnOffline:             guild.ShouldDeleteMessageOnOffline,
						AdditionalUsersIdsForLiveCheck:           guild.AdditionalUsersIdsForLiveCheck,
					},
				)
				guildsMu.Unlock()

				return nil
			},
		)
	}

	if err := guildsGroup.Wait(); err != nil {
		return nil, fmt.Errorf("failed to get guilds: %w", err)
	}

	slices.SortFunc(
		guilds,
		func(a, b *integrations_discord.DiscordGuild) int {
			return cmp.Compare(strings.ToLower(a.GetName()), strings.ToLower(b.GetName()))
		},
	)

	return &integrations_discord.GetDataResponse{
		Guilds: guilds,
	}, nil
}

func (c *Integrations) IntegrationsDiscordUpdate(
	ctx context.Context,
	req *integrations_discord.UpdateMessage,
) (*empty.Empty, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	channelIntegration, err := c.getChannelIntegrationByService(
		ctx,
		model.IntegrationServiceDiscord,
		dashboardId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel integration: %w", err)
	}

	newGuilds := make([]model.ChannelIntegrationDataDiscordGuild, 0, len(req.GetGuilds()))
	for _, guild := range req.GetGuilds() {
		newGuilds = append(
			newGuilds,
			model.ChannelIntegrationDataDiscordGuild{
				ID:                               guild.GetId(),
				LiveNotificationEnabled:          guild.GetLiveNotificationEnabled(),
				LiveNotificationChannelsIds:      guild.GetLiveNotificationChannelsIds(),
				LiveNotificationShowTitle:        guild.GetLiveNotificationShowTitle(),
				LiveNotificationShowCategory:     guild.GetLiveNotificationShowCategory(),
				LiveNotificationShowViewers:      guild.GetLiveNotificationShowViewers(),
				LiveNotificationMessage:          guild.GetLiveNotificationMessage(),
				OfflineNotificationMessage:       guild.GetOfflineNotificationMessage(),
				LiveNotificationShowPreview:      guild.GetLiveNotificationShowPreview(),
				LiveNotificationShowProfileImage: guild.GetLiveNotificationShowProfileImage(),
				ShouldDeleteMessageOnOffline:     guild.GetShouldDeleteMessageOnOffline(),
				AdditionalUsersIdsForLiveCheck:   guild.GetAdditionalUsersIdsForLiveCheck(),
			},
		)
	}

	channelIntegration.Data.Discord = &model.ChannelIntegrationDataDiscord{
		Guilds: newGuilds,
	}

	if err := c.Db.WithContext(ctx).Save(&channelIntegration).Error; err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

type DiscordPostCodeResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Guild        struct {
		Id string `json:"id"`
	} `json:"guild"`
}

func (c *Integrations) IntegrationDiscordConnectGuild(
	ctx context.Context,
	data *integrations_discord.PostCodeRequest,
) (*empty.Empty, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	res := DiscordPostCodeResponse{}
	r, err := req.
		SetBasicAuth(c.Config.DiscordClientID, c.Config.DiscordClientSecret).
		SetSuccessResult(&res).
		SetFormData(
			map[string]string{
				"grant_type": "authorization_code",
				"code":       data.GetCode(),
				"redirect_uri": fmt.Sprintf(
					"http://%s/dashboard/integrations/discord",
					c.Config.SiteBaseUrl,
				),
			},
		).
		Post("https://discord.com/api/oauth2/token")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	if !r.IsSuccessState() {
		return nil, fmt.Errorf("failed to get token: %s", r.String())
	}
	if res.Guild.Id == "" {
		return nil, fmt.Errorf("failed to get guild id")
	}

	channelIntegration, err := c.getChannelIntegrationByService(
		ctx,
		model.IntegrationServiceDiscord,
		dashboardId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel integration: %w", err)
	}

	if channelIntegration.Data.Discord == nil {
		channelIntegration.Data.Discord = &model.ChannelIntegrationDataDiscord{
			Guilds: []model.ChannelIntegrationDataDiscordGuild{},
		}
	}

	if !lo.SomeBy(
		channelIntegration.Data.Discord.Guilds,
		func(item model.ChannelIntegrationDataDiscordGuild) bool {
			return item.ID == res.Guild.Id
		},
	) || channelIntegration.Data.Discord == nil {
		channelIntegration.Data.Discord.Guilds = append(
			channelIntegration.Data.Discord.Guilds,
			model.ChannelIntegrationDataDiscordGuild{
				ID:                               res.Guild.Id,
				LiveNotificationShowTitle:        true,
				LiveNotificationShowCategory:     true,
				LiveNotificationShowViewers:      true,
				LiveNotificationShowProfileImage: true,
				LiveNotificationShowPreview:      true,
			},
		)

		if err := c.Db.WithContext(ctx).Save(&channelIntegration).Error; err != nil {
			return nil, fmt.Errorf("failed to save channel integration: %w", err)
		}
	}

	return &emptypb.Empty{}, nil
}

func (c *Integrations) IntegrationsDiscordDisconnectGuild(
	ctx context.Context,
	req *integrations_discord.DisconnectGuildMessage,
) (*empty.Empty, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	channelIntegration, err := c.getChannelIntegrationByService(
		ctx,
		model.IntegrationServiceDiscord,
		dashboardId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel integration: %w", err)
	}

	if channelIntegration.Data == nil || channelIntegration.Data.Discord == nil {
		return nil, fmt.Errorf("failed to get channel integration data")
	}

	channelIntegration.Data.Discord.Guilds = lo.Filter(
		channelIntegration.Data.Discord.Guilds,
		func(guild model.ChannelIntegrationDataDiscordGuild, _ int) bool {
			return guild.ID != req.GuildId
		},
	)

	if err := c.Db.WithContext(ctx).Save(&channelIntegration).Error; err != nil {
		return nil, fmt.Errorf("failed to save channel integration: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (c *Integrations) IntegrationsDiscordGetGuildChannels(
	ctx context.Context,
	req *integrations_discord.GetGuildChannelsRequest,
) (
	*integrations_discord.GetGuildChannelsResponse,
	error,
) {
	channelsReq, err := c.Grpc.Discord.GetGuildChannels(
		ctx, &discord.GetGuildChannelsRequest{
			GuildId: req.GuildId,
		},
	)
	if err != nil {
		return nil, err
	}

	channels := make([]*integrations_discord.GuildChannel, 0, len(channelsReq.Channels))
	for _, channel := range channelsReq.Channels {
		channels = append(
			channels,
			&integrations_discord.GuildChannel{
				Id:              channel.Id,
				Name:            channel.Name,
				Type:            integrations_discord.ChannelType(channel.Type.Number()),
				CanSendMessages: channel.CanSendMessages,
			},
		)
	}

	return &integrations_discord.GetGuildChannelsResponse{
		Channels: channels,
	}, nil
}

func (c *Integrations) IntegrationsDiscordGetGuildInfo(
	ctx context.Context,
	req *integrations_discord.GetGuildInfoRequest,
) (*integrations_discord.GetGuildInfoResponse, error) {
	guildInfo, err := c.Grpc.Discord.GetGuildInfo(
		ctx, &discord.GetGuildInfoRequest{
			GuildId: req.GuildId,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get guild info: %w", err)
	}

	channels := make([]*integrations_discord.GuildChannel, 0, len(guildInfo.Channels))
	for _, channel := range guildInfo.Channels {
		channels = append(
			channels,
			&integrations_discord.GuildChannel{
				Id:              channel.Id,
				Name:            channel.Name,
				Type:            integrations_discord.ChannelType(channel.Type.Number()),
				CanSendMessages: channel.CanSendMessages,
			},
		)
	}

	roles := make([]*integrations_discord.GuildRole, 0, len(guildInfo.Roles))
	for _, role := range guildInfo.Roles {
		roles = append(
			roles,
			&integrations_discord.GuildRole{
				Id:    role.Id,
				Name:  role.Name,
				Color: role.Color,
			},
		)
	}

	return &integrations_discord.GetGuildInfoResponse{
		Id:       guildInfo.Id,
		Name:     guildInfo.Name,
		Icon:     guildInfo.Icon,
		Channels: channels,
		Roles:    roles,
	}, nil
}
