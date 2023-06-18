package interceptors

import (
	"context"
	"github.com/twitchtv/twirp"
)

func (s *Service) SessionInterceptor(next twirp.Method) twirp.Method {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		user := s.sessionManager.Get(ctx, "user")
		// TODO: check api key and attach user to context from api key
		if user == nil {
			return nil, twirp.Unauthenticated.Error("not authenticated")
		}
		return next(ctx, req)
	}
}
