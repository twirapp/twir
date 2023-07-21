package hooks

import (
	"context"
	"github.com/twitchtv/twirp"
	"go.uber.org/zap"
)

func NewLoggingServerHooks(logger *zap.Logger) *twirp.ServerHooks {
	return &twirp.ServerHooks{
		Error: func(ctx context.Context, twerr twirp.Error) context.Context {
			method, _ := twirp.MethodName(ctx)
			logger.Sugar().Errorw("Error in method", zap.String("method", method), zap.Error(twerr))
			return ctx
		},
	}
}
