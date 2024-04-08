package interceptors

import (
	"context"

	"github.com/satont/twir/apps/api/internal/helpers"
	"github.com/twitchtv/twirp"
)

func (s *Service) AdminInterceptor(next twirp.Method) twirp.Method {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		user, err := helpers.GetUserModelFromCtx(ctx)
		if err != nil {
			return nil, twirp.NewError(twirp.Unauthenticated, "failed to get user from context")
		}

		if !user.IsBotAdmin {
			return nil, twirp.NewError(twirp.PermissionDenied, "user is not admin")
		}

		return next(ctx, req)
	}
}
