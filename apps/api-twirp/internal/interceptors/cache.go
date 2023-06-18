package interceptors

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/twitchtv/twirp"
	"go.uber.org/zap"
	"time"
)

type CacheOpts struct {
	CacheMethod       string
	CacheDuration     time.Duration
	WithChannelHeader bool
	NewCastTo         func() any
}

func (s *Service) NewCacheInterceptor(opts CacheOpts) twirp.Interceptor {
	return func(next twirp.Method) twirp.Method {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			method, ok := twirp.MethodName(ctx)
			if method != opts.CacheMethod || !ok {
				return next(ctx, req)
			}

			channelId := ctx.Value("dashboard_id").(string)

			cacheKey := fmt.Sprintf("api-cache-twirp-%s", opts.CacheMethod)
			if opts.WithChannelHeader {
				cacheKey += "-channel-" + channelId
			}

			cached, _ := s.redis.Get(ctx, cacheKey).Result()
			if cached != "" {
				castedData := opts.NewCastTo()
				unmarshalErr := json.Unmarshal([]byte(cached), castedData)
				if unmarshalErr != nil {
					zap.S().Error(unmarshalErr)
					return nil, unmarshalErr
				}

				return castedData, nil
			}

			result, err := next(ctx, req)

			if err == nil {
				bytes, marshallErr := json.Marshal(result)
				if marshallErr != nil {
					zap.S().Error(marshallErr)
				} else {
					redisSetErr := s.redis.Set(ctx, cacheKey, bytes, opts.CacheDuration).Err()
					if redisSetErr != nil {
						zap.S().Error(redisSetErr)
					}
				}
			}

			return result, err
		}
	}
}
