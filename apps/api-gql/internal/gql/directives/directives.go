package directives

import (
	"github.com/twirapp/twir/apps/api-gql/internal/sessions"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Sessions *sessions.Sessions
	Gorm     *gorm.DB
}

func New(opts Opts) *Directives {
	return &Directives{
		sessions: opts.Sessions,
		gorm:     opts.Gorm,
	}
}

type Directives struct {
	sessions *sessions.Sessions
	gorm     *gorm.DB
}
