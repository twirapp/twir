package helpers

import (
	"context"
	"net"
	"time"

	config "github.com/satont/twir/libs/config"
	"golang.ngrok.com/ngrok"
	ngrok_config "golang.ngrok.com/ngrok/config"
)

func GetAppTunnel(ctx context.Context, cfg *config.Config) (net.Listener, error) {
	if cfg.AppEnv != "production" {
		ngrokCtx, cancelNgrokCtx := context.WithTimeout(ctx, 5*time.Second)
		tun, err := ngrok.Listen(
			ngrokCtx,
			ngrok_config.HTTPEndpoint(),
		)
		if err != nil {
			cancelNgrokCtx()
			return createDefaultTun()
		}

		return tun, nil
	} else {
		return createDefaultTun()
	}
}

func createDefaultTun() (net.Listener, error) {
	return net.Listen("tcp", ":3003")
}
