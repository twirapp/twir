package integrations

import (
	"context"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/api/integrations_donatello"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Integrations) IntegrationsDonatelloGet(ctx context.Context, _ *emptypb.Empty) (*integrations_donatello.GetResponse, error) {
	integration, err := c.getChannelIntegrationByService(ctx, model.IntegrationServiceDonatello)
	if err != nil {
		return nil, err
	}

	return &integrations_donatello.GetResponse{
		IntegrationId: integration.ID,
	}, nil
}
