package cacher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiUrl, nil)
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return nil
	}
	req.Header.Set("Authorization", c.services.Config.ValorantHenrikApiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return nil
	}

	var data types.ValorantMatchesResponse
	if err := json.Unmarshal(body, &data); err != nil {
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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiUrl, nil)
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return nil
	}
	req.Header.Set("Authorization", c.services.Config.ValorantHenrikApiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return nil
	}

	c.cache.valorantProfile = &types.ValorantMMR{}
	if err := json.Unmarshal(body, c.cache.valorantProfile); err != nil {
		c.services.Logger.Sugar().Error(err)
		return nil
	}

	return c.cache.valorantProfile
}
