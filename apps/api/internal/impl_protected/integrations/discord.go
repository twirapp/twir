package integrations

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/imroc/req/v3"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/integrations_discord"
	"github.com/satont/twir/libs/grpc/generated/discord"
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

	redirectUrl := fmt.Sprintf("https://%s/dashboard/integrations/discord", c.Config.HostName)

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
	dashboardId := ctx.Value("dashboardId").(string)

	channelIntegration, err := c.getChannelIntegrationByService(
		ctx, model.IntegrationServiceDiscord,
		dashboardId,
	)
	if err != nil {
		return nil, err
	}

	if channelIntegration.ID == "" {
		return nil, errors.New("integration not found")
	}

	if channelIntegration.Data == nil || channelIntegration.Data.UserId == nil {
		return &integrations_discord.GetDataResponse{}, nil
	}
	return &integrations_discord.GetDataResponse{}, nil
}

func (c *Integrations) IntegrationsDiscordUpdate(
	ctx context.Context,
	in *integrations_discord.UpdateMessage,
) (*empty.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	channelIntegration := model.ChannelsIntegrations{}
	if err := c.Db.WithContext(ctx).Where(
		`"channelId" = ?`,
		dashboardId,
	).Find(&channelIntegration).Error; err != nil {
		return nil, err
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
	dashboardId := ctx.Value("dashboardId").(string)
	fmt.Println("dashboardId", dashboardId)

	res := DiscordPostCodeResponse{}
	r, err := req.
		SetBasicAuth(c.Config.DiscordClientID, c.Config.DiscordClientSecret).
		SetSuccessResult(&res).
		SetFormData(
			map[string]string{
				"grant_type":   "authorization_code",
				"code":         data.Code,
				"redirect_uri": fmt.Sprintf("https://%s/dashboard/integrations/discord", c.Config.HostName),
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

	newData := &model.ChannelsIntegrationsData{
		Discord: &model.ChannelIntegrationDataDiscord{},
	}

	if channelIntegration.Data != nil && channelIntegration.Data.Discord != nil {
		newData = channelIntegration.Data
	}

	newData.Discord.Guilds = append(
		newData.Discord.Guilds,
		model.ChannelIntegrationDataDiscordGuild{
			ID: res.Guild.Id,
		},
	)

	channelIntegration.Data = newData
	if err := c.Db.WithContext(ctx).Save(&channelIntegration).Error; err != nil {
		return nil, fmt.Errorf("failed to save channel integration: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (c *Integrations) IntegrationsDiscordDisconnectGuild(
	ctx context.Context,
	req *integrations_discord.DisconnectGuildMessage,
) (*empty.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	fmt.Println("dashboardId", dashboardId)

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
				Id:   channel.Id,
				Name: channel.Name,
				Type: integrations_discord.ChannelType(channel.Type.Number()),
			},
		)
	}

	return &integrations_discord.GetGuildChannelsResponse{
		Channels: channels,
	}, nil
}
