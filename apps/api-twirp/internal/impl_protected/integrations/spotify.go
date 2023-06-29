package integrations

import (
	"context"
	"fmt"
	"github.com/guregu/null"
	"github.com/imroc/req/v3"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/integrations_spotify"
	"github.com/satont/twir/libs/integrations/spotify"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/url"
)

type spotifyTokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

func (c *Integrations) IntegrationsSpotifyGetAuthLink(
	ctx context.Context,
	_ *emptypb.Empty,
) (*integrations_spotify.GetAuthLink, error) {
	integration, err := c.getIntegrationByService(ctx, model.IntegrationServiceSpotify)
	if err != nil {
		return nil, err
	}

	if !integration.ClientID.Valid || !integration.ClientSecret.Valid || !integration.RedirectURL.Valid {
		return nil, fmt.Errorf("spotify integration not configured")
	}

	link, _ := url.Parse("https://accounts.spotify.com/authorize")

	q := link.Query()
	q.Add("response_type", "code")
	q.Add("client_id", integration.ClientID.String)
	q.Add("scope", "user-read-currently-playing")
	q.Add("redirect_uri", integration.RedirectURL.String)
	link.RawQuery = q.Encode()

	return &integrations_spotify.GetAuthLink{
		Link: link.String(),
	}, nil
}

func (c *Integrations) IntegrationsSpotifyGetData(
	ctx context.Context,
	_ *emptypb.Empty,
) (*integrations_spotify.GetDataResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	integration, err := c.getChannelIntegrationByService(ctx, model.IntegrationServiceSpotify, dashboardId)
	if err != nil {
		return nil, err
	}

	return &integrations_spotify.GetDataResponse{
		UserName: integration.Data.UserName,
		Avatar:   integration.Data.Avatar,
	}, nil
}

func (c *Integrations) IntegrationsSpotifyPostCode(
	ctx context.Context,
	request *integrations_spotify.PostCodeRequest,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	integration, err := c.getChannelIntegrationByService(ctx, model.IntegrationServiceSpotify, dashboardId)
	if err != nil {
		return nil, err
	}

	data := spotifyTokensResponse{}
	resp, err := req.R().
		SetContext(ctx).
		SetFormData(
			map[string]string{
				"grant_type":   "authorization_code",
				"redirect_uri": integration.Integration.RedirectURL.String,
				"code":         request.Code,
			},
		).
		SetBasicAuth(integration.Integration.ClientID.String, integration.Integration.ClientSecret.String).
		SetSuccessResult(&data).
		SetContentType("application/x-www-form-urlencoded").
		Post("https://accounts.spotify.com/api/token")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("failed to get spotify tokens: %s", resp.String())
	}

	integration.AccessToken = null.StringFrom(data.AccessToken)
	integration.RefreshToken = null.StringFrom(data.RefreshToken)

	userSpotify := spotify.New(integration, c.Db)
	profile, err := userSpotify.GetProfile()
	if err != nil {
		return nil, err
	}

	integration.Data.UserName = &profile.DisplayName
	integration.Data.Avatar = &profile.Images[0].URL

	if err = c.Db.WithContext(ctx).Save(integration).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Integrations) IntegrationsSpotifyLogout(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboard_id").(string)
	integration, err := c.getChannelIntegrationByService(
		ctx, model.IntegrationServiceSpotify, dashboardId,
	)
	if err != nil {
		return nil, err
	}

	integration.Data = nil
	integration.APIKey = null.String{}

	if err = c.Db.WithContext(ctx).Save(&integration).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
