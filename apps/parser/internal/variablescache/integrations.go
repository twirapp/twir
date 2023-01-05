package variables_cache

import (
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	model "github.com/satont/tsuwari/libs/gomodels"
	"gorm.io/gorm"
)

func (c *VariablesCacheService) GetEnabledIntegrations() []model.ChannelsIntegrations {
	c.locks.integrations.Lock()
	defer c.locks.integrations.Unlock()

	if c.cache.Integrations != nil {
		return c.cache.Integrations
	}

	db := do.MustInvoke[gorm.DB](di.Provider)

	result := []model.ChannelsIntegrations{}
	err := db.Where(`"channelId" = ? AND enabled = ?`, c.ChannelId, true).
		Preload("Integration").
		Find(&result).
		Error

	if err == nil {
		c.cache.Integrations = result
	}

	return c.cache.Integrations
}
