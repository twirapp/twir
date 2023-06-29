package helpers

import (
	"context"
	config "github.com/satont/twir/libs/config"
	"golang.ngrok.com/ngrok"
	ngrok_config "golang.ngrok.com/ngrok/config"
	"net"
)

func GetAppTunnel(ctx context.Context, cfg *config.Config) (net.Listener, error) {
	if cfg.AppEnv != "production" {
		tun, err := ngrok.Listen(
			ctx,
			ngrok_config.HTTPEndpoint(),
		)
		if err != nil {
			return nil, err
		}

		return tun, nil
	} else {
		tun, err := net.Listen("tcp", ":3003")
		if err != nil {
			return nil, err
		}
		return tun, nil
	}
}
