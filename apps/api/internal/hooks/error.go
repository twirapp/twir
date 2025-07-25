package hooks

import (
	"context"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twitchtv/twirp"
	"log/slog"
)

func NewLoggingServerHooks(logger logger.Logger) *twirp.ServerHooks {
	return &twirp.ServerHooks{
		Error: func(ctx context.Context, twerr twirp.Error) context.Context {
			method, _ := twirp.MethodName(ctx)
			logger.Error(
				"Error in method",
				slog.String("method", method), slog.Any(
					"err",
					twerr,
				),
			)
			return ctx
		},
	}
}
