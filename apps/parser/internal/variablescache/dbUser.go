package variables_cache

import (
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	model "github.com/satont/tsuwari/libs/gomodels"
	"gorm.io/gorm"
)

func (c *VariablesCacheService) GetGbUser() *model.UsersStats {
	c.locks.dbUser.Lock()
	defer c.locks.dbUser.Unlock()

	if c.cache.DbUserStats != nil {
		return c.cache.DbUserStats
	}

	db := do.MustInvoke[gorm.DB](di.Provider)

	result := &model.UsersStats{}

	err := db.
		Where(`"userId" = ? AND "channelId" = ?`, c.SenderId, c.ChannelId).
		Find(result).
		Error
	if err == nil {
		c.cache.DbUserStats = result
	}

	return c.cache.DbUserStats
}
