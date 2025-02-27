package integrations

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/guregu/null"
	"github.com/imroc/req/v3"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/integrations_vk"
	"google.golang.org/protobuf/types/known/emptypb"
)

type vkProfile struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	PhotoMaxOrig string `json:"photo_max_orig"`
	ID           int    `json:"id"`
}

type vkProfileResponse struct {
	Error *struct {
		Msg  string `json:"error_msg"`
		Code int    `json:"error_code"`
	}
	Response []vkProfile `json:"response"`
}

type vkTokensResponse struct {
	AccessToken string `json:"access_token"`
}

func (c *Integrations) IntegrationsVKGetAuthLink(
	ctx context.Context, _ *emptypb.Empty,
) (*integrations_vk.GetAuthLink, error) {
	integration, err := c.getIntegrationByService(ctx, model.IntegrationServiceVK)
	if err != nil {
		return nil, err
	}

	if !integration.ClientID.Valid || !integration.ClientSecret.Valid || !integration.RedirectURL.Valid {
		return nil, errors.New("vk not enabled on our side, please be patient")
	}

	link, _ := url.Parse("https://oauth.vk.com/authorize")

	q := link.Query()
	q.Add("client_id", integration.ClientID.String)
	q.Add("display", "page")
	q.Add("response_type", "code")
	q.Add("scope", "status offline")
	q.Add("redirect_uri", integration.RedirectURL.String)
	link.RawQuery = q.Encode()

	return &integrations_vk.GetAuthLink{
		Link: link.String(),
	}, nil
}

func (c *Integrations) IntegrationsVKGetData(
	ctx context.Context, _ *emptypb.Empty,
) (*integrations_vk.GetDataResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	integration, err := c.getChannelIntegrationByService(ctx, model.IntegrationServiceVK, dashboardId)
	if err != nil {
		return nil, err
	}

	return &integrations_vk.GetDataResponse{
		Avatar:   integration.Data.Avatar,
		UserName: integration.Data.UserName,
	}, nil
}

func (c *Integrations) IntegrationsVKPostCode(
	ctx context.Context, request *integrations_vk.PostCodeRequest,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	integration, err := c.getChannelIntegrationByService(ctx, model.IntegrationServiceVK, dashboardId)
	if err != nil {
		return nil, err
	}

	tokensData := &vkTokensResponse{}
	resp, err := req.R().
		SetContext(ctx).
		SetQueryParams(
			map[string]string{
				"grant_type":    "authorization_code",
				"client_id":     integration.Integration.ClientID.String,
				"client_secret": integration.Integration.ClientSecret.String,
				"redirect_uri":  integration.Integration.RedirectURL.String,
				"code":          request.Code,
			},
		).
		SetSuccessResult(tokensData).
		Get("https://oauth.vk.com/access_token")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("vk auth error: %s", resp.String())
	}

	profileData := &vkProfileResponse{}
	resp, err = req.R().
		SetContext(ctx).
		SetQueryParams(
			map[string]string{
				"v":            "5.131",
				"fields":       "photo_max_orig",
				"access_token": tokensData.AccessToken,
			},
		).
		SetSuccessResult(profileData).
		Get("https://api.vk.com/method/users.get")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("vk auth error: %s", resp.String())
	}

	userName := profileData.Response[0].FirstName + " " + profileData.Response[0].LastName
	integration.AccessToken = null.StringFrom(tokensData.AccessToken)
	integration.Data = &model.ChannelsIntegrationsData{
		Avatar:   &profileData.Response[0].PhotoMaxOrig,
		UserName: &userName,
	}
	integration.Enabled = true
	if err = c.Db.WithContext(ctx).Save(integration).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Integrations) IntegrationsVKLogout(
	ctx context.Context, _ *emptypb.Empty,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	integration, err := c.getChannelIntegrationByService(
		ctx, model.IntegrationServiceVK, dashboardId,
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

	return &emptypb.Empty{}, nil
}
