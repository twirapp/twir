package integrations

import (
	"context"
	"github.com/guregu/null"
	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/integrations_faceit"
	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/url"
)

func (c *Integrations) IntegrationsFaceitGetAuthLink(
	ctx context.Context, _ *emptypb.Empty,
) (*integrations_faceit.GetAuthLink, error) {
	integration, err := c.getIntegrationByService(ctx, model.IntegrationServiceFaceit)
	if err != nil {
		return nil, err
	}

	if !integration.ClientID.Valid || !integration.RedirectURL.Valid {
		return nil, twirp.NewError(twirp.Internal, "faceit not enabled on our side. Please be patient.")
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
	integration, err := c.getChannelIntegrationByService(ctx, model.IntegrationServiceFaceit, dashboardId)
	if err != nil {
		return nil, err
	}

	return &integrations_faceit.GetDataResponse{
		UserName: integration.Data.Name,
		Avatar:   integration.Data.Avatar,
		Game:     integration.Data.Game,
	}, nil
}

func (c *Integrations) IntegrationsFaceitPostCode(
	ctx context.Context, request *integrations_faceit.PostCodeRequest,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	integration, err := c.getChannelIntegrationByService(ctx, model.IntegrationServiceFaceit, dashboardId)
	if err != nil {
		return nil, err
	}

	tokensData := make(map[string]any)

	resp, err := req.
		C().EnableForceHTTP1().
		R().
		SetContext(ctx).
		SetFormData(
			map[string]string{
				"grant_type": "authorization_code",
				"code":       request.Code,
			},
		).
		SetSuccessResult(&tokensData).
		SetBasicAuth(integration.Integration.ClientID.String, integration.Integration.ClientSecret.String).
		Post("https://api.faceit.com/auth/v1/oauth/token")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() {
		return nil, resp.Err
	}

	integration.AccessToken = null.StringFrom(tokensData["access_token"].(string))
	integration.RefreshToken = null.StringFrom(tokensData["refresh_token"].(string))

	userInfoResult := make(map[string]any)
	resp, err = req.
		C().EnableForceHTTP1().
		R().
		SetContext(ctx).
		SetBearerAuthToken(integration.AccessToken.String).
		SetSuccessResult(&userInfoResult).
		Get("https://api.faceit.com/auth/v1/resources/userinfo")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() {
		return nil, resp.Err
	}

	integrationData := model.ChannelsIntegrationsData{
		UserId: lo.ToPtr(userInfoResult["guid"].(string)),
		Name:   lo.ToPtr(userInfoResult["nickname"].(string)),
		Game:   lo.ToPtr("csgo"),
	}

	profileResult := make(map[string]any)
	resp, err = req.
		C().EnableForceHTTP1().
		R().
		SetContext(ctx).
		SetBearerAuthToken(integration.Integration.APIKey.String).
		SetSuccessResult(&profileResult).
		Get("https://open.faceit.com/data/v4/players/" + *integrationData.UserId)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() {
		return nil, resp.Err
	}

	integrationData.Avatar = lo.ToPtr(profileResult["avatar"].(string))
	integration.Data = &integrationData

	if err = c.Db.WithContext(ctx).Save(integration).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Integrations) IntegrationsFaceitLogout(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboard_id").(string)
	integration, err := c.getChannelIntegrationByService(
		ctx, model.IntegrationServiceFaceit, dashboardId,
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

	return &emptypb.Empty{}, nil
}
