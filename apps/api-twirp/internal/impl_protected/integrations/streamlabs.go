package integrations

import (
	"context"
	"fmt"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/api/integrations_streamlabs"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/url"
)

type streamlabsTokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
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
		return nil, fmt.Errorf("spotify integration not configured")
	}

	link, _ := url.Parse("https://www.streamlabs.com/api/v1.0/authorize")

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
	integration, err := c.getChannelIntegrationByService(ctx, model.IntegrationServiceStreamLabs, dashboardId)
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
	//TODO implement me
	panic("implement me")
}

func (c *Integrations) IntegrationsStreamlabsLogout(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
