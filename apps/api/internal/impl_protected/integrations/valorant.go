package integrations

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/guregu/null"
	"github.com/imroc/req/v3"
	"github.com/twirapp/twir/apps/api/internal/helpers"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/integrations_valorant"
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

	tokenResponse := ValorantTokenResponse{}
	resp, err := req.R().
		SetContext(ctx).
		SetFormData(
			map[string]string{
				// "client_assertion_type": "urn:ietf:params:oauth:client-assertion-type:jwt-bearer",
				// "client_assertion":      "",
				"grant_type":   "authorization_code",
				"redirect_uri": integration.Integration.RedirectURL.String,
				"code":         request.GetCode(),
			},
		).
		SetBasicAuth(
			integration.Integration.ClientID.String,
			integration.Integration.ClientSecret.String,
		).
		SetSuccessResult(&tokenResponse).
		SetContentType("application/x-www-form-urlencoded").
		Post("https://auth.riotgames.com/token")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("failed to get valorant tokens: %s", resp.String())
	}

	accountResponse := ValorantAccountResp{}
	accountReq, err := req.
		SetSuccessResult(&accountResponse).
		SetBearerAuthToken(tokenResponse.AccessToken).
		Get(api_base + "/riot/account/v1/accounts/me")
	if err != nil {
		return nil, err
	}
	if !accountReq.IsSuccessState() {
		return nil, fmt.Errorf("cannot get valorant account info: %s", accountReq.String())
	}

	shardResponse := ValorantShardResponse{}
	shardReq, err := req.
		SetHeader("X-Riot-Token", integration.Integration.APIKey.String).
		SetSuccessResult(&shardResponse).
		Get(
			fmt.Sprintf(
				"%s/riot/account/v1/active-shards/by-game/%s/by-puuid/%s",
				api_base,
				"val",
				accountResponse.Puuid,
			),
		)
	if err != nil {
		return nil, err
	}
	if !shardReq.IsSuccessState() {
		return nil, fmt.Errorf(
			"cannot get valorant shard info: %v, %s",
			shardReq.StatusCode,
			shardReq.String(),
		)
	}

	henrikResponse := ValorantHenrikResponse{}
	henrikReq, err := req.
		SetSuccessResult(&henrikResponse).
		SetHeader("Authorization", c.Config.ValorantHenrikApiKey).
		Get(
			fmt.Sprintf(
				"https://api.henrikdev.xyz/valorant/v1/account/%s/%s",
				accountResponse.UserName,
				accountResponse.TagLine,
			),
		)
	if err != nil {
		return nil, err
	}
	if !henrikReq.IsSuccessState() {
		return nil, fmt.Errorf(
			"cannot get valorant shard info: %d %s",
			henrikReq.StatusCode,
			henrikReq.String(),
		)
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
