package integrations

import (
	"context"
	"time"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/integrations_donate_stream"
	"google.golang.org/protobuf/types/known/emptypb"
)

const donateStreamConfirmationKey = "donate_stream_confirmation"

func (c *Integrations) IntegrationsDonateStreamGet(
	ctx context.Context, _ *emptypb.Empty,
) (*integrations_donate_stream.DonateStreamResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	integration, err := c.getChannelIntegrationByService(ctx, model.IntegrationServiceDonateStream, dashboardId)
	if err != nil {
		return nil, err
	}

	return &integrations_donate_stream.DonateStreamResponse{
		IntegrationId: integration.ID,
	}, nil
}

func (c *Integrations) IntegrationsDonateStreamPostSecret(
	ctx context.Context, request *integrations_donate_stream.DonateStreamPostSecretRequest,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	integration, err := c.getChannelIntegrationByService(ctx, model.IntegrationServiceDonateStream, dashboardId)
	if err != nil {
		return nil, err
	}

	if err = c.Redis.
		Set(ctx, donateStreamConfirmationKey+integration.IntegrationID, request.Secret, 1*time.Hour).
		Err(); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
