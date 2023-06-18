package integrations

import (
	"context"
	"github.com/satont/tsuwari/libs/grpc/generated/api/integrations_lastfm"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Integrations) IntegrationsLastFMGetAuthLink(ctx context.Context, empty *emptypb.Empty) (*integrations_lastfm.GetAuthLink, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Integrations) IntegrationsLastFMGetData(ctx context.Context, empty *emptypb.Empty) (*integrations_lastfm.GetDataResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Integrations) IntegrationsLastFMPostCode(ctx context.Context, request *integrations_lastfm.PostCodeRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Integrations) IntegrationsLastFMLogout(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
