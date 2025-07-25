package manager

import (
	"context"

	cfg "github.com/twirapp/twir/libs/config"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"go.uber.org/fx"
)

type Creds struct {
	cfg     cfg.Config
	twirBus *buscore.Bus
}

type CredsOpts struct {
	fx.In

	Config  cfg.Config
	TwirBus *buscore.Bus
}

func NewCreds(opts CredsOpts) *Creds {
	return &Creds{
		cfg:     opts.Config,
		twirBus: opts.TwirBus,
	}
}

func (c *Creds) ClientID() (string, error) {
	return c.cfg.TwitchClientId, nil
}
func (c *Creds) AppToken() (string, error) {
	appToken, err := c.twirBus.Tokens.RequestAppToken.Request(context.TODO(), struct{}{})
	if err != nil {
		return "", err
	}

	return appToken.Data.AccessToken, nil
}
