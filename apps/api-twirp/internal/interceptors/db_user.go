package interceptors

import (
	"context"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twitchtv/twirp"
	"go.uber.org/zap"
)

func (s *Service) getUserByApiKey(apiKey string) (*model.Users, error) {
	user := model.Users{}
	err := s.db.
		Where(`"apiKey" = ?`, apiKey).
		Preload("Roles").
		Preload("Roles.Role").
		Find(&user).Error
	if err != nil {
		return nil, err
	}
	if user.ID == "" {
		return nil, nil
	}

	return &user, nil
}

func (s *Service) DbUserInterceptor(next twirp.Method) twirp.Method {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		apiKey := ctx.Value("apiKey")
		sessionUser := s.sessionManager.Get(ctx, "dbUser")
		if sessionUser == nil {
			return nil, twirp.Unauthenticated.Error("not authenticated")
		}

		user := sessionUser.(model.Users)
		if user.ID != "" {
			ctx = context.WithValue(ctx, "dbUser", user)
			return next(ctx, req)
		}

		if apiKey == nil || apiKey == "" {
			return nil, twirp.Unauthenticated.Error("not authenticated")
		}

		castedApiKey, ok := apiKey.(string)
		if !ok {
			return nil, twirp.Internal.Error("internal error")
		}

		dbUser, err := s.getUserByApiKey(castedApiKey)
		if err != nil {
			zap.S().Error(err)
			return nil, twirp.Internal.Error("internal error")
		}
		if dbUser == nil {
			return nil, twirp.Unauthenticated.Error("not authenticated")
		}
		ctx = context.WithValue(ctx, "dbUser", dbUser)

		return next(ctx, req)
	}
}
