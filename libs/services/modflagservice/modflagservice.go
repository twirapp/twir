package modflagservice

import (
	"context"
	"time"

	"github.com/twirapp/kv"
	kvoptions "github.com/twirapp/kv/options"
)

func New(store kv.KV) *ModFlagService {
	return &ModFlagService{store: store}
}

type ModFlagService struct {
	store kv.KV
}

func buildKey(channelId, moderatorUserId string) string {
	return "modFlag:" + channelId + ":" + moderatorUserId
}

func (c *ModFlagService) FlagModeratorToMod(
	ctx context.Context,
	channelId,
	moderatorUserId string,
	ttl time.Duration,
	value bool,
) error {
	return c.store.Set(ctx, buildKey(channelId, moderatorUserId), value, kvoptions.WithExpire(ttl))
}

func (c *ModFlagService) IsModeratorFlagged(
	ctx context.Context,
	channelId, moderatorUserId string,
) (bool, error) {
	return c.store.Get(ctx, buildKey(channelId, moderatorUserId)).Bool()
}
