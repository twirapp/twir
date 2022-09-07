package variables_cache

import model "tsuwari/parser/internal/models"

func (c *VariablesCacheService) GetEnabledIntegrations() *[]model.ChannelInegrationWithRelation {
	c.locks.integrations.Lock()
	defer c.locks.integrations.Unlock()

	if c.cache.Integrations != nil {
		return c.cache.Integrations
	}

	result := &[]model.ChannelInegrationWithRelation{}
	err := c.Services.Db.Where(`"channelId" = ? AND enabled = ?`, c.ChannelId, true).Joins("Integration").Find(result).Error

	if err == nil {
		c.cache.Integrations = result
	}

	return c.cache.Integrations
}
