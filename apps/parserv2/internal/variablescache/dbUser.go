package variables_cache

import model "tsuwari/parser/internal/models"

func (c *VariablesCacheService) GetGbUser() *model.UsersStats {
	c.locks.dbUser.Lock()
	defer c.locks.dbUser.Unlock()

	if c.cache.DbUserStats != nil {
		return c.cache.DbUserStats
	}

	result := model.UsersStats{}
	err := c.Services.Db.Where(`"userId" = ? AND "channelId" = ?`, c.SenderId, c.ChannelId).Find(&result).Error
	if err == nil {
		c.cache.DbUserStats = &result
	}

	return c.cache.DbUserStats
}
