package integrations

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/guregu/null"
	"github.com/imroc/req/v3"
	"github.com/kr/pretty"
	"github.com/satont/twir/apps/api/internal/helpers"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/integrations_valorant"
	"google.golang.org/protobuf/types/known/emptypb"
)

type valorantTokenResponse struct {
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
	SubSid       string `json:"sub_sid"`
	AccessToken  string `json:"access_token"`
}

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
	dashboardId := ctx.Value("dashboardId").(string)
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

	data := valorantTokenResponse{}
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
		SetSuccessResult(&data).
		SetContentType("application/x-www-form-urlencoded").
		Post("https://auth.riotgames.com/token")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("failed to get valorant tokens: %s", resp.String())
	}

	integration.AccessToken = null.StringFrom(data.AccessToken)
	integration.RefreshToken = null.StringFrom(data.RefreshToken)

	if err = c.Db.WithContext(ctx).Save(integration).Error; err != nil {
		return nil, err
	}

	type valAccountResp struct {
		Puuid string `json:"puuid"`
	}

	v := &valAccountResp{}

	req.SetSuccessResult(v).SetBearerAuthToken(data.AccessToken).Get(
		"https://europe.api." +
			"riotgames.com/riot/account/v1/accounts/me",
	)

	pretty.Println(v)

	q, err := req.SetBearerAuthToken(data.AccessToken).Get(
		"https://eu.api.riotgames.com/val/match/v1/matchlists/by-puuid/" + v.Puuid,
	)
	fmt.Println(err, q)

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
