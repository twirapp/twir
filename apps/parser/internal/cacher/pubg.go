package cacher

import (
	"context"
	"errors"

	"github.com/NovikovRoman/pubg"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	tpubg "github.com/twirapp/twir/libs/pubg"
)

func (c *cacher) GetPubgCurrentSeason(ctx context.Context) (*string, error) {
	c.locks.pubgCurrentSeason.Lock()
	defer c.locks.pubgCurrentSeason.Unlock()

	if c.cache.pubgCurrentSeason != nil {
		return c.cache.pubgCurrentSeason, nil
	}

	season, err := c.services.PubgClient.GetCurrentSeason(ctx)
	if err != nil {
		return nil, err
	}

	c.cache.pubgCurrentSeason = season
	return c.cache.pubgCurrentSeason, nil

}

func (c *cacher) GetPubgLifetimeData(ctx context.Context) (*pubg.LifetimeStatsPlayer, error) {
	c.locks.pubgLifetimeData.Lock()
	defer c.locks.pubgLifetimeData.Unlock()

	if c.cache.pubgLifetimeData != nil {
		return c.cache.pubgLifetimeData, nil
	}

	c.cache.pubgLifetimeData = &pubg.LifetimeStatsPlayer{}

	if c.cache.pubgCurrentSeason == nil {
		_, err := c.GetPubgCurrentSeason(ctx)
		if err != nil {
			c.services.Logger.Sugar().Error(err)
			return nil, err
		}
	}

	integrations := c.GetEnabledChannelIntegrations(ctx)

	if integrations == nil {
		return nil, errors.New("no enabled integrations")
	}

	integration, ok := lo.Find(
		integrations, func(i *model.ChannelsIntegrations) bool {
			return i.Integration.Service == "PUBG" && i.Enabled
		},
	)
	if !ok {
		return nil, errors.New("pubg integration not enabled")
	}

	userLifetimeStats, err := c.services.PubgClient.GetLifetimeStats(
		ctx,
		lo.FromPtr(integration.Data.UserId),
	)
	if err != nil {
		if errors.Is(err, tpubg.ErrOverloaded) {
			return nil, tpubg.ErrOverloaded
		}
		if errors.Is(err, tpubg.ErrPubg) {
			return nil, tpubg.ErrPubg
		}

		c.services.Logger.Sugar().Error(err)
		return nil, errors.Unwrap(err)
	}
	c.cache.pubgLifetimeData = userLifetimeStats

	return c.cache.pubgLifetimeData, nil
}
