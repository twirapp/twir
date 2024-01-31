package cacher

import (
	"context"
	"fmt"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
)

// GetValorantMatches implements types.VariablesCacher
func (c *cacher) GetValorantMatches(ctx context.Context) []*types.ValorantMatch {
	c.locks.valorantMatch.Lock()
	defer c.locks.valorantMatch.Unlock()

	if c.cache.valorantMatches != nil {
		return c.cache.valorantMatches
	}

	var data *types.ValorantMatchesResponse

	integrations := c.GetEnabledChannelIntegrations(ctx)
	integration, ok := lo.Find(
		integrations, func(item *model.ChannelsIntegrations) bool {
			return item.Integration.Service == "VALORANT"
		},
	)

	if !ok || integration.Data == nil || integration.Data.UserName == nil ||
		integration.Data.ValorantActiveRegion == nil {
		return nil
	}

	apiUrl := fmt.Sprintf(
		"https://api.henrikdev.xyz/valorant/v3/matches/%s/",
		*integration.Data.ValorantActiveRegion,
	)

	_, err := req.R().
		SetContext(ctx).
		SetSuccessResult(&data).
		Get(apiUrl + strings.Replace(*integration.Data.UserName, "#", "/", 1))
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

	if !ok || integration.Data == nil || integration.Data.UserName == nil ||
		integration.Data.ValorantActiveRegion == nil {
		return nil
	}

	apiUrl := fmt.Sprintf(
		"https://api.henrikdev.xyz/valorant/v3/matches/%s/",
		*integration.Data.ValorantActiveRegion,
	)

	c.cache.valorantProfile = &types.ValorantProfile{}
	_, err := req.R().
		SetContext(ctx).
		SetSuccessResult(c.cache.valorantProfile).
		Get(apiUrl + strings.Replace(*integration.Data.UserName, "#", "/", 1))
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return nil
	}

	return c.cache.valorantProfile
}
