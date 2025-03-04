package middlewares

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Auth *auth.Auth
	Gorm *gorm.DB
	Huma huma.API
}

func New(opts Opts) *Middlewares {
	return &Middlewares{
		auth: opts.Auth,
		gorm: opts.Gorm,
		huma: opts.Huma,
	}
}

type Middlewares struct {
	auth *auth.Auth
	gorm *gorm.DB
	huma huma.API
}
