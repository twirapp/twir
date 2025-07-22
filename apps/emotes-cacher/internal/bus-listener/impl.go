package bus_listener

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	bttvfetcher "github.com/satont/twir/apps/emotes-cacher/internal/services/bttv/fetcher"
	"github.com/satont/twir/apps/emotes-cacher/internal/services/ffz/fetcher"
	seventvfetcher "github.com/satont/twir/apps/emotes-cacher/internal/services/seventv/fetcher"
	"github.com/satont/twir/libs/logger"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"go.uber.org/fx"
)

type BusListener struct {
	redis  *redis.Client
	logger logger.Logger
	bus    *buscore.Bus
}

type Opts struct {
	fx.In
	Lc fx.Lifecycle

	Redis  *redis.Client
	Logger logger.Logger
	Bus    *buscore.Bus
}

func New(opts Opts) {
	impl := &BusListener{
		redis:  opts.Redis,
		logger: opts.Logger,
		bus:    opts.Bus,
	}

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := impl.bus.EmotesCacher.CacheGlobalEmotes.SubscribeGroup(
					"emotes-cacher",
					impl.cacheGlobalEmotes,
				); err != nil {
					return err
				}
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return nil
			},
		},
	)
}

func (c *BusListener) cacheGlobalEmotes(ctx context.Context, _ struct{}) (struct{}, error) {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	resultEmotes := make([]string, 300)

	wg.Add(3)

	go func() {
		defer wg.Done()
		em, err := seventvfetcher.GetGlobalSevenTvEmotes(ctx)
		if err != nil || em == nil || len(em) == 0 {
			return
		}

		mu.Lock()
		defer mu.Unlock()

		emotesNames := make([]string, 0, len(em))
		for _, e := range em {
			if e.Name == "" {
				continue
			}
			emotesNames = append(emotesNames, e.Name)
		}

		resultEmotes = append(resultEmotes, emotesNames...)
	}()

	go func() {
		defer wg.Done()
		em, err := fetcher.GetGlobalFfzEmotes(ctx)
		if err != nil || em == nil || len(em) == 0 {
			return
		}

		mu.Lock()
		defer mu.Unlock()
		resultEmotes = append(resultEmotes, em...)
	}()

	go func() {
		defer wg.Done()
		em, err := bttvfetcher.GetGlobalBttvEmotes(ctx)
		if err != nil || em == nil || len(em) == 0 {
			return
		}

		mu.Lock()
		defer mu.Unlock()
		resultEmotes = append(resultEmotes, em...)
	}()

	wg.Wait()

	c.redis.Pipelined(
		ctx, func(pipe redis.Pipeliner) error {
			for _, emote := range resultEmotes {
				if emote == "" {
					continue
				}

				pipe.Set(
					context.Background(),
					fmt.Sprintf("emotes:global:%s", emote),
					emote,
					10*time.Minute,
				)
			}

			return nil
		},
	)

	return struct{}{}, nil
}
