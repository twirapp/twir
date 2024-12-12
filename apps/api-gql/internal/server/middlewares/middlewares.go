package middlewares

import (
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Sessions *auth.Auth
	Logger   logger.Logger
}

func New(opts Opts) *Middlewares {
	return &Middlewares{
		sessions: opts.Sessions,
		logger:   opts.Logger,
	}
}

type Middlewares struct {
	sessions *auth.Auth
	logger   logger.Logger
}
