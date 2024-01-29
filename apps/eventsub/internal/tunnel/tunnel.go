package tunnel

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"regexp"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/imroc/req/v3"
	"github.com/localtunnel/go-localtunnel"
	"github.com/samber/lo"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"go.uber.org/fx"
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
		lis, err := retry.DoWithData(
			func() (*localtunnel.Listener, error) {
				return localtunnel.Listen(
					localtunnel.Options{},
				)
			},
			retry.Attempts(5),
			retry.Delay(1*time.Second),
		)
		if err != nil {
			return nil, err
		}

		if err := validateLocalTunnel(lis); err != nil {
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
			c.Listener.Addr().String(),
		).
		Else(fmt.Sprintf("https://eventsub.%s", c.cfg.SiteBaseUrl))
}

func createDefaultTun() (net.Listener, error) {
	return net.Listen("tcp", ":3003")
}

var validateRgx = regexp.MustCompile(`url: "(.+)"`)

const agent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

func validateLocalTunnel(listener net.Listener) error {
	addr := listener.Addr().String()
	reqClient := req.C().SetUserAgent(agent).SetTimeout(5 * time.Second)

	res, err := reqClient.R().Get(addr)
	if err != nil {
		return err
	}
	if !res.IsSuccessState() && res.StatusCode != 511 {
		return fmt.Errorf("failed to get localtunnel: %s", res.String())
	}

	matches := validateRgx.FindStringSubmatch(res.String())
	if len(matches) < 2 {
		return err
	}

	validateUrl := addr + matches[1]

	res, err = reqClient.R().Get("https://loca.lt/mytunnelpassword")
	if err != nil {
		return err
	}
	if !res.IsSuccessState() {
		return fmt.Errorf("failed to get localtunnel password: %s", res.String())
	}

	ip := res.String()

	res, err = reqClient.R().SetFormData(
		map[string]string{
			"endpoint": ip,
		},
	).Post(validateUrl)
	if err != nil {
		return err
	}
	if !res.IsSuccessState() {
		return fmt.Errorf("failed to validate localtunnel: %s", res.String())
	}

	return nil
}
