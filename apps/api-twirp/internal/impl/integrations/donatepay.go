package integrations

import (
	"context"
	"github.com/satont/tsuwari/libs/grpc/generated/api/integrations_donatepay"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Integrations) IntegrationsDonatepayGet(ctx context.Context, empty *emptypb.Empty) (*integrations_donatepay.GetResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Integrations) IntegrationsDonatepayPut(ctx context.Context, request *integrations_donatepay.PostRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
