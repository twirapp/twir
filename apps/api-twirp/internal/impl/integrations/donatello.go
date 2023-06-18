package integrations

import (
	"context"
	"github.com/satont/tsuwari/libs/grpc/generated/api/integrations_donatello"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Integrations) IntegrationsDonatelloGet(ctx context.Context, empty *emptypb.Empty) (*integrations_donatello.GetResponse, error) {
	//TODO implement me
	panic("implement me")
}
