package integrations

import (
	"context"

	"github.com/guregu/null"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/integrations_donatepay"
	"github.com/twirapp/twir/libs/bus-core/integrations"
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
	entity, err := c.getChannelIntegrationByService(
		ctx,
		model.IntegrationServiceDonatePay,
		dashboardId,
	)
	if err != nil {
		return nil, err
	}

	entity.APIKey = null.StringFrom(request.ApiKey)
	entity.Enabled = true
	if err = c.Db.WithContext(ctx).Save(&entity).Error; err != nil {
		return nil, err
	}

	c.Bus.Integrations.Add.Publish(
		ctx, integrations.Request{
			ID: entity.ID,
		},
	)

	return &emptypb.Empty{}, nil
}
