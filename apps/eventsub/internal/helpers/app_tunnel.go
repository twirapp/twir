package helpers

import (
	"context"
	"net"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/localtunnel/go-localtunnel"
	config "github.com/satont/twir/libs/config"
)

func GetAppTunnel(ctx context.Context, cfg *config.Config) (net.Listener, error) {
	if cfg.AppEnv != "production" {
		return retry.DoWithData(
			func() (*localtunnel.Listener, error) {
				return localtunnel.Listen(
					localtunnel.Options{},
				)
			},
			retry.Attempts(5),
			retry.Delay(1*time.Second),
		)
	} else {
		return createDefaultTun()
	}
}

func createDefaultTun() (net.Listener, error) {
	return net.Listen("tcp", ":3003")
}
