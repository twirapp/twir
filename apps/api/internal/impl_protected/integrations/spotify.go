package integrations

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/imroc/req/v3"
	deprecatedgormmodel "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/integrations_spotify"
	"github.com/twirapp/twir/libs/integrations/spotify"
	channelsintegrationsspotify "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	"github.com/twirapp/twir/libs/repositories/channels_integrations_spotify/model"
	"google.golang.org/protobuf/types/known/emptypb"
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
	integration, err := c.getIntegrationByService(ctx, deprecatedgormmodel.IntegrationServiceSpotify)
	if err != nil {
		return nil, err
	}

	if !integration.ClientID.Valid || !integration.ClientSecret.Valid || !integration.RedirectURL.Valid {
		return nil, errors.New("spotify not enabled on our side, please be patient")
	}

	link, _ := url.Parse("https://accounts.spotify.com/authorize")

	q := link.Query()
	q.Add("response_type", "code")
	q.Add("client_id", integration.ClientID.String)
	q.Add("scope", "user-read-currently-playing user-read-playback-state")
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
	integration, err := c.SpotifyRepo.GetByChannelID(ctx, dashboardId)
	if err != nil {
		return nil, err
	}

	return &integrations_spotify.GetDataResponse{
		UserName: &integration.Username,
		Avatar:   &integration.AvatarURI,
	}, nil
}

func (c *Integrations) IntegrationsSpotifyPostCode(
	ctx context.Context,
	request *integrations_spotify.PostCodeRequest,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	integration, err := c.getIntegrationByService(
		ctx,
		deprecatedgormmodel.IntegrationServiceSpotify,
	)
	if err != nil {
		return nil, err
	}

	data := spotifyTokensResponse{}
	resp, err := req.R().
		SetContext(ctx).
		SetFormData(
			map[string]string{
				"grant_type":   "authorization_code",
				"redirect_uri": integration.RedirectURL.String,
				"code":         request.Code,
			},
		).
		SetBasicAuth(
			integration.ClientID.String,
			integration.ClientSecret.String,
		).
		SetSuccessResult(&data).
		SetContentType("application/x-www-form-urlencoded").
		Post("https://accounts.spotify.com/api/token")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("failed to get spotify tokens: %s", resp.String())
	}

	createInput := channelsintegrationsspotify.CreateInput{
		ChannelID:    dashboardId,
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
		Scopes:       []string{"user-read-currently-playing", "user-read-playback-state"},
	}

	userSpotify := spotify.New(
		*integration,
		model.ChannelIntegrationSpotify{
			AccessToken:  data.AccessToken,
			RefreshToken: data.RefreshToken,
		}, c.SpotifyRepo,
	)
	profile, err := userSpotify.GetProfile(ctx)
	if err != nil {
		return nil, err
	}

	createInput.AvatarURI = profile.Images[0].URL
	createInput.Username = profile.DisplayName

	if _, err := c.SpotifyRepo.Create(ctx, createInput); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Integrations) IntegrationsSpotifyLogout(
	ctx context.Context,
	_ *emptypb.Empty,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	integration, err := c.SpotifyRepo.GetByChannelID(ctx, dashboardId)
	if err != nil {
		return nil, err
	}
	if integration.AccessToken == "" {
		return nil, errors.New("not found")
	}

	if err = c.SpotifyRepo.Delete(ctx, integration.ID); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
