package interceptors

import (
	"context"
	"fmt"
	json "github.com/bytedance/sonic"
	"github.com/twitchtv/twirp"
	"go.uber.org/zap"
	"time"
)

type CacheOpts struct {
	CacheMethod       string
	CacheDuration     time.Duration
	ClearMethods      []string
	WithChannelHeader bool
	NewCastTo         func() any
}

func (s *Service) NewCacheInterceptor(options ...CacheOpts) twirp.Interceptor {
	return func(next twirp.Method) twirp.Method {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			method, ok := twirp.MethodName(ctx)
			if !ok {
				return nil, twirp.InternalError("failed to get method name")
			}

			var option *CacheOpts
			for _, opt := range options {
				if opt.CacheMethod == method {
					option = &opt
					break
				}

				for _, clearMethod := range opt.ClearMethods {
					if clearMethod == method {
						option = &opt
						break
					}
				}
			}

			if option == nil {
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
					res, err := next(ctx, req)
					if err == nil {
						s.redis.Del(ctx, cacheKey)
					}
					return res, err
				}
			}

			if method != option.CacheMethod {
				return next(ctx, req)
			}

			cached, _ := s.redis.Get(ctx, cacheKey).Bytes()
			if cached != nil {
				castedData := option.NewCastTo()
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
