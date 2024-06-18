package duplicate_tracker

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Opts struct {
	Redis *redis.Client
}

func New(opts Opts) *DuplicateTracker {
	return &DuplicateTracker{redis: opts.Redis}
}

type DuplicateTracker struct {
	redis *redis.Client
}

// AddAndCheckIfDuplicate errors ignored for case of redis failure, and we not surethat we need return error to twitch
func (m *DuplicateTracker) AddAndCheckIfDuplicate(ctx context.Context, id string) (bool, error) {
	redisKey := fmt.Sprintf("eventsub:duplicates-track:%s", id)

	exists, _ := m.redis.Exists(ctx, redisKey).Result()
	if exists == 1 {
		return true, nil
	}

	m.redis.Set(ctx, redisKey, 1, 24*time.Hour).Err()

	return false, nil
}
