package integrations

import (
	"context"
	"github.com/guregu/null"
	"github.com/imroc/req/v3"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/api/integrations_donationalerts"
	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/url"
)

type donationAlertsTokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type donationAlertsProfileResponse struct {
	Data struct {
		Name   string `json:"name"`
		Code   string `json:"code"`
		Avatar string `json:"avatar"`
	} `json:"data"`
}

func (c *Integrations) IntegrationsDonationAlertsGetAuthLink(
	ctx context.Context, _ *emptypb.Empty,
) (*integrations_donationalerts.GetAuthLink, error) {
	integration, err := c.getIntegrationByService(ctx, model.IntegrationServiceDonationAlerts)
	if err != nil {
		return nil, err
	}

	if !integration.ClientID.Valid || !integration.ClientSecret.Valid || !integration.RedirectURL.Valid {
		return nil, twirp.Internal.Error("internal error")
	}

	u, _ := url.Parse("https://www.donationalerts.com/oauth/authorize")

	q := u.Query()
	q.Add("client_id", integration.ClientID.String)
	q.Add("response_type", "code")
	q.Add("scope", "oauth-user-show oauth-donation-subscribe")
	q.Add("redirect_uri", integration.RedirectURL.String)
	u.RawQuery = q.Encode()

	str := u.String()

	return &integrations_donationalerts.GetAuthLink{Link: str}, nil
}

func (c *Integrations) IntegrationsDonationAlertsGetData(
	ctx context.Context, _ *emptypb.Empty,
) (*integrations_donationalerts.GetDataResponse, error) {
	dashboardId := ctx.Value("dashboard_id").(string)
	integration, err := c.getChannelIntegrationByService(
		ctx, model.IntegrationServiceDonationAlerts, dashboardId,
	)
	if err != nil {
		return nil, err
	}

	return &integrations_donationalerts.GetDataResponse{
		UserName: integration.Data.Name,
		Avatar:   integration.Data.Avatar,
	}, nil
}

func (c *Integrations) IntegrationsDonationAlertsPostCode(
	ctx context.Context, request *integrations_donationalerts.PostCodeRequest,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboard_id").(string)
	integration, err := c.getChannelIntegrationByService(
		ctx, model.IntegrationServiceDonationAlerts, dashboardId,
	)
	if err != nil {
		return nil, err
	}

	data := donationAlertsTokensResponse{}
	resp, err := req.R().
		SetContext(ctx).
		SetFormData(
			map[string]string{
				"grant_type":    "authorization_code",
				"client_id":     integration.Integration.ClientID.String,
				"client_secret": integration.Integration.ClientSecret.String,
				"redirect_uri":  integration.Integration.RedirectURL.String,
				"code":          request.Code,
			},
		).
		SetSuccessResult(&data).
		SetContentType("application/x-www-form-urlencoded").
		Post("https://www.donationalerts.com/oauth/token")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() {
		return nil, twirp.InternalErrorWith(resp.Err)
	}

	profile := donationAlertsProfileResponse{}
	profileResp, err := req.R().
		SetContext(ctx).
		SetSuccessResult(&profile).
		SetBearerAuthToken(data.AccessToken).
		Get("https://www.donationalerts.com/api/v1/user/oauth")
	if err != nil {
		return nil, err
	}
	if !profileResp.IsSuccessState() {
		return nil, twirp.InternalErrorWith(resp.Err)
	}

	integration.Data = &model.ChannelsIntegrationsData{
		Name:   &profile.Data.Name,
		Code:   &profile.Data.Code,
		Avatar: &profile.Data.Avatar,
	}
	integration.AccessToken = null.StringFrom(data.AccessToken)
	integration.RefreshToken = null.StringFrom(data.RefreshToken)

	if err = c.Db.WithContext(ctx).Save(integration).Error; err != nil {
		return nil, err
	}

	if err = c.sendGrpcEvent(ctx, integration.ID, integration.Enabled); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Integrations) IntegrationsDonationAlertsLogout(ctx context.Context, _ *emptypb.Empty) (
	*emptypb.Empty, error,
) {
	dashboardId := ctx.Value("dashboard_id").(string)
	integration, err := c.getChannelIntegrationByService(
		ctx, model.IntegrationServiceDonationAlerts, dashboardId,
	)
	if err != nil {
		return nil, err
	}

	integration.Data = nil
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
