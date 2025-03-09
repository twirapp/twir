package streamelements

import (
	"context"
	"fmt"
	"net/url"

	config "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/libs/integrations/streamelements"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
)

type Opts struct {
	fx.In

	Config config.Config
}

func New(opts Opts) (*Service, error) {
	s := &Service{
		config: opts.Config,
	}

	siteBaseUrl, err := url.Parse(opts.Config.SiteBaseUrl)
	if err != nil {
		return nil, err
	}

	s.redirectURL = siteBaseUrl.JoinPath("/dashboard/integrations/streamelements").String()

	return s, nil
}

type Service struct {
	config      config.Config
	redirectURL string
}

func (c *Service) GetAuthLink() (string, error) {
	if c.config.StreamElementsClientId == "" || c.config.StreamElementsClientSecret == "" {
		return "", fmt.Errorf("service not configured")
	}

	i := streamelements.New(
		c.config.StreamElementsClientId,
		c.config.StreamElementsClientSecret,
	)

	return i.GetAuthLink(c.redirectURL), nil
}

func (c *Service) ExchangeDataByCode(
	ctx context.Context,
	code string,
) (*Data, error) {
	if c.config.StreamElementsClientId == "" || c.config.StreamElementsClientSecret == "" {
		return nil, nil
	}

	i := streamelements.New(
		c.config.StreamElementsClientId,
		c.config.StreamElementsClientSecret,
	)

	_, err := i.ExchangeCode(ctx, code, c.redirectURL)
	if err != nil {
		return nil, err
	}

	profile, err := i.GetProfile(ctx)
	if err != nil {
		return nil, err
	}

	data := &Data{
		Commands: nil,
		Timers:   nil,
	}
	var errgp errgroup.Group

	errgp.Go(
		func() error {
			cmds, err := i.GetCommands(ctx, profile.ID)
			if err != nil {
				return err
			}

			data.Commands = cmds
			return nil
		},
	)

	errgp.Go(
		func() error {
			timers, err := i.GetTimers(ctx, profile.ID)
			if err != nil {
				return err
			}

			data.Timers = timers
			return nil
		},
	)

	if err := errgp.Wait(); err != nil {
		return nil, fmt.Errorf("failed to exchange data: %w", err)
	}

	return data, nil
}
