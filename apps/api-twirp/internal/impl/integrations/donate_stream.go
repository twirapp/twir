package integrations

import (
	"context"
	"github.com/satont/tsuwari/libs/grpc/generated/api/integrations_donate_stream"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Integrations) IntegrationsDonateStreamGet(ctx context.Context, empty *emptypb.Empty) (*integrations_donate_stream.DonateStreamResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Integrations) IntegrationsDonateStreamPostSecret(ctx context.Context, request *integrations_donate_stream.DonateStreamPostSecretRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
