package integrations

import (
	"context"
	"github.com/satont/tsuwari/libs/grpc/generated/api/integrations_spotify"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Integrations) IntegrationsSpotifyGetAuthLink(ctx context.Context, empty *emptypb.Empty) (*integrations_spotify.GetAuthLink, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Integrations) IntegrationsSpotifyGetData(ctx context.Context, empty *emptypb.Empty) (*integrations_spotify.GetDataResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Integrations) IntegrationsSpotifyPostCode(ctx context.Context, request *integrations_spotify.PostCodeRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Integrations) IntegrationsSpotifyLogout(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
