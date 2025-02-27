package integrations

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/guregu/null"
	"github.com/imroc/req/v3"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/integrations_streamlabs"
	"google.golang.org/protobuf/types/known/emptypb"
)

type streamlabsTokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type streamlabsProfileResponse struct {
	StreamLabs struct {
		DisplayName string `json:"display_name"`
		ThumbNail   string `json:"thumbnail"`
		ID          int    `json:"id"`
	} `json:"streamlabs"`
}

func (c *Integrations) IntegrationsStreamlabsGetAuthLink(
	ctx context.Context,
	_ *emptypb.Empty,
) (*integrations_streamlabs.GetAuthLink, error) {
	integration, err := c.getIntegrationByService(ctx, model.IntegrationServiceStreamLabs)
	if err != nil {
		return nil, err
	}

	if !integration.ClientID.Valid || !integration.ClientSecret.Valid || !integration.RedirectURL.Valid {
		return nil, errors.New("streamlabs not enabled on our side, please be patient")
	}

	link, _ := url.Parse("https://www.streamlabs.com/api/v2.0/authorize")

	q := link.Query()
	q.Add("response_type", "code")
	q.Add("client_id", integration.ClientID.String)
	q.Add("scope", "socket.token donations.read")
	q.Add("redirect_uri", integration.RedirectURL.String)
	link.RawQuery = q.Encode()

	return &integrations_streamlabs.GetAuthLink{
		Link: link.String(),
	}, nil
}

func (c *Integrations) IntegrationsStreamlabsGetData(
	ctx context.Context,
	_ *emptypb.Empty,
) (*integrations_streamlabs.GetDataResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	integration, err := c.getChannelIntegrationByService(
		ctx,
		model.IntegrationServiceStreamLabs,
		dashboardId,
	)
	if err != nil {
		return nil, err
	}

	return &integrations_streamlabs.GetDataResponse{
		UserName: integration.Data.UserName,
		Avatar:   integration.Data.Avatar,
	}, nil
}

func (c *Integrations) IntegrationsStreamlabsPostCode(
	ctx context.Context,
	request *integrations_streamlabs.PostCodeRequest,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	channelIntegration, err := c.getChannelIntegrationByService(
		ctx,
		model.IntegrationServiceStreamLabs,
		dashboardId,
	)
	if err != nil {
		return nil, err
	}

	tokensData := streamlabsTokensResponse{}
	resp, err := req.R().
		SetContext(ctx).
		SetFormData(
			map[string]string{
				"grant_type":    "authorization_code",
				"client_id":     channelIntegration.Integration.ClientID.String,
				"client_secret": channelIntegration.Integration.ClientSecret.String,
				"redirect_uri":  channelIntegration.Integration.RedirectURL.String,
				"code":          request.Code,
			},
		).
		SetSuccessResult(&tokensData).
		SetContentType("application/x-www-form-urlencoded").
		Post("https://streamlabs.com/api/v2.0/token")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("streamlabs token request failed: %s", resp.String())
	}

	profileData := &streamlabsProfileResponse{}
	resp, err = req.R().
		SetContext(ctx).
		SetSuccessResult(profileData).
		SetBearerAuthToken(tokensData.AccessToken).
		Get("https://streamlabs.com/api/v2.0/user")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("streamlabs token request failed: %s", resp.String())
	}

	channelIntegration.Data = &model.ChannelsIntegrationsData{
		UserName: &profileData.StreamLabs.DisplayName,
		Avatar:   &profileData.StreamLabs.ThumbNail,
	}
	channelIntegration.AccessToken = null.StringFrom(tokensData.AccessToken)
	channelIntegration.RefreshToken = null.StringFrom(tokensData.RefreshToken)
	channelIntegration.Enabled = true

	if err = c.Db.WithContext(ctx).Save(channelIntegration).Error; err != nil {
		return nil, err
	}

	if err = c.sendGrpcEvent(ctx, channelIntegration.ID, channelIntegration.Enabled); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Integrations) IntegrationsStreamlabsLogout(
	ctx context.Context,
	empty *emptypb.Empty,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	integration, err := c.getChannelIntegrationByService(
		ctx, model.IntegrationServiceStreamLabs, dashboardId,
	)
	if err != nil {
		return nil, err
	}

	integration.Data = &model.ChannelsIntegrationsData{}
	integration.AccessToken = null.String{}
	integration.RefreshToken = null.String{}
	integration.Enabled = false

	if err = c.Db.WithContext(ctx).Save(&integration).Error; err != nil {
		return nil, err
	}

	if err = c.sendGrpcEvent(ctx, integration.ID, integration.Enabled); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
