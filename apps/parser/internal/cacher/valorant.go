package cacher

import (
	"context"
	"fmt"

	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
)

// GetValorantMatches implements types.VariablesCacher
func (c *cacher) GetValorantMatches(ctx context.Context) []types.ValorantMatch {
	c.locks.valorantMatches.Lock()
	defer c.locks.valorantMatches.Unlock()

	if c.cache.valorantMatches != nil {
		return c.cache.valorantMatches
	}

	integrations := c.GetEnabledChannelIntegrations(ctx)
	integration, ok := lo.Find(
		integrations, func(item *model.ChannelsIntegrations) bool {
			return item.Integration.Service == "VALORANT"
		},
	)

	if !ok || !integration.Enabled || integration.Data == nil || integration.Data.UserName == nil ||
		integration.Data.ValorantActiveRegion == nil || integration.Data.ValorantPuuid == nil {
		return nil
	}

	apiUrl := fmt.Sprintf(
		"https://api.henrikdev.xyz/valorant/v3/by-puuid/matches/%s/%s",
		*integration.Data.ValorantActiveRegion,
		*integration.Data.ValorantPuuid,
	)

	var data *types.ValorantMatchesResponse

	r := req.R().
		SetContext(ctx).
		SetHeader("Authorization", c.services.Config.ValorantHenrikApiKey).
		SetSuccessResult(&data)

	_, err := r.Get(apiUrl)
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return nil
	}

	c.cache.valorantMatches = data.Data

	return c.cache.valorantMatches
}

// GetValorantProfile implements types.VariablesCacher
func (c *cacher) GetValorantMMR(ctx context.Context) *types.ValorantMMR {
	c.locks.valorantProfile.Lock()
	defer c.locks.valorantProfile.Unlock()

	if c.cache.valorantProfile != nil {
		return c.cache.valorantProfile
	}

	integrations := c.GetEnabledChannelIntegrations(ctx)
	integration, ok := lo.Find(
		integrations, func(item *model.ChannelsIntegrations) bool {
			return item.Integration.Service == "VALORANT"
		},
	)

	if !ok || !integration.Enabled || integration.Data == nil || integration.Data.UserName == nil ||
		integration.Data.ValorantActiveRegion == nil || integration.Data.ValorantPuuid == nil {
		return nil
	}

	apiUrl := fmt.Sprintf(
		"https://api.henrikdev.xyz/valorant/v3/by-puuid/mmr/%s/pc/%s",
		*integration.Data.ValorantActiveRegion,
		*integration.Data.ValorantPuuid,
	)

	c.cache.valorantProfile = &types.ValorantMMR{}

	r := req.R().
		SetContext(ctx).
		SetHeader("Authorization", c.services.Config.ValorantHenrikApiKey).
		SetSuccessResult(c.cache.valorantProfile)

	_, err := r.Get(apiUrl)
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return nil
	}

	return c.cache.valorantProfile
}
