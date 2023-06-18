package integrations

import (
	"context"
	"github.com/satont/tsuwari/libs/grpc/generated/api/integrations_vk"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Integrations) IntegrationsVKGetAuthLink(ctx context.Context, empty *emptypb.Empty) (*integrations_vk.GetAuthLink, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Integrations) IntegrationsVKGetData(ctx context.Context, empty *emptypb.Empty) (*integrations_vk.GetDataResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Integrations) IntegrationsVKPostCode(ctx context.Context, request *integrations_vk.PostCodeRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Integrations) IntegrationsVKLogout(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
