package cacher

import (
	"context"

	"github.com/satont/twir/libs/twitch"
)

func (c *cacher) GetSubAgeInfo(ctx context.Context, channelName, userName string) (
	*twitch.UserSubscribePayload,
	error,
) {
	c.locks.subage.Lock()
	defer c.locks.subage.Unlock()

	if c.cache.cachedSubAgeInfo != nil {
		return c.cache.cachedSubAgeInfo, nil
	}

	info, err := twitch.GetUserSubAge(ctx, channelName, userName)
	if err != nil {
		return nil, err
	}

	c.cache.cachedSubAgeInfo = info

	return info, nil
}
