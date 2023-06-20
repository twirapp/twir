package integrations

import (
	"context"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/api/integrations_donate_stream"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

const donateStreamConfirmationKey = "donate_stream_confirmation"

func (c *Integrations) IntegrationsDonateStreamGet(ctx context.Context, _ *emptypb.Empty) (*integrations_donate_stream.DonateStreamResponse, error) {
	integration, err := c.getChannelIntegrationByService(ctx, model.IntegrationServiceDonateStream)
	if err != nil {
		return nil, err
	}

	return &integrations_donate_stream.DonateStreamResponse{
		IntegrationId: integration.ID,
	}, nil
}

func (c *Integrations) IntegrationsDonateStreamPostSecret(ctx context.Context, request *integrations_donate_stream.DonateStreamPostSecretRequest) (*emptypb.Empty, error) {
	integration, err := c.getChannelIntegrationByService(ctx, model.IntegrationServiceDonateStream)
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
