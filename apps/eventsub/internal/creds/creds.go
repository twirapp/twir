package creds

import (
	"context"

	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Creds struct {
	appCtx context.Context

	clientId   string
	tokensGrpc tokens.TokensClient
}

func NewCreds(ctx context.Context, cfg *config.Config, tokensGrpc tokens.TokensClient) *Creds {
	return &Creds{
		appCtx:     ctx,
		clientId:   cfg.TwitchClientId,
		tokensGrpc: tokensGrpc,
	}
}

func (c *Creds) ClientID() (string, error) {
	return c.clientId, nil
}
func (c *Creds) AppToken() (string, error) {
	appToken, err := c.tokensGrpc.RequestAppToken(c.appCtx, &emptypb.Empty{})
	if err != nil {
		return "", err
	}

	return appToken.AccessToken, nil
}
