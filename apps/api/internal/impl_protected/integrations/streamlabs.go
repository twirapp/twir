package integrations

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/guregu/null"
	"github.com/twirapp/twir/libs/api/messages/integrations_streamlabs"
	"github.com/twirapp/twir/libs/bus-core/integrations"
	model "github.com/twirapp/twir/libs/gomodels"
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

	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("client_id", channelIntegration.Integration.ClientID.String)
	formData.Set("client_secret", channelIntegration.Integration.ClientSecret.String)
	formData.Set("redirect_uri", channelIntegration.Integration.RedirectURL.String)
	formData.Set("code", request.Code)

	tokenReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://streamlabs.com/api/v2.0/token",
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		return nil, err
	}
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tokenResp, err := http.DefaultClient.Do(tokenReq)
	if err != nil {
		return nil, err
	}
	defer tokenResp.Body.Close()

	if tokenResp.StatusCode < 200 || tokenResp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(tokenResp.Body)
		return nil, fmt.Errorf("streamlabs token request failed: %s", string(bodyBytes))
	}

	tokenBodyBytes, err := io.ReadAll(tokenResp.Body)
	if err != nil {
		return nil, err
	}

	tokensData := streamlabsTokensResponse{}
	if err := json.Unmarshal(tokenBodyBytes, &tokensData); err != nil {
		return nil, err
	}

	userReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://streamlabs.com/api/v2.0/user", nil)
	if err != nil {
		return nil, err
	}
	userReq.Header.Set("Authorization", "Bearer "+tokensData.AccessToken)

	userResp, err := http.DefaultClient.Do(userReq)
	if err != nil {
		return nil, err
	}
	defer userResp.Body.Close()

	if userResp.StatusCode < 200 || userResp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(userResp.Body)
		return nil, fmt.Errorf("streamlabs token request failed: %s", string(bodyBytes))
	}

	userBodyBytes, err := io.ReadAll(userResp.Body)
	if err != nil {
		return nil, err
	}

	profileData := &streamlabsProfileResponse{}
	if err := json.Unmarshal(userBodyBytes, profileData); err != nil {
		return nil, err
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

	if err = c.sendBusEvent(ctx, channelIntegration.ID, channelIntegration.Enabled, integrations.Streamlabs); err != nil {
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

	if err = c.sendBusEvent(ctx, integration.ID, integration.Enabled, integrations.Streamlabs); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
