package integrations

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/guregu/null"
	"github.com/twirapp/twir/apps/api/internal/helpers"
	"github.com/twirapp/twir/libs/api/messages/integrations_valorant"
	model "github.com/twirapp/twir/libs/gomodels"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ValorantTokenResponse struct {
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
	SubSid       string `json:"sub_sid"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type ValorantAccountResp struct {
	Puuid    string `json:"puuid"`
	UserName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

type ValorantShardResponse struct {
	Puuid       string `json:"puuid"`
	Game        string `json:"game"`
	ActiveShard string `json:"activeShard"`
}

type ValorantHenrikResponse struct {
	Data struct {
		Puuid        string `json:"puuid"`
		Region       string `json:"region"`
		Name         string `json:"name"`
		Tag          string `json:"tag"`
		AccountLevel int    `json:"account_level"`
	} `json:"data"`
	Status int `json:"status"`
}

const api_base = "https://europe.api.riotgames.com"

func (c *Integrations) IntegrationsValorantGetAuthLink(
	ctx context.Context,
	_ *emptypb.Empty,
) (*integrations_valorant.GetAuthLink, error) {
	integration, err := c.getIntegrationByService(ctx, model.IntegrationServiceValorant)
	if err != nil {
		return nil, err
	}

	if !integration.ClientID.Valid || !integration.ClientSecret.Valid || !integration.RedirectURL.Valid {
		return nil, errors.New("valorant not enabled on our side, please be patient")
	}

	link, _ := url.Parse("https://auth.riotgames.com/authorize")

	q := link.Query()
	q.Add("response_type", "code")
	q.Add("client_id", integration.ClientID.String)
	q.Add("scope", strings.Join([]string{"openid", "offline_access", "cpid"}, "+"))
	q.Add("redirect_uri", integration.RedirectURL.String)
	link.RawQuery = q.Encode()

	return &integrations_valorant.GetAuthLink{
		Link: link.String(),
	}, nil
}

func (c *Integrations) IntegrationsValorantGetData(
	ctx context.Context,
	_ *emptypb.Empty,
) (*integrations_valorant.GetDataResponse, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	integration, err := c.getChannelIntegrationByService(
		ctx,
		model.IntegrationServiceValorant,
		dashboardId,
	)
	if err != nil {
		return nil, err
	}

	return &integrations_valorant.GetDataResponse{
		UserName: integration.Data.UserName,
		Avatar:   integration.Data.Avatar,
	}, nil
}

func (c *Integrations) IntegrationsValorantPostCode(
	ctx context.Context,
	request *integrations_valorant.PostCodeRequest,
) (*emptypb.Empty, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	integration, err := c.getChannelIntegrationByService(
		ctx,
		model.IntegrationServiceValorant,
		dashboardId,
	)
	if err != nil {
		return nil, err
	}

	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("redirect_uri", integration.Integration.RedirectURL.String)
	formData.Set("code", request.GetCode())

	tokenReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://auth.riotgames.com/token",
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
		return nil, fmt.Errorf("failed to get valorant tokens: %s", string(bodyBytes))
	}

	tokenBodyBytes, err := io.ReadAll(tokenResp.Body)
	if err != nil {
		return nil, err
	}

	tokenResponse := ValorantTokenResponse{}
	if err := json.Unmarshal(tokenBodyBytes, &tokenResponse); err != nil {
		return nil, err
	}

	accountReq, err := http.NewRequestWithContext(ctx, http.MethodGet, api_base+"/riot/account/v1/accounts/me", nil)
	if err != nil {
		return nil, err
	}
	accountReq.Header.Set("Authorization", "Bearer "+tokenResponse.AccessToken)

	accountResp, err := http.DefaultClient.Do(accountReq)
	if err != nil {
		return nil, err
	}
	defer accountResp.Body.Close()

	if accountResp.StatusCode < 200 || accountResp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(accountResp.Body)
		return nil, fmt.Errorf("cannot get valorant account info: %s", string(bodyBytes))
	}

	accountBodyBytes, err := io.ReadAll(accountResp.Body)
	if err != nil {
		return nil, err
	}

	accountResponse := ValorantAccountResp{}
	if err := json.Unmarshal(accountBodyBytes, &accountResponse); err != nil {
		return nil, err
	}

	shardUrl := fmt.Sprintf(
		"%s/riot/account/v1/active-shards/by-game/%s/by-puuid/%s",
		api_base,
		"val",
		accountResponse.Puuid,
	)
	shardReq, err := http.NewRequestWithContext(ctx, http.MethodGet, shardUrl, nil)
	if err != nil {
		return nil, err
	}
	shardReq.Header.Set("X-Riot-Token", integration.Integration.APIKey.String)

	shardResp, err := http.DefaultClient.Do(shardReq)
	if err != nil {
		return nil, err
	}
	defer shardResp.Body.Close()

	if shardResp.StatusCode < 200 || shardResp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(shardResp.Body)
		return nil, fmt.Errorf(
			"cannot get valorant shard info: %v, %s",
			shardResp.StatusCode,
			string(bodyBytes),
		)
	}

	shardBodyBytes, err := io.ReadAll(shardResp.Body)
	if err != nil {
		return nil, err
	}

	shardResponse := ValorantShardResponse{}
	if err := json.Unmarshal(shardBodyBytes, &shardResponse); err != nil {
		return nil, err
	}

	henrikUrl := fmt.Sprintf(
		"https://api.henrikdev.xyz/valorant/v1/account/%s/%s",
		accountResponse.UserName,
		accountResponse.TagLine,
	)
	henrikReq, err := http.NewRequestWithContext(ctx, http.MethodGet, henrikUrl, nil)
	if err != nil {
		return nil, err
	}
	henrikReq.Header.Set("Authorization", c.Config.ValorantHenrikApiKey)

	henrikResp, err := http.DefaultClient.Do(henrikReq)
	if err != nil {
		return nil, err
	}
	defer henrikResp.Body.Close()

	if henrikResp.StatusCode < 200 || henrikResp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(henrikResp.Body)
		return nil, fmt.Errorf(
			"cannot get valorant shard info: %d %s",
			henrikResp.StatusCode,
			string(bodyBytes),
		)
	}

	henrikBodyBytes, err := io.ReadAll(henrikResp.Body)
	if err != nil {
		return nil, err
	}

	henrikResponse := ValorantHenrikResponse{}
	if err := json.Unmarshal(henrikBodyBytes, &henrikResponse); err != nil {
		return nil, err
	}

	userName := fmt.Sprintf(
		"%s#%s",
		accountResponse.UserName,
		accountResponse.TagLine,
	)

	integration.AccessToken = null.StringFrom(tokenResponse.AccessToken)
	integration.RefreshToken = null.StringFrom(tokenResponse.RefreshToken)
	integration.Data.ValorantActiveRegion = &shardResponse.ActiveShard
	integration.Data.UserName = &userName
	integration.Data.ValorantPuuid = &henrikResponse.Data.Puuid
	integration.Enabled = true

	if err = c.Db.WithContext(ctx).Save(integration).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Integrations) IntegrationsValorantLogout(
	ctx context.Context,
	_ *emptypb.Empty,
) (*emptypb.Empty, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	integration, err := c.getChannelIntegrationByService(
		ctx,
		model.IntegrationServiceValorant,
		dashboardId,
	)
	if err != nil {
		return nil, err
	}

	integration.Data = &model.ChannelsIntegrationsData{}
	integration.APIKey = null.String{}
	integration.AccessToken = null.String{}
	integration.RefreshToken = null.String{}
	integration.Enabled = false

	if err = c.Db.WithContext(ctx).Save(&integration).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
