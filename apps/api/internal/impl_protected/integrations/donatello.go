package integrations

import (
	"context"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/integrations_donatello"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Integrations) IntegrationsDonatelloGet(
	ctx context.Context, _ *emptypb.Empty,
) (*integrations_donatello.GetResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	integration, err := c.getChannelIntegrationByService(
		ctx,
		model.IntegrationServiceDonatello,
		dashboardId,
	)
	if err != nil {
		return nil, err
	}

	return &integrations_donatello.GetResponse{
		IntegrationId: integration.ID,
	}, nil
}
