package integrations

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/guregu/null"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/api/messages/integrations_faceit"
	model "github.com/twirapp/twir/libs/gomodels"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Integrations) IntegrationsFaceitGetAuthLink(
	ctx context.Context, _ *emptypb.Empty,
) (*integrations_faceit.GetAuthLink, error) {
	integration, err := c.getIntegrationByService(ctx, model.IntegrationServiceFaceit)
	if err != nil {
		return nil, err
	}

	if !integration.ClientID.Valid || !integration.ClientSecret.Valid || !integration.RedirectURL.Valid {
		return nil, errors.New("faceit not enabled on our side, please be patient")
	}

	link, _ := url.Parse("https://cdn.faceit.com/widgets/sso/index.html")

	q := link.Query()
	q.Add("response_type", "code")
	q.Add("client_id", integration.ClientID.String)
	q.Add("redirect_popup", integration.RedirectURL.String)
	link.RawQuery = q.Encode()

	return &integrations_faceit.GetAuthLink{Link: link.String()}, nil
}

func (c *Integrations) IntegrationsFaceitGetData(
	ctx context.Context, _ *emptypb.Empty,
) (*integrations_faceit.GetDataResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	integration, err := c.getChannelIntegrationByService(
		ctx,
		model.IntegrationServiceFaceit,
		dashboardId,
	)
	if err != nil {
		return nil, err
	}

	return &integrations_faceit.GetDataResponse{
		UserName: integration.Data.UserName,
		Avatar:   integration.Data.Avatar,
		Game:     integration.Data.Game,
	}, nil
}

func (c *Integrations) IntegrationsFaceitUpdate(
	ctx context.Context,
	req *integrations_faceit.UpdateDataRequest,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	integration, err := c.getChannelIntegrationByService(
		ctx,
		model.IntegrationServiceFaceit,
		dashboardId,
	)
	if err != nil {
		return nil, err
	}

	integration.Data.Game = lo.ToPtr(req.Game)

	if err = c.Db.WithContext(ctx).Save(integration).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Integrations) IntegrationsFaceitPostCode(
	ctx context.Context, request *integrations_faceit.PostCodeRequest,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	integration, err := c.getChannelIntegrationByService(
		ctx,
		model.IntegrationServiceFaceit,
		dashboardId,
	)
	if err != nil {
		return nil, err
	}

	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("code", request.Code)

	tokenReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.faceit.com/auth/v1/oauth/token",
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		return nil, err
	}
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	authStr := base64.StdEncoding.EncodeToString(
		[]byte(integration.Integration.ClientID.String + ":" + integration.Integration.ClientSecret.String),
	)
	tokenReq.Header.Set("Authorization", "Basic "+authStr)

	tokenResp, err := http.DefaultClient.Do(tokenReq)
	if err != nil {
		return nil, err
	}
	defer tokenResp.Body.Close()

	if tokenResp.StatusCode < 200 || tokenResp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(tokenResp.Body)
		return nil, errors.New(string(bodyBytes))
	}

	tokenBodyBytes, err := io.ReadAll(tokenResp.Body)
	if err != nil {
		return nil, err
	}

	tokensData := make(map[string]any)
	if err := json.Unmarshal(tokenBodyBytes, &tokensData); err != nil {
		return nil, err
	}

	integration.AccessToken = null.StringFrom(tokensData["access_token"].(string))
	integration.RefreshToken = null.StringFrom(tokensData["refresh_token"].(string))

	userInfoReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.faceit.com/auth/v1/resources/userinfo", nil)
	if err != nil {
		return nil, err
	}
	userInfoReq.Header.Set("Authorization", "Bearer "+integration.AccessToken.String)

	userInfoResp, err := http.DefaultClient.Do(userInfoReq)
	if err != nil {
		return nil, err
	}
	defer userInfoResp.Body.Close()

	if userInfoResp.StatusCode < 200 || userInfoResp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(userInfoResp.Body)
		return nil, errors.New(string(bodyBytes))
	}

	userInfoBodyBytes, err := io.ReadAll(userInfoResp.Body)
	if err != nil {
		return nil, err
	}

	userInfoResult := make(map[string]any)
	if err := json.Unmarshal(userInfoBodyBytes, &userInfoResult); err != nil {
		return nil, err
	}

	integrationData := model.ChannelsIntegrationsData{
		UserId:   lo.ToPtr(userInfoResult["guid"].(string)),
		UserName: lo.ToPtr(userInfoResult["nickname"].(string)),
		Game:     lo.ToPtr("cs2"),
	}

	profileReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://open.faceit.com/data/v4/players/"+*integrationData.UserId, nil)
	if err != nil {
		return nil, err
	}
	profileReq.Header.Set("Authorization", "Bearer "+integration.Integration.APIKey.String)

	profileResp, err := http.DefaultClient.Do(profileReq)
	if err != nil {
		return nil, err
	}
	defer profileResp.Body.Close()

	if profileResp.StatusCode < 200 || profileResp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(profileResp.Body)
		return nil, errors.New(string(bodyBytes))
	}

	profileBodyBytes, err := io.ReadAll(profileResp.Body)
	if err != nil {
		return nil, err
	}

	profileResult := make(map[string]any)
	if err := json.Unmarshal(profileBodyBytes, &profileResult); err != nil {
		return nil, err
	}

	integrationData.Avatar = lo.ToPtr(profileResult["avatar"].(string))
	integration.Data = &integrationData
	integration.Enabled = true

	if err = c.Db.WithContext(ctx).Save(integration).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Integrations) IntegrationsFaceitLogout(
	ctx context.Context,
	empty *emptypb.Empty,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	integration, err := c.getChannelIntegrationByService(
		ctx, model.IntegrationServiceFaceit, dashboardId,
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
