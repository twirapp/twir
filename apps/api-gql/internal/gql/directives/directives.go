package directives

import (
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Sessions *auth.Auth
	Gorm     *gorm.DB
}

func New(opts Opts) *Directives {
	return &Directives{
		sessions: opts.Sessions,
		gorm:     opts.Gorm,
	}
}

type Directives struct {
	sessions *auth.Auth
	gorm     *gorm.DB
}
