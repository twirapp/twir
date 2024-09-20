package interceptors

import (
	"context"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/twitchtv/twirp"
)

func (s *Service) DashboardId(next twirp.Method) twirp.Method {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		dashboardId := s.sessionManager.Get(ctx, "dashboardId")
		dbUser := ctx.Value("dbUser")

		if dashboardId == nil {
			ctx = context.WithValue(ctx, "dashboardId", dbUser.(*model.Users).ID)
		} else {
			ctx = context.WithValue(ctx, "dashboardId", dashboardId)
		}

		return next(ctx, req)
	}
}
