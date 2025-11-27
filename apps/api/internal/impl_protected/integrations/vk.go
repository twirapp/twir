package integrations

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/guregu/null"
	"github.com/twirapp/twir/libs/api/messages/integrations_vk"
	model "github.com/twirapp/twir/libs/gomodels"
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

	tokenUrl, _ := url.Parse("https://oauth.vk.com/access_token")
	q := tokenUrl.Query()
	q.Set("grant_type", "authorization_code")
	q.Set("client_id", integration.Integration.ClientID.String)
	q.Set("client_secret", integration.Integration.ClientSecret.String)
	q.Set("redirect_uri", integration.Integration.RedirectURL.String)
	q.Set("code", request.Code)
	tokenUrl.RawQuery = q.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, tokenUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("vk auth error: %s", string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	tokensData := &vkTokensResponse{}
	if err := json.Unmarshal(bodyBytes, tokensData); err != nil {
		return nil, err
	}

	profileUrl, _ := url.Parse("https://api.vk.com/method/users.get")
	pq := profileUrl.Query()
	pq.Set("v", "5.131")
	pq.Set("fields", "photo_max_orig")
	pq.Set("access_token", tokensData.AccessToken)
	profileUrl.RawQuery = pq.Encode()

	profileReq, err := http.NewRequestWithContext(ctx, http.MethodGet, profileUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	profileResp, err := http.DefaultClient.Do(profileReq)
	if err != nil {
		return nil, err
	}
	defer profileResp.Body.Close()

	if profileResp.StatusCode < 200 || profileResp.StatusCode >= 300 {
		profileBodyBytes, _ := io.ReadAll(profileResp.Body)
		return nil, fmt.Errorf("vk auth error: %s", string(profileBodyBytes))
	}

	profileBodyBytes, err := io.ReadAll(profileResp.Body)
	if err != nil {
		return nil, err
	}

	profileData := &vkProfileResponse{}
	if err := json.Unmarshal(profileBodyBytes, profileData); err != nil {
		return nil, err
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
