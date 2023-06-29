package integrations

import (
	"context"
	"github.com/satont/twir/libs/grpc/generated/api/integrations_valorant"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Integrations) IntegrationsValorantGet(ctx context.Context, empty *emptypb.Empty) (*integrations_valorant.GetResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Integrations) IntegrationsValorantUpdate(ctx context.Context, request *integrations_valorant.PostRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
