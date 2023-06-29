package interceptors

import (
	"context"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twitchtv/twirp"
	"net/http"
)

func (s *Service) ChannelAccessInterceptor(next twirp.Method) twirp.Method {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		user := ctx.Value("dbUser")
		if user == nil {
			return nil, twirp.Internal.Error("internal error")
		}
		castedUser, ok := user.(*model.Users)
		if !ok {
			return nil, twirp.Internal.Error("internal error")
		}

		if castedUser.IsBotAdmin {
			return next(ctx, req)
		}

		dashboardId := ctx.Value("dashboardId")
		if dashboardId == nil {
			return nil, twirp.NewError(twirp.ErrorCode(http.StatusBadRequest), "no dashboardId provided")
		}
		castedDashboardId, ok := dashboardId.(string)
		if !ok {
			return nil, twirp.NewError(twirp.ErrorCode(http.StatusBadRequest), "wrong type of dashboardId")
		}

		if castedUser.ID == castedDashboardId {
			return next(ctx, req)
		}

		_, ok = lo.Find(castedUser.Roles, func(a model.ChannelRoleUser) bool {
			return a.Role.ChannelID == castedDashboardId &&
				lo.Contains(a.Role.Permissions, model.RolePermissionCanAccessDashboard.String())
		})

		if ok {
			return next(ctx, req)
		}

		return nil, twirp.NewError(twirp.ErrorCode(http.StatusForbidden), "not authorized to that dashboard")
	}
}
