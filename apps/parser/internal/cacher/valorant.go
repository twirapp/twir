package cacher

import (
	"context"
	"fmt"

	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
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
		SetSuccessResult(&data)

	if c.services.Config.ValorantHenrikApiKey != "" {
		r.SetHeader("Authorization", c.services.Config.ValorantHenrikApiKey)
	}

	_, err := r.Get(apiUrl)
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return nil
	}

	c.cache.valorantMatches = data.Data

	return c.cache.valorantMatches
}

// GetValorantProfile implements types.VariablesCacher
func (c *cacher) GetValorantProfile(ctx context.Context) *types.ValorantProfile {
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
		"https://api.henrikdev.xyz/valorant/v2/by-puuid/mmr/%s/%s",
		*integration.Data.ValorantActiveRegion,
		*integration.Data.ValorantPuuid,
	)

	c.cache.valorantProfile = &types.ValorantProfile{}

	r := req.R().
		SetContext(ctx).
		SetSuccessResult(c.cache.valorantProfile)

	if c.services.Config.ValorantHenrikApiKey != "" {
		r.SetHeader("Authorization", c.services.Config.ValorantHenrikApiKey)
	}

	_, err := r.Get(apiUrl)
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return nil
	}

	return c.cache.valorantProfile
}
