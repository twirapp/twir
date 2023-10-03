package interceptors

import (
	"context"

	"github.com/twitchtv/twirp"
)

func (s *Service) ChannelAccessInterceptor(next twirp.Method) twirp.Method {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return next(ctx, req)
		// method, _ := twirp.MethodName(ctx)
		// if method == "AuthUserProfile" || method == "AuthGetDashboards" {
		// 	return next(ctx, req)
		// }
		//
		// user := ctx.Value("dbUser").(model.Users)
		// if user.ID == "" {
		// 	return nil, twirp.Internal.Error("internal error")
		// }
		//
		// if user.IsBotAdmin {
		// 	return next(ctx, req)
		// }
		//
		// dashboardId := ctx.Value("dashboardId")
		// if dashboardId == "" {
		// 	dashboardId = s.sessionManager.Get(ctx, "dashboardId").(string)
		// 	ctx = context.WithValue(ctx, "dashboardId", dashboardId)
		// }
		//
		// if user.ID == dashboardId {
		// 	return next(ctx, req)
		// }
		//
		// var roles []model.ChannelRoleUser
		//
		// if err := s.db.
		// 	Preload("Role", &model.ChannelRole{ChannelID: dashboardId.(string)}).
		// 	Where(`"userId" = ?`, user.ID).
		// 	Find(&roles).Error; err != nil {
		// 	return nil, fmt.Errorf("cannot get user roles: %w", err)
		// }
		//
		// ok := lo.SomeBy(
		// 	roles,
		// 	func(a model.ChannelRoleUser) bool {
		// 		if a.Role == nil {
		// 			return false
		// 		}
		// 		return a.Role.ChannelID == dashboardId && a.UserID == user.ID
		// 	},
		// )
		//
		// if ok {
		// 	return next(ctx, req)
		// }
		//
		// return nil, twirp.PermissionDenied.Error("not authorized to that dashboard")
	}
}
