package integrations

import (
	"context"
	"github.com/satont/tsuwari/libs/grpc/generated/api/integrations_faceit"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Integrations) IntegrationsFaceitGetAuthLink(ctx context.Context, empty *emptypb.Empty) (*integrations_faceit.GetAuthLink, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Integrations) IntegrationsFaceitGetData(ctx context.Context, empty *emptypb.Empty) (*integrations_faceit.GetDataResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Integrations) IntegrationsFaceitPostCode(ctx context.Context, request *integrations_faceit.PostCodeRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Integrations) IntegrationsFaceitLogout(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
