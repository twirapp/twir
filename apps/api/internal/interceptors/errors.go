package interceptors

import (
	"context"
	"log/slog"

	"gorm.io/gorm"

	"github.com/twitchtv/twirp"
)

func (s *Service) Errors(next twirp.Method) twirp.Method {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		resp, err := next(ctx, req)

		if err != nil && err != gorm.ErrRecordNotFound {
			s.logger.Error("unexpected error", slog.Any("err", err))
		}
		return resp, err
	}
}
