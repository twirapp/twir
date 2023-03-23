package grpc_impl

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/emotes-cacher/internal/di"
	"github.com/satont/tsuwari/apps/emotes-cacher/internal/emotes"
	"github.com/satont/tsuwari/libs/grpc/generated/emotes_cacher"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
	"time"
)

type EmotesCacherImpl struct {
	emotes_cacher.UnimplementedEmotesCacherServer

	redis  redis.Client
	logger zap.Logger
}

func NewEmotesCacher() *EmotesCacherImpl {
	redisClient := do.MustInvoke[redis.Client](di.Provider)
	logger := do.MustInvoke[zap.Logger](di.Provider)

	return &EmotesCacherImpl{
		redis:  redisClient,
		logger: logger,
	}
}

func (c *EmotesCacherImpl) CacheChannelEmotes(_ context.Context, req *emotes_cacher.Request) (*emptypb.Empty, error) {
	if req.ChannelId == "" {
		return &emptypb.Empty{}, nil
	}

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	resultEmotes := make([]string, 300)

	wg.Add(3)

	go func() {
		defer wg.Done()
		em, err := emotes.GetChannelSevenTvEmotes(req.ChannelId)
		if err != nil || em == nil || len(em) == 0 {
			return
		}

		mu.Lock()
		defer mu.Unlock()
		resultEmotes = append(resultEmotes, em...)
	}()

	go func() {
		defer wg.Done()
		em, err := emotes.GetChannelBttvEmotes(req.ChannelId)
		if err != nil || em == nil || len(em) == 0 {
			return
		}

		mu.Lock()
		defer mu.Unlock()
		resultEmotes = append(resultEmotes, em...)
	}()

	go func() {
		defer wg.Done()
		em, err := emotes.GetChannelFfzEmotes(req.ChannelId)
		if err != nil || em == nil || len(em) == 0 {
			return
		}

		mu.Lock()
		defer mu.Unlock()
		resultEmotes = append(resultEmotes, em...)
	}()

	wg.Wait()

	for _, emote := range resultEmotes {
		if emote == "" {
			continue
		}
		go func(emote string) {
			c.redis.Set(
				context.Background(),
				fmt.Sprintf("emotes:channel:%s:%s", req.ChannelId, emote),
				emote,
				10*time.Minute,
			)
		}(emote)
	}

	return &emptypb.Empty{}, nil
}

func (c *EmotesCacherImpl) CacheGlobalEmotes(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	resultEmotes := make([]string, 300)

	wg.Add(3)

	go func() {
		defer wg.Done()
		em, err := emotes.GetGlobalSevenTvEmotes()
		if err != nil || em == nil || len(em) == 0 {
			return
		}

		mu.Lock()
		defer mu.Unlock()
		resultEmotes = append(resultEmotes, em...)
	}()

	go func() {
		defer wg.Done()
		em, err := emotes.GetGlobalFfzEmotes()
		if err != nil || em == nil || len(em) == 0 {
			return
		}

		mu.Lock()
		defer mu.Unlock()
		resultEmotes = append(resultEmotes, em...)
	}()

	go func() {
		defer wg.Done()
		em, err := emotes.GetGlobalBttvEmotes()
		if err != nil || em == nil || len(em) == 0 {
			return
		}

		mu.Lock()
		defer mu.Unlock()
		resultEmotes = append(resultEmotes, em...)
	}()

	wg.Wait()

	for _, emote := range resultEmotes {
		if emote == "" {
			continue
		}
		go func(emote string) {
			c.redis.Set(
				context.Background(),
				fmt.Sprintf("emotes:global:%s", emote),
				emote,
				10*time.Minute,
			)
		}(emote)
	}

	return &emptypb.Empty{}, nil
}
