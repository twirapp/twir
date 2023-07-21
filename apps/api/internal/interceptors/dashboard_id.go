package interceptors

import (
	"context"
	"github.com/twitchtv/twirp"
)

func (s *Service) DashboardId(next twirp.Method) twirp.Method {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		dashboardId := s.sessionManager.Get(ctx, "dashboardId")
		ctx = context.WithValue(ctx, "dashboardId", dashboardId)

		return next(ctx, req)
	}
}
