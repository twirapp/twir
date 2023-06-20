package integrations

import (
	"context"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/api/integrations_donationalerts"
	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/url"
)

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
	ctx context.Context, empty *emptypb.Empty,
) (*integrations_donationalerts.GetDataResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Integrations) IntegrationsDonationAlertsPostCode(
	ctx context.Context, request *integrations_donationalerts.PostCodeRequest,
) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Integrations) IntegrationsDonationAlertsLogout(ctx context.Context, empty *emptypb.Empty) (
	*emptypb.Empty, error,
) {
	//TODO implement me
	panic("implement me")
}
