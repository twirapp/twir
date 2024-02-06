package manager

import (
	"context"

	cfg "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.uber.org/fx"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Creds struct {
	cfg        cfg.Config
	tokensGrpc tokens.TokensClient
}

type CredsOpts struct {
	fx.In

	Config     cfg.Config
	TokensGrpc tokens.TokensClient
}

func NewCreds(opts CredsOpts) *Creds {
	return &Creds{
		cfg:        opts.Config,
		tokensGrpc: opts.TokensGrpc,
	}
}

func (c *Creds) ClientID() (string, error) {
	return c.cfg.TwitchClientId, nil
}
func (c *Creds) AppToken() (string, error) {
	appToken, err := c.tokensGrpc.RequestAppToken(context.TODO(), &emptypb.Empty{})
	if err != nil {
		return "", err
	}

	return appToken.AccessToken, nil
}
