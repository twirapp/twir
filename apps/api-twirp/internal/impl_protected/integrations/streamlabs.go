package integrations

import (
	"context"
	"github.com/satont/tsuwari/libs/grpc/generated/api/integrations_streamlabs"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Integrations) IntegrationsStreamlabsGetAuthLink(ctx context.Context, empty *emptypb.Empty) (*integrations_streamlabs.GetAuthLink, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Integrations) IntegrationsStreamlabsGetData(ctx context.Context, empty *emptypb.Empty) (*integrations_streamlabs.GetDataResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Integrations) IntegrationsStreamlabsPostCode(ctx context.Context, request *integrations_streamlabs.PostCodeRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Integrations) IntegrationsStreamlabsLogout(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
