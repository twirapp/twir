package interceptors

import (
	"context"
	"encoding/json"
	model "github.com/satont/tsuwari/libs/gomodels"
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
		user := s.sessionManager.Get(ctx, "dbUser")
		if user != nil {
			parsedUser := &model.Users{}
			err := json.Unmarshal([]byte(user.(string)), parsedUser)
			if err != nil {
				zap.S().Error(err)
				return nil, twirp.Internal.Error("internal error")
			}

			ctx = context.WithValue(ctx, "dbUser", parsedUser)
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
