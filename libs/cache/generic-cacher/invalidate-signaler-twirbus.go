package generic_cacher

import (
	"context"

	buscore "github.com/twirapp/twir/libs/bus-core"
	cache_invalidator "github.com/twirapp/twir/libs/bus-core/cache-invalidator"
)

var _ InvalidateSignaler = (*BusCoreInvalidator)(nil)

func NewBusCoreInvalidator(bus *buscore.Bus) *BusCoreInvalidator {
	return &BusCoreInvalidator{bus: bus}
}

type BusCoreInvalidator struct {
	bus *buscore.Bus
}

func (i *BusCoreInvalidator) Send(key string) error {
	err := i.bus.CacheInvalidator.Publish(
		context.Background(),
		cache_invalidator.InvalidateRequest{
			Key: key,
		},
	)

	return err
}

func (i *BusCoreInvalidator) Receiver() <-chan string {
	ch := make(chan string)
	i.bus.CacheInvalidator.Subscribe(
		func(ctx context.Context, data cache_invalidator.InvalidateRequest) (struct{}, error) {
			ch <- data.Key

			return struct{}{}, nil
		},
	)

	return ch
}
