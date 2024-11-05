package integrations

import (
	"context"

	"github.com/guregu/null"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/integrations_donatepay"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Integrations) IntegrationsDonatepayGet(
	ctx context.Context,
	_ *emptypb.Empty,
) (*integrations_donatepay.GetResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	integration, err := c.getChannelIntegrationByService(
		ctx,
		model.IntegrationServiceDonatePay,
		dashboardId,
	)
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
	dashboardId := ctx.Value("dashboardId").(string)
	integration, err := c.getChannelIntegrationByService(
		ctx,
		model.IntegrationServiceDonatePay,
		dashboardId,
	)
	if err != nil {
		return nil, err
	}

	integration.APIKey = null.StringFrom(request.ApiKey)
	integration.Enabled = true
	if err = c.Db.WithContext(ctx).Save(&integration).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
