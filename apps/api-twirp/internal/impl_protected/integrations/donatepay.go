package integrations

import (
	"context"
	"github.com/guregu/null"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/api/integrations_donatepay"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Integrations) IntegrationsDonatepayGet(
	ctx context.Context,
	_ *emptypb.Empty,
) (*integrations_donatepay.GetResponse, error) {
	integration, err := c.getChannelIntegrationByService(ctx, model.IntegrationServiceDonateStream)
	if err != nil {
		return nil, err
	}

	return &integrations_donatepay.GetResponse{
		ApiKey: integration.APIKey.String,
	}, nil
}

func (c *Integrations) IntegrationsDonatepayPut(
	ctx context.Context,
	request *integrations_donatepay.PostRequest,
) (*emptypb.Empty, error) {
	integration, err := c.getChannelIntegrationByService(ctx, model.IntegrationServiceDonateStream)
	if err != nil {
		return nil, err
	}

	integration.APIKey = null.StringFrom(request.ApiKey)
	if err = c.Db.WithContext(ctx).Save(&integration).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
