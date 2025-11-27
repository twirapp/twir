package hooks

import (
	"context"
	"log/slog"

	"github.com/twitchtv/twirp"
)

func NewLoggingServerHooks(logger *slog.Logger) *twirp.ServerHooks {
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
