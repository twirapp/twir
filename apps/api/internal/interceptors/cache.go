package interceptors

import (
	"context"
	"fmt"
	"time"

	json "github.com/bytedance/sonic"
	"github.com/twitchtv/twirp"
	"go.uber.org/zap"
)

type CacheOpts struct {
	CacheMethod       string
	CacheDuration     time.Duration
	ClearMethods      []string
	WithChannelHeader bool
	CastTo            func() any
}

func (s *Service) NewCacheInterceptor(options ...CacheOpts) twirp.Interceptor {
	interceptors := make(map[string]CacheOpts)
	for _, option := range options {
		interceptors[option.CacheMethod] = option
		for _, clearMethod := range option.ClearMethods {
			interceptors[clearMethod] = option
		}
	}

	return func(next twirp.Method) twirp.Method {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			method, ok := twirp.MethodName(ctx)
			if !ok {
				return nil, twirp.InternalError("failed to get method name")
			}

			option, ok := interceptors[method]
			if !ok {
				return next(ctx, req)
			}

			channelId := ctx.Value("dashboardId").(string)
			cacheKey := fmt.Sprintf("api:cache:twirp-%s", option.CacheMethod)
			if option.WithChannelHeader {
				cacheKey += "-channel-" + channelId
			}

			// delete cache if returns with no error
			for _, clearMethod := range option.ClearMethods {
				if clearMethod == method {
					s.redis.Del(ctx, cacheKey)
					return next(ctx, req)
				}
			}

			if method != option.CacheMethod {
				return next(ctx, req)
			}

			cached, _ := s.redis.Get(ctx, cacheKey).Bytes()
			if cached != nil {
				castedData := option.CastTo()
				unmarshalErr := json.Unmarshal(cached, castedData)
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
					redisSetErr := s.redis.Set(ctx, cacheKey, bytes, option.CacheDuration).Err()
					if redisSetErr != nil {
						zap.S().Error(redisSetErr)
					}
				}
			}

			return result, err
		}
	}
}
