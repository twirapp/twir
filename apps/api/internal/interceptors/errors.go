package interceptors

import (
	"context"

	"github.com/twirapp/twir/libs/logger"
	"gorm.io/gorm"

	"github.com/twitchtv/twirp"
)

func (s *Service) Errors(next twirp.Method) twirp.Method {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		resp, err := next(ctx, req)

		if err != nil && err != gorm.ErrRecordNotFound {
			s.logger.Error("unexpected error", logger.Error(err))
		}
		return resp, err
	}
}
