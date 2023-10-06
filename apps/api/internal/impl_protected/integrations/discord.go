package integrations

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/url"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/guregu/null"
	"github.com/imroc/req/v3"
	"github.com/kr/pretty"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/integrations_discord"
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
	q.Add("scope", "guilds guilds.members.read identify")
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

	if channelIntegration.ID == "" || !channelIntegration.Enabled {
		return nil, errors.New("integration not found")
	}

	return &integrations_discord.GetDataResponse{
		Id:                       *channelIntegration.Data.UserId,
		UserName:                 *channelIntegration.Data.UserName,
		Avatar:                   *channelIntegration.Data.Avatar,
		Guilds:                   channelIntegration.Data.DiscordGuilds,
		Channels:                 channelIntegration.Data.DiscordChannels,
		NotificationShowTitle:    channelIntegration.Data.DiscordNotificationShowTitle,
		NotificationShowCategory: channelIntegration.Data.DiscordNotificationShowCategory,
		NotificationMessage:      channelIntegration.Data.DiscordNotificationMessage,
	}, nil
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

	channelIntegration.Data.DiscordNotificationShowTitle = in.NotificationShowTitle
	channelIntegration.Data.DiscordNotificationShowCategory = in.NotificationShowCategory
	channelIntegration.Data.DiscordNotificationMessage = in.NotificationMessage
	channelIntegration.Data.DiscordGuilds = in.Guilds
	channelIntegration.Data.DiscordChannels = in.Channels

	if err := c.Db.WithContext(ctx).Save(&channelIntegration).Error; err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

type DiscordPostCodeResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type DiscordGetUserResponse struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
}

func (c *Integrations) IntegrationsDiscordPostCode(
	ctx context.Context,
	in *integrations_discord.PostCodeRequest,
) (*empty.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	postResult := DiscordPostCodeResponse{}

	r, err := req.
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Accept", "application/json").
		SetBasicAuth(c.Config.DiscordClientID, c.Config.DiscordClientSecret).
		SetSuccessResult(&postResult).
		SetFormData(
			map[string]string{
				"grant_type":   "authorization_code",
				"code":         in.Code,
				"redirect_uri": fmt.Sprintf("https://%s/dashboard/integrations/discord", c.Config.HostName),
			},
		).
		Post("https://discord.com/api/oauth2/token")
	if err != nil {
		c.Logger.Error("failed to get access token", slog.Any("err", err))
		return nil, err
	}
	if !r.IsSuccessState() {
		body, _ := io.ReadAll(r.Body)
		c.Logger.Error("failed to get access token", slog.String("err", string(body)))
		return nil, errors.New("failed to get access token")
	}

	channelIntegration, err := c.getChannelIntegrationByService(
		ctx, model.IntegrationServiceDiscord,
		dashboardId,
	)
	if err != nil {
		c.Logger.Error("failed to get channel integration", slog.Any("err", err))
		return nil, err
	}

	getUserResult := DiscordGetUserResponse{}

	r, err = req.
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Accept", "application/json").
		SetBearerAuthToken(postResult.AccessToken).
		SetSuccessResult(&getUserResult).
		Get("https://discord.com/api/v10/users/@me")
	if err != nil {
		c.Logger.Error("failed to get user info", slog.Any("err", err))
		return nil, err
	}
	if !r.IsSuccessState() {
		body, _ := io.ReadAll(r.Body)
		c.Logger.Error("failed to get user info", slog.String("err", string(body)))
		return nil, errors.New("failed to get user info")
	}

	pretty.Println(getUserResult)

	avatar := fmt.Sprintf(
		"https://cdn.discordapp.com/avatars/%s/%s.png",
		getUserResult.Id,
		getUserResult.Avatar,
	)

	channelIntegration.AccessToken = null.StringFrom(postResult.AccessToken)
	channelIntegration.RefreshToken = null.StringFrom(postResult.RefreshToken)
	channelIntegration.Enabled = true
	channelIntegration.Data = &model.ChannelsIntegrationsData{
		UserId:   &getUserResult.Id,
		UserName: &getUserResult.Username,
		Avatar:   &avatar,
	}

	if err := c.Db.WithContext(ctx).Save(&channelIntegration).Error; err != nil {
		c.Logger.Error("failed to save channel integration", slog.Any("err", err))
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (c *Integrations) IntegrationsDiscordLogout(
	ctx context.Context,
	_ *empty.Empty,
) (*empty.Empty, error) {
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

	if err := c.Db.WithContext(ctx).Delete(&channelIntegration).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
