package interceptors

import (
	"context"
	"gorm.io/gorm"

	"github.com/twitchtv/twirp"
	"go.uber.org/zap"
)

func (s *Service) Errors(next twirp.Method) twirp.Method {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		resp, err := next(ctx, req)

		if err != nil && err != gorm.ErrRecordNotFound {
			zap.S().Error(err)
		}
		return resp, err
	}
}
