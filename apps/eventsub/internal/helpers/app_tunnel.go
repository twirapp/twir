package helpers

import (
	"context"
	"net"

	"github.com/localtunnel/go-localtunnel"
	config "github.com/satont/twir/libs/config"
)

func GetAppTunnel(ctx context.Context, cfg *config.Config) (net.Listener, error) {
	if cfg.AppEnv != "production" {
		tun, err := localtunnel.Listen(
			localtunnel.Options{},
		)
		if err != nil {
			return nil, err
		}

		return tun, nil
	} else {
		return createDefaultTun()
	}
}

func createDefaultTun() (net.Listener, error) {
	return net.Listen("tcp", ":3003")
}
