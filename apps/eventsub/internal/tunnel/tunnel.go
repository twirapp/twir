package tunnel

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/samber/lo"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"go.uber.org/fx"
	"golang.ngrok.com/ngrok"
	ngconfig "golang.ngrok.com/ngrok/config"
)

type AppTunnel struct {
	net.Listener

	cfg config.Config
}

func New(cfg config.Config, lc fx.Lifecycle, log logger.Logger) (*AppTunnel, error) {
	tunn := &AppTunnel{
		cfg: cfg,
	}

	if cfg.AppEnv != "production" {
		if cfg.NgrokAuthToken == "" {
			panic("NGROK_AUTH_TOKEN is required")
		}

		lis, err := retry.DoWithData(
			func() (ngrok.Tunnel, error) {
				return ngrok.Listen(
					context.Background(),
					ngconfig.HTTPEndpoint(),
					ngrok.WithAuthtoken(cfg.NgrokAuthToken),
				)
			},
			retry.Attempts(5),
			retry.Delay(1*time.Second),
		)
		if err != nil {
			return nil, err
		}

		tunn.Listener = lis
	} else {
		lis, err := createDefaultTun()
		if err != nil {
			return nil, err
		}
		tunn.Listener = lis
	}

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				log.Info("Tunnel runned", slog.String("addr", tunn.GetAddr()))
				return nil
			},
		},
	)

	return tunn, nil
}

func (c *AppTunnel) GetAddr() string {
	return lo.
		If(
			c.cfg.AppEnv != "production",
			"https://"+c.Listener.Addr().String(),
		).
		Else(fmt.Sprintf("https://eventsub.%s", c.cfg.SiteBaseUrl))
}

func createDefaultTun() (net.Listener, error) {
	return net.Listen("tcp", ":3003")
}
