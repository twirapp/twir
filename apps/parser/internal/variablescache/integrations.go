package variables_cache

import (
	model "tsuwari/models"
)

func (c *VariablesCacheService) GetEnabledIntegrations() *[]model.ChannelsIntegrations {
	c.locks.integrations.Lock()
	defer c.locks.integrations.Unlock()

	if c.cache.Integrations != nil {
		return c.cache.Integrations
	}

	result := &[]model.ChannelsIntegrations{}
	err := c.Services.Db.Where(`"channelId" = ? AND enabled = ?`, c.ChannelId, true).
		Preload("Integrations").
		Find(result).
		Error

	if err == nil {
		c.cache.Integrations = result
	}

	return c.cache.Integrations
}
