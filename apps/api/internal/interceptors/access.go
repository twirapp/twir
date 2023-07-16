package interceptors

import (
	"context"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twitchtv/twirp"
)

func (s *Service) ChannelAccessInterceptor(next twirp.Method) twirp.Method {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		method, _ := twirp.MethodName(ctx)
		if method == "AuthUserProfile" || method == "AuthGetDashboards" {
			return next(ctx, req)
		}

		user := ctx.Value("dbUser").(model.Users)
		if user.ID == "" {
			return nil, twirp.Internal.Error("internal error")
		}

		if user.IsBotAdmin {
			return next(ctx, req)
		}

		dashboardId := ctx.Value("dashboardId")
		if dashboardId == "" {
			dashboardId = s.sessionManager.Get(ctx, "dashboardId").(string)
			ctx = context.WithValue(ctx, "dashboardId", dashboardId)
		}

		if user.ID == dashboardId {
			return next(ctx, req)
		}

		_, ok := lo.Find(user.Roles, func(a model.ChannelRoleUser) bool {
			return a.Role.ChannelID == dashboardId &&
				lo.Contains(a.Role.Permissions, model.RolePermissionCanAccessDashboard.String())
		})

		if ok {
			return next(ctx, req)
		}

		return nil, twirp.NewError(twirp.ErrorCode(twirp.PermissionDenied), "not authorized to that dashboard")
	}
}
