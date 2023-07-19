package integrations

import (
	"context"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/integrations_valorant"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Integrations) IntegrationsValorantGet(
	ctx context.Context, _ *emptypb.Empty,
) (*integrations_valorant.GetResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	integration, err := c.getChannelIntegrationByService(ctx, model.IntegrationServiceValorant, dashboardId)
	if err != nil {
		return nil, err
	}

	return &integrations_valorant.GetResponse{
		UserName: integration.Data.UserName,
	}, nil
}

func (c *Integrations) IntegrationsValorantUpdate(
	ctx context.Context,
	request *integrations_valorant.PostRequest,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	integration, err := c.getChannelIntegrationByService(ctx, model.IntegrationServiceValorant, dashboardId)
	if err != nil {
		return nil, err
	}

	integration.Data.UserName = &request.UserName
	if err = c.Db.WithContext(ctx).Save(integration).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
